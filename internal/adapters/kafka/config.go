package kafka

import "time"

type config struct {
	AddrsStr string

	Addrs       []string
	maxRetries  int
	maxWaitTime time.Duration

	// Authentication
	SASLUser string
	SASLPass string
}
