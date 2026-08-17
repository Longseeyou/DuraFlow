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
	GetWorkflowDefinitionByID(
		ctx context.Context,
		workflowDefinitionID uuid.UUID,
	) (workflow.WorkflowDefinition, error)
}

type WorkflowRunStoreInternal interface {
	CreateWorkflowRun(
		ctx context.Context,
		workflowRun workflow.WorkflowRun,
	) (workflow.WorkflowRun, error)
}

//

type TaskDefinitionStoreInternal interface {
	GetTaskDefinitionByID(
		ctx context.Context,
		taskDefinitionID uuid.UUID,
	) (task.TaskDefinition, error)
	GetTaskDefinitionsByWorkflowDefinition(
		ctx context.Context,
		workflowDefinitionID uuid.UUID,
	) ([]task.TaskDefinition, error)
}

type TaskDependencyStoreInternal interface {
	// GetTaskDependenciesByWorkflowDefinition(
	// 	ctx context.Context,
	// 	workflowDefinitionID uuid.UUID,
	// ) ([]task.TaskDependency, error)
}

type TaskRunStoreInternal interface {
	CreateTaskRun(ctx context.Context, tR task.TaskRun) (task.TaskRun, error)
	GetTaskRunByID(ctx context.Context, taskRunID uuid.UUID) (task.TaskRun, error)
	// GetTaskRunsByWorkflowRunID(ctx context.Context, workflowRunID uuid.UUID) ([]TaskRun, error)
	UpdateTaskRunByID(
		ctx context.Context,
		taskRunID uuid.UUID,
		newTR map[string]any,
	) (task.TaskRun, error)
	// HardDeleteTaskRun(ctx context.Context, taskRunID uuid.UUID) (TaskRun, error)
	TaskRunIdempotency(
		ctx context.Context,
		taskRunID uuid.UUID,
		newStatus task.TaskRunStatus,
	) (bool, error)

	GetPredecessorTaskRuns(ctx context.Context, taskRunID uuid.UUID) ([]task.TaskRun, error)
}

type TaskAttemptStoreInternal interface {
	CreateTaskAttempt(ctx context.Context, tA task.TaskAttempt) (task.TaskAttempt, error)
	UpdateTaskAttemptByID(
		ctx context.Context,
		taskAttemptID uuid.UUID,
		newTA map[string]any,
	) (task.TaskAttempt, error)
	TaskAttemptIdempotency(
		ctx context.Context,
		taskAttemptID uuid.UUID,
		newStatus task.TaskAttemptStatus,
	) (bool, error)
}
