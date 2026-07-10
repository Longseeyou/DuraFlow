package orchestrator

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

// WorkflowRepository persists workflow-run state required by the
// orchestrator by reusing the workflow package's internal store contract.
type WorkflowRepository interface {
	GetWorkflowDefinitionById(ctx context.Context, id uuid.UUID) (workflow.WorkflowDefinition, error)
	GetWorkflowRunById(ctx context.Context, id uuid.UUID) (workflow.WorkflowRun, error)
	UpdateWorkflowRunByIdAndStatusInternal(
		ctx context.Context,
		id uuid.UUID,
		statuses []workflow.WorkflowRunStatus,
		updates map[string]any,
	) (workflow.WorkflowRun, bool, error)
}

// TaskRepository contains only the persistence operations used by Service.
type TaskRepository interface {
	GetTaskDefinitionById(ctx context.Context, id uuid.UUID) (task.TaskDefinition, error)
	GetTaskDefinitionsByWorkflowDefinition(
		ctx context.Context,
		workflowDefinitionID uuid.UUID,
	) ([]task.TaskDefinition, error)
	GetTaskDependenciesByWorkflowDefinitionId(
		ctx context.Context,
		workflowDefinitionID uuid.UUID,
	) ([]task.TaskDependency, error)
	GetTaskRunById(ctx context.Context, id uuid.UUID) (task.TaskRun, error)
	GetTaskRunsByWorkflowRunId(ctx context.Context, workflowRunID uuid.UUID) ([]task.TaskRun, error)
	UpdateTaskRunByIdAndStatus(
		ctx context.Context,
		id uuid.UUID,
		statuses []task.TaskRunStatus,
		updates map[string]any,
	) (task.TaskRun, bool, error)
	CreateTaskAttempt(ctx context.Context, attempt task.TaskAttempt) (task.TaskAttempt, error)
	HardDeleteTaskAttempt(ctx context.Context, id uuid.UUID) (task.TaskAttempt, error)
	GetTaskAttemptByTaskRunAndNumber(
		ctx context.Context,
		taskRunID uuid.UUID,
		attemptNumber uint,
	) (task.TaskAttempt, error)
	UpdateTaskAttemptById(
		ctx context.Context,
		id uuid.UUID,
		updates map[string]any,
	) (task.TaskAttempt, error)
	TaskEventExists(ctx context.Context, eventID uuid.UUID) (bool, error)
	CreateTaskEvent(ctx context.Context, event task.TaskEvent) (task.TaskEvent, error)
}

// TaskPollingRepository is the smaller write contract needed by TaskPoller.
type TaskPollingRepository interface {
	GetTaskDefinitionById(ctx context.Context, id uuid.UUID) (task.TaskDefinition, error)
	UpdateTaskRunByIdAndStatus(
		ctx context.Context,
		id uuid.UUID,
		statuses []task.TaskRunStatus,
		updates map[string]any,
	) (task.TaskRun, bool, error)
	CreateTaskAttempt(ctx context.Context, attempt task.TaskAttempt) (task.TaskAttempt, error)
	HardDeleteTaskAttempt(ctx context.Context, id uuid.UUID) (task.TaskAttempt, error)
}

// TaskPublisher dispatches runnable task commands to workers.
type TaskPublisher interface {
	PublishTask(ctx context.Context, command task.Command) error
}
