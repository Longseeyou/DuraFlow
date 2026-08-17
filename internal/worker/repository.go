package worker

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/google/uuid"
)

type WorkerRepository interface {
	TaskRunStoreInternal
	TaskAttemptStoreInternal
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
	) (bool, error)
}

type TaskAttemptStoreInternal interface {
	CreateTaskAttempt(ctx context.Context, tA task.TaskAttempt) (task.TaskAttempt, error)
	// HardDeleteTaskAttempt(ctx context.Context, taskAttemptID uuid.UUID) (TaskAttempt, error)
	// GetTaskAttemptByTaskRunAndNumber(
	// 	ctx context.Context,
	// 	taskRunID uuid.UUID,
	// 	attemptNumber uint,
	// ) (TaskAttempt, error)
	// UpdateTaskAttemptByID(
	// 	ctx context.Context,
	// 	taskAttemptID uuid.UUID,
	// 	newTaskAttempt map[string]any,
	// ) (TaskAttempt, error)
	// TaskAttemptIdempotency(
	// 	ctx context.Context,
	// 	taskAttemptID uuid.UUID,
	// 	newStatus task.TaskAttemptStatus,
	// ) (bool, error)
	// GetTimedOutTaskAttempts(ctx context.Context) ([]TaskAttempt, error)
}
