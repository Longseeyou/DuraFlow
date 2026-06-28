package kafka

import "time"

type KafkaConfig struct {
	AddrsStr string

	Addrs       []string
	MaxRetries  int
	MaxWaitTime time.Duration

	// Authentication
	SASLUser string
	SASLPass string
}
