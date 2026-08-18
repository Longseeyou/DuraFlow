package worker

import (
	"context"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/google/uuid"
)

type WorkerRepository interface {
	TaskRunStoreInternal
	TaskAttemptStoreInternal

	WorkerStoreInternal
}

type TaskRunStoreInternal interface {
	GetTaskRunByID(ctx context.Context, taskRunID uuid.UUID) (task.TaskRun, error)
	UpdateTaskRunByID(
		ctx context.Context,
		taskRunID uuid.UUID,
		newTaskRun map[string]any,
	) (task.TaskRun, error)
	TaskRunIdempotency(
		ctx context.Context,
		taskRunID uuid.UUID,
		newStatus task.TaskRunStatus,
	) (bool, error)
	MarkTaskRunRunning(
		ctx context.Context,
		taskRunID uuid.UUID,
		attemptNumber uint,
		startedAt time.Time,
	) (bool, task.TaskRun, error)
}

type TaskAttemptStoreInternal interface {
	CreateTaskAttempt(ctx context.Context, tA task.TaskAttempt) (task.TaskAttempt, error)
}

type WorkerStoreInternal interface {
	UpdateTaskRunByIDAndCreateTaskAttempt(
		ctx context.Context,
		taskRunID uuid.UUID,
		newTaskRun map[string]any,
		taskAttempt task.TaskAttempt,
	) (task.TaskRun, task.TaskAttempt, error)
}
