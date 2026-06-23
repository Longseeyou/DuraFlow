package kafka

import (
	"context"
	"log/slog"
	"strings"

	"github.com/IBM/sarama"
)

type producerKafka struct {
	id string
	*config
	producer *sarama.SyncProducer
}

func NewProducerKafka(id string) *producerKafka {
	return &producerKafka{
		id:     id,
		config: new(config),
	}
}

func (p *producerKafka) Start(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Retry.Max = p.maxRetries
	config.Producer.Timeout = p.maxWaitTime

	// set authentication
	if p.SASLUser != "" && p.SASLPass != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = p.SASLUser
		config.Net.SASL.Password = p.SASLPass
	}

	// Parse the addresses
	p.Addrs = strings.Split(p.AddrsStr, ",")

	// Create the producer
	slog.Info("Creating Kafka producer", "addresses", p.Addrs)
	producer, err := sarama.NewSyncProducer(p.Addrs, config)
	if err != nil {
		return err
	}

	p.producer = &producer
	slog.Info("Kafka producer created successfully", "addresses", p.Addrs)

	return nil
}

func (p *producerKafka) Stop() error {
	if p.producer != nil {
		slog.Info("Stopping Kafka producer")
		if err := (*p.producer).Close(); err != nil {
			slog.Error("Failed to close Kafka producer", "error", err)
		}
	}
	return nil
}

func (p *producerKafka) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	if p.producer == nil {
		return sarama.ErrNotConnected
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	slog.Info(
		"Sending message to Kafka",
		"topic",
		topic,
		"key",
		string(key),
		"value",
		string(value),
	)

	partition, offset, err := (*p.producer).SendMessage(msg)
	if err != nil {
		return err
	}

	slog.Info("Message sent successfully", "topic", topic, "partition", partition, "offset", offset)

	return nil
}
