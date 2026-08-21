package kafka

import "time"

type KafkaConfig struct {
	Addrress    []string
	MaxRetries  int
	MaxWaitTime time.Duration

	// Authentication
	SASLUser string
	SASLPass string
}

func NewKafkaConfig(
	address []string,
	maxRetries int,
	maxWaitTime time.Duration,
	SASLUser *string,
	SASLPass *string,
) KafkaConfig {
	var kafkaConfig KafkaConfig
	kafkaConfig.Addrress = address
	kafkaConfig.MaxRetries = maxRetries
	kafkaConfig.MaxWaitTime = maxWaitTime

	if SASLUser != nil && SASLPass != nil {
		kafkaConfig.SASLUser = *SASLUser
		kafkaConfig.SASLPass = *SASLPass
	}

	return kafkaConfig
}
