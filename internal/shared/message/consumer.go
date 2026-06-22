package message

import (
	"context"
)

type Message struct {
	Topic     string
	Key       []byte
	Value     []byte
	Partition int32
	Offset    int64
}

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
	ReceiveMessage(ctx context.Context) (*Message, error)
}
