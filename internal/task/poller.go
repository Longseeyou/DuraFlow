package task

import (
	"context"
)

type TaskPoller interface {
	PollTask(ctx context.Context) (TaskRun, error)
}
