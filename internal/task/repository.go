package task

import (
	"context"

	"github.com/google/uuid"
)

type TaskRepository interface {
	CreateTaskDefinition(ctx context.Context, tD TaskDefinition) (TaskDefinition, error)
	GetTaskDefinitionByWorkflowDefinition(ctx context.Context, wDId uuid.UUID) ([]TaskDefinition, error)
	GetTaskDefinitionByWorkflowDefinitionAndId(ctx context.Context, wDId uuid.UUID, tDId uuid.UUID) (TaskDefinition, error)
	UpdateTaskDefinitionById(ctx context.Context, wDId uuid.UUID, tDId uuid.UUID, newTD map[string]any) (TaskDefinition, error)
	SoftDeleteTaskDefinition(ctx context.Context, wDId uuid.UUID, tDId uuid.UUID) (TaskDefinition, error)
	HardDeleteTaskDefinition(ctx context.Context, wDId uuid.UUID, tDId uuid.UUID) (TaskDefinition, error)

	CreateTaskDependency(ctx context.Context, tDp TaskDependency) (TaskDependency, error)
	GetTaskDependencyByWorkflowDefinition(ctx context.Context, wDId uuid.UUID) ([]TaskDependency, error)
	GetTaskDependencyByWorkflowDefinitionAndId(ctx context.Context, wDId uuid.UUID, tDpId uuid.UUID) (TaskDependency, error)
	UpdateTaskDependencyById(ctx context.Context, tDpId uuid.UUID, newTDp map[string]any) (TaskDependency, error)
	SoftDeleteTaskDependency(ctx context.Context, tDpId uuid.UUID) (TaskDependency, error)
	HardDeleteTaskDependency(ctx context.Context, tDpId uuid.UUID) (TaskDependency, error)

	CreateTaskRun(ctx context.Context, tR TaskRun) (TaskRun, error)
	GetTaskRunByWorkflowRun(ctx context.Context, wRId uuid.UUID) ([]TaskRun, error)
	GetTaskRunByWorkflowRunAndId(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID) (TaskRun, error)
	UpdateTaskRunById(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID, newTR map[string]any) (TaskRun, error)
	SoftDeleteTaskRun(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID) (TaskRun, error)
	HardDeleteTaskRun(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID) (TaskRun, error)

	CreateTaskAttempt(ctx context.Context, tA TaskAttempt) (TaskAttempt, error)
	GetTaskAttemptByTaskRun(ctx context.Context, tRId uuid.UUID) ([]TaskAttempt, error)
	GetTaskAttemptByTaskRunAndId(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID) (TaskAttempt, error)
	UpdateTaskAttemptById(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID, newTA map[string]any) (TaskAttempt, error)
	SoftDeleteTaskAttempt(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID) (TaskAttempt, error)
	HardDeleteTaskAttempt(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID) (TaskAttempt, error)

	// CreateTaskEvent(ctx context.Context, tE TaskEvent) (TaskEvent, error)
	// GetTaskEventByWorkflowRun(ctx context.Context, wRId uuid.UUID) ([]TaskEvent, error)
	// GetTaskEventByWorkflowRunAndId(ctx context.Context, wRId uuid.UUID, tEId uuid.UUID) (TaskEvent, error)
	// UpdateTaskEventById(ctx context.Context, wRId uuid.UUID, tEId uuid.UUID, newTE map[string]any) (TaskEvent, error)
	// SoftDeleteTaskEvent(ctx context.Context, wRId uuid.UUID, tEId uuid.UUID) (TaskEvent, error)
	// HardDeleteTaskEvent(ctx context.Context, wRId uuid.UUID, tEId uuid.UUID) (TaskEvent, error)
}

type TaskRepositoryInternal interface {
	GetTaskDefinitionById(ctx context.Context, tDId uuid.UUID) (TaskDefinition, error)

	GetTaskDependencyById(ctx context.Context, tDpId uuid.UUID) (TaskDependency, error)
}
