package message

import (
	"context"
)

type Producer interface {
	Start(ctx context.Context) error
	Stop() error
	SendMessage(ctx context.Context, topic string, key string, value []byte) error
}
