package task

import (
	"context"
)

type TaskExecutor interface {
	Execute(ctx context.Context, input string) (string, string, TaskAttemptStatus, error)
}
