package orchestrator

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/task"
)

type TaskPoller interface {
	PollTaskRun(ctx context.Context) ([]task.TaskRun, error)
}
