package kafka

import (
	"context"
	"log/slog"
	"strings"

	"github.com/IBM/sarama"
	"github.com/Longseeyou/DuraFlow/internal/shared/message"
)

type consumerKafka struct {
	id string
	*config

	groupID string
	topics  []string

	consumerGroup sarama.ConsumerGroup

	messages chan *message.Message

	ctxCancel context.CancelFunc
}

func NewConsumerKafka(id, groupID string, topics []string) *consumerKafka {
	return &consumerKafka{
		id:       id,
		groupID:  groupID,
		topics:   topics,
		config:   new(config),
		messages: make(chan *message.Message, 100),
	}
}

func (c *consumerKafka) Start(ctx context.Context) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0

	// Authentication
	if c.SASLUser != "" && c.SASLPass != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = c.SASLUser
		cfg.Net.SASL.Password = c.SASLPass
	}

	c.Addrs = strings.Split(c.AddrsStr, ",")

	group, err := sarama.NewConsumerGroup(c.Addrs, c.groupID, cfg)
	if err != nil {
		slog.Error("Failed to create consumer group", "error", err)
		return
	}

	c.consumerGroup = group

	runCtx, cancel := context.WithCancel(ctx)
	c.ctxCancel = cancel

	go func() {
		handler := &consumerGroupHandler{
			consumer: c,
		}

		for {
			if err := group.Consume(runCtx, c.topics, handler); err != nil {
				slog.Error("Consumer error", "error", err)
			}

			if runCtx.Err() != nil {
				return
			}
		}
	}()

	slog.Info(
		"Kafka consumer started",
		"groupID",
		c.groupID,
		"topics",
		c.topics,
	)
}

func (c *consumerKafka) Stop() {
	if c.ctxCancel != nil {
		c.ctxCancel()
	}

	if c.consumerGroup != nil {
		if err := c.consumerGroup.Close(); err != nil {
			slog.Error("Failed to close consumer group", "error", err)
		}
	}

	slog.Info("Kafka consumer stopped")
}

func (c *consumerKafka) ReceiveMessage(
	ctx context.Context,
) (*message.Message, error) {
	select {
	case msg := <-c.messages:
		return msg, nil

	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

type consumerGroupHandler struct {
	consumer *consumerKafka
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		select {
		case h.consumer.messages <- &message.Message{
			Topic:     msg.Topic,
			Key:       msg.Key,
			Value:     msg.Value,
			Partition: msg.Partition,
			Offset:    msg.Offset,
		}:
			session.MarkMessage(msg, "")

		case <-session.Context().Done():
			return nil
		}
	}

	return nil
}
