package orchestrator

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type OrchestratorRepository interface {
	WorkflowDefinitionStoreInternal
	WorkflowRunStoreInternal

	TaskDefinitionStoreInternal
	TaskDependencyStoreInternal
	TaskRunStoreInternal
	TaskAttemptStoreInternal
}

type WorkflowDefinitionStoreInternal interface {
	GetWorkflowDefinitionById(
		ctx context.Context,
		wDId uuid.UUID,
	) (workflow.WorkflowDefinition, error)
}

type WorkflowRunStoreInternal interface {
	CreateWorkflowRun(
		ctx context.Context,
		wR workflow.WorkflowRun,
	) (workflow.WorkflowRun, error)
}

//

type TaskDefinitionStoreInternal interface {
	GetTaskDefinitionsByWorkflowDefinition(
		ctx context.Context,
		wDId uuid.UUID,
	) ([]task.TaskDefinition, error)
}

type TaskDependencyStoreInternal interface {
	// GetTaskDependenciesByWorkflowDefinition(
	// 	ctx context.Context,
	// 	wDId uuid.UUID,
	// ) ([]task.TaskDependency, error)
}

type TaskRunStoreInternal interface {
	CreateTaskRun(ctx context.Context, tR task.TaskRun) (task.TaskRun, error)
	GetTaskRunById(ctx context.Context, tRId uuid.UUID) (task.TaskRun, error)
	// GetTaskRunsByWorkflowRunId(ctx context.Context, wRId uuid.UUID) ([]TaskRun, error)
	UpdateTaskRunById(
		ctx context.Context,
		tRId uuid.UUID,
		newTR map[string]any,
	) (task.TaskRun, error)
	// HardDeleteTaskRun(ctx context.Context, tRId uuid.UUID) (TaskRun, error)
	TaskRunIdempotency(
		ctx context.Context,
		tRId uuid.UUID,
		newStatus task.TaskRunStatus,
	) (bool, error)

	GetPredecessorTaskRuns(ctx context.Context, tRId uuid.UUID) ([]task.TaskRun, error)
}

type TaskAttemptStoreInternal interface {
	GetTaskAttemptId(ctx context.Context, tAId uuid.UUID) (task.TaskAttempt, error)
	CreateTaskAttempt(ctx context.Context, tA task.TaskAttempt) (task.TaskAttempt, error)
	UpdateTaskAttemptById(
		ctx context.Context,
		tAId uuid.UUID,
		newTA map[string]any,
	) (task.TaskAttempt, error)
	TaskAttemptIdempotency(
		ctx context.Context,
		tRId uuid.UUID,
		newStatus task.TaskRunStatus,
	) (bool, error)
}
