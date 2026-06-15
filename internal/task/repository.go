package task

import (
	"context"

	"github.com/google/uuid"
)

type TaskRepository interface {
	TaskDefinitionStore
	TaskDependencyStore
	TaskRunStore
	TaskAttemptStore
}

type TaskDefinitionStore interface {
	CreateTaskDefinition(ctx context.Context, tD TaskDefinition) (TaskDefinition, error)
	GetTaskDefinitionByWorkflowDefinition(
		ctx context.Context,
		uId uuid.UUID,
		wDId uuid.UUID,
	) ([]TaskDefinition, error)
	GetTaskDefinitionByUserAndId(
		ctx context.Context,
		uId uuid.UUID,
		tDId uuid.UUID,
	) (TaskDefinition, error)
	UpdateTaskDefinitionById(
		ctx context.Context,
		uId uuid.UUID,
		tDId uuid.UUID,
		newTD map[string]any,
	) (TaskDefinition, error)
	SoftDeleteTaskDefinition(
		ctx context.Context,
		uId uuid.UUID,
		tDId uuid.UUID,
	) (TaskDefinition, error)
	HardDeleteTaskDefinition(
		ctx context.Context,
		uId uuid.UUID,
		tDId uuid.UUID,
	) (TaskDefinition, error)
}

type TaskDependencyStore interface {
	CreateTaskDependency(ctx context.Context, tDp TaskDependency) (TaskDependency, error)
	GetTaskDependencyByWorkflowDefinition(
		ctx context.Context,
		uId uuid.UUID,
		wDId uuid.UUID,
	) ([]TaskDependency, error)
	GetTaskDependencyByUserAndId(
		ctx context.Context,
		uId uuid.UUID,
		tDpId uuid.UUID,
	) (TaskDependency, error)
	UpdateTaskDependencyById(
		ctx context.Context,
		uId uuid.UUID,
		tDpId uuid.UUID,
		newTDp map[string]any,
	) (TaskDependency, error)
	SoftDeleteTaskDependency(
		ctx context.Context,
		uId uuid.UUID,
		tDpId uuid.UUID,
	) (TaskDependency, error)
	HardDeleteTaskDependency(
		ctx context.Context,
		uId uuid.UUID,
		tDpId uuid.UUID,
	) (TaskDependency, error)
}

type TaskRunStore interface {
	CreateTaskRun(ctx context.Context, tR TaskRun) (TaskRun, error)
	GetTaskRunByWorkflowRun(ctx context.Context, uId uuid.UUID, wRId uuid.UUID) ([]TaskRun, error)
	GetTaskRunByUserAndId(
		ctx context.Context,
		uId uuid.UUID,
		tRId uuid.UUID,
	) (TaskRun, error)
	UpdateTaskRunById(
		ctx context.Context,
		uId uuid.UUID,
		tRId uuid.UUID,
		newTR map[string]any,
	) (TaskRun, error)
	SoftDeleteTaskRun(ctx context.Context, uId uuid.UUID, tRId uuid.UUID) (TaskRun, error)
	HardDeleteTaskRun(ctx context.Context, uId uuid.UUID, tRId uuid.UUID) (TaskRun, error)
}

type TaskAttemptStore interface {
	CreateTaskAttempt(ctx context.Context, tA TaskAttempt) (TaskAttempt, error)
	GetTaskAttemptByTaskRun(
		ctx context.Context,
		uId uuid.UUID,
		tRId uuid.UUID,
	) ([]TaskAttempt, error)
	GetTaskAttemptByUserAndId(
		ctx context.Context,
		uId uuid.UUID,
		tAId uuid.UUID,
	) (TaskAttempt, error)
	UpdateTaskAttemptById(
		ctx context.Context,
		uId uuid.UUID,
		tAId uuid.UUID,
		newTA map[string]any,
	) (TaskAttempt, error)
	SoftDeleteTaskAttempt(ctx context.Context, uId uuid.UUID, tAId uuid.UUID) (TaskAttempt, error)
	HardDeleteTaskAttempt(ctx context.Context, uId uuid.UUID, tAId uuid.UUID) (TaskAttempt, error)
}

type TaskRepositoryInternal interface {
	TaskDefinitionStoreInternal
	TaskDependencyStoreInternal
}

type TaskDefinitionStoreInternal interface {
	GetTaskDefinitionById(ctx context.Context, tDId uuid.UUID) (TaskDefinition, error)
}

type TaskDependencyStoreInternal interface {
	GetTaskDependencyById(ctx context.Context, tDpId uuid.UUID) (TaskDependency, error)
}
