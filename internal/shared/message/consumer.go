package message

import (
	"context"
)

type Consumer interface {
	Start(ctx context.Context)
	Stop()
	ReceiveMessage(ctx context.Context) (*Message, error)
}
