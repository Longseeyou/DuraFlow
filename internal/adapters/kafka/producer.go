package kafka

import (
	"context"
	"log/slog"

	"github.com/IBM/sarama"
)

type kafkaProducer struct {
	id string
	*KafkaConfig
	producer *sarama.SyncProducer
}

func NewKafkaProducer(id string, config *KafkaConfig) *kafkaProducer {
	return &kafkaProducer{
		id:          id,
		KafkaConfig: config,
	}
}

func (p *kafkaProducer) Start(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Retry.Max = p.MaxRetries
	config.Producer.Timeout = p.MaxWaitTime

	// set authentication
	if p.SASLUser != "" && p.SASLPass != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = p.SASLUser
		config.Net.SASL.Password = p.SASLPass
	}

	// Create the producer
	slog.Info("Kafka producer", "addresses", p.Addrress)
	producer, err := sarama.NewSyncProducer(p.Addrress, config)
	if err != nil {
		return err
	}

	p.producer = &producer

	return nil
}

func (p *kafkaProducer) Stop() error {
	if p.producer != nil {
		err := (*p.producer).Close()
		if err != nil {
			slog.Error("Failed to close Kafka producer", "error", err)
		}
	}
	return nil
}

func (p *kafkaProducer) SendMessage(
	ctx context.Context,
	topic string,
	key string,
	value []byte,
) error {
	if p.producer == nil {
		return sarama.ErrNotConnected
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
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
