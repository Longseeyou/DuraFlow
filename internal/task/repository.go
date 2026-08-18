package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskRepository interface {
	TaskDefinitionStore
	TaskDependencyStore
	TaskRunStore
	TaskAttemptStore
}

type TaskDefinitionStore interface {
	CreateTaskDefinition(
		ctx context.Context,
		userID uuid.UUID,
		workflowDefinitionID uuid.UUID,
		tD TaskDefinition,
	) (TaskDefinition, error)
	GetTaskDefinitionsByWorkflowDefinition(
		ctx context.Context,
		userID uuid.UUID,
		workflowDefinitionID uuid.UUID,
	) ([]TaskDefinition, error)
	GetTaskDefinitionByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		taskDefinitionID uuid.UUID,
	) (TaskDefinition, error)
	UpdateTaskDefinitionByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		taskDefinitionID uuid.UUID,
		newTaskDefinition map[string]any,
	) (TaskDefinition, error)
	SoftDeleteTaskDefinition(
		ctx context.Context,
		userID uuid.UUID,
		taskDefinitionID uuid.UUID,
	) (TaskDefinition, error)
}

type TaskDependencyStore interface {
	CreateTaskDependency(ctx context.Context, taskDependency TaskDependency) (TaskDependency, error)
	GetTaskDependenciesByWorkflowDefinition(
		ctx context.Context,
		userID uuid.UUID,
		workflowDefinitionID uuid.UUID,
	) ([]TaskDependency, error)
	GetTaskDependencyByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		taskDependencyID uuid.UUID,
	) (TaskDependency, error)
	UpdateTaskDependencyByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		taskDependencyID uuid.UUID,
		newTaskDependency map[string]any,
	) (TaskDependency, error)
	SoftDeleteTaskDependency(
		ctx context.Context,
		userID uuid.UUID,
		taskDependencyID uuid.UUID,
	) (TaskDependency, error)
}

type TaskRunStore interface {
	GetTaskRunsByWorkflowRun(
		ctx context.Context,
		userID uuid.UUID,
		workflowRunID uuid.UUID,
	) ([]TaskRun, error)
	GetTaskRunByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		taskRunID uuid.UUID,
	) (TaskRun, error)
	UpdateTaskRunByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		taskRunID uuid.UUID,
		newTaskRun map[string]any,
	) (TaskRun, error)
	SoftDeleteTaskRun(ctx context.Context, userID uuid.UUID, taskRunID uuid.UUID) (TaskRun, error)
}

type TaskAttemptStore interface {
	GetTaskAttemptsByTaskRun(
		ctx context.Context,
		userID uuid.UUID,
		taskRunID uuid.UUID,
	) ([]TaskAttempt, error)
	GetTaskAttemptByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		taskAttemptID uuid.UUID,
	) (TaskAttempt, error)
	SoftDeleteTaskAttempt(
		ctx context.Context,
		userID uuid.UUID,
		taskAttemptID uuid.UUID,
	) (TaskAttempt, error)
}

type TaskRepositoryInternal interface {
	TaskDefinitionStoreInternal
	TaskDependencyStoreInternal
	TaskRunStoreInternal
	TaskAttemptStoreInternal
}

type TaskDefinitionStoreInternal interface {
	GetTaskDefinitionByID(ctx context.Context, taskDefinitionID uuid.UUID) (TaskDefinition, error)
	GetTaskDefinitionsByWorkflowDefinition(
		ctx context.Context,
		workflowDefinitionID uuid.UUID,
	) ([]TaskDefinition, error)
	HardDeleteTaskDefinition(
		ctx context.Context,
		taskDefinitionID uuid.UUID,
	) (TaskDefinition, error)
}

type TaskDependencyStoreInternal interface {
	GetTaskDependencyByID(ctx context.Context, taskDependencyID uuid.UUID) (TaskDependency, error)
	GetTaskDependenciesByWorkflowDefinition(
		ctx context.Context,
		workflowDefinitionID uuid.UUID,
	) ([]TaskDependency, error)
	HardDeleteTaskDependency(
		ctx context.Context,
		taskDependencyID uuid.UUID,
	) (TaskDependency, error)
}

type TaskRunStoreInternal interface {
	CreateTaskRun(ctx context.Context, tR TaskRun) (TaskRun, error)
	GetTaskRunByID(ctx context.Context, taskRunID uuid.UUID) (TaskRun, error)
	UpdateTaskRunByID(
		ctx context.Context,
		taskRunID uuid.UUID,
		newTaskRun map[string]any,
	) (TaskRun, error)
	HardDeleteTaskRun(ctx context.Context, taskRunID uuid.UUID) (TaskRun, error)
	TaskRunIdempotency(
		ctx context.Context,
		taskRunID uuid.UUID,
		newStatus TaskRunStatus,
	) (bool, error)
	MarkTaskRunRunning(
		ctx context.Context,
		taskRunID uuid.UUID,
		attemptNumber uint,
		startedAt time.Time,
	) (bool, TaskRun, error)

	GetPredecessorTaskRuns(ctx context.Context, taskRunID uuid.UUID) ([]TaskRun, error)
	GetTimedOutTaskRuns(ctx context.Context) ([]TaskRun, error)
}

type TaskAttemptStoreInternal interface {
	CreateTaskAttempt(ctx context.Context, tA TaskAttempt) (TaskAttempt, error)
	HardDeleteTaskAttempt(ctx context.Context, taskAttemptID uuid.UUID) (TaskAttempt, error)
	GetTaskAttemptByTaskRunAndNumber(
		ctx context.Context,
		taskRunID uuid.UUID,
		attemptNumber uint,
	) (TaskAttempt, error)
	UpdateTaskAttemptByID(
		ctx context.Context,
		taskAttemptID uuid.UUID,
		newTaskAttempt map[string]any,
	) (TaskAttempt, error)
}
