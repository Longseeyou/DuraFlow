package orchestrator

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
)

// WorkflowRepository persists workflow-run state required by the
// orchestrator by reusing the workflow package's internal store contract.
type WorkflowRepository interface {
	workflow.WorkflowDefinitionStoreInternal
	workflow.WorkflowRunStoreInternal
}

// TaskRepository reuses the task package's internal aggregate repository for
// definitions, runs, dependencies, attempts, and events.
type TaskRepository interface {
	task.TaskRepositoryInternal
}

// TaskPublisher dispatches runnable task commands to workers.
type TaskPublisher interface {
	PublishTask(ctx context.Context, command task.Command) error
}
