package message

import (
	"context"
)

type Producer interface {
	Start(ctx context.Context)
	Stop()
	SendMessage(ctx context.Context, topic string, key []byte, value []byte) error
}
