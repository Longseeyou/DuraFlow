package message

import (
	"context"
)

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
	ReceiveMessage(ctx context.Context) (*Message, error)
}
