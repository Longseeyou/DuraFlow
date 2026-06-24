package kafka

import (
	"context"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/Longseeyou/DuraFlow/internal/shared/message"
)

type kafkaConsumer struct {
	id string
	*KafkaConfig

	groupID string
	topics  []string

	consumerGroup sarama.ConsumerGroup

	messages chan *message.Message

	ctxCancel context.CancelFunc
}

func NewKafkaConsumer(
	id string,
	config *KafkaConfig,
	groupID string,
	topics []string,
) message.Consumer {
	return &kafkaConsumer{
		id:          id,
		KafkaConfig: config,
		groupID:     groupID,
		topics:      topics,
		messages:    make(chan *message.Message, 100),
	}
}

func (c *kafkaConsumer) Start(ctx context.Context) error {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0

	// Authentication
	if c.SASLUser != "" && c.SASLPass != "" {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = c.SASLUser
		cfg.Net.SASL.Password = c.SASLPass
	}

	group, err := sarama.NewConsumerGroup(c.Addrress, c.groupID, cfg)
	if err != nil {
		return err
	}

	c.consumerGroup = group

	runCtx, cancel := context.WithCancel(ctx)
	c.ctxCancel = cancel

	go func() {
		handler := &consumerGroupHandler{
			consumer: c,
		}

		for {
			err := group.Consume(runCtx, c.topics, handler)
			if err != nil {
				slog.Error("Kafka consumer error", "error", err)
			}

			if runCtx.Err() != nil {
				return
			}
		}
	}()
	return nil
}

func (c *kafkaConsumer) Stop() error {
	if c.ctxCancel != nil {
		c.ctxCancel()
	}

	if c.consumerGroup != nil {
		err := c.consumerGroup.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *kafkaConsumer) ReceiveMessage(
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
	consumer *kafkaConsumer
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
