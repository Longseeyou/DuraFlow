package workflow

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowRepository interface {
	WorkflowStore
	WorkflowDefinitionStore
	WorkflowRunStore
}

type WorkflowStore interface {
	CreateWorkflow(ctx context.Context, w Workflow) (Workflow, error)
	GetWorkflowByUser(ctx context.Context, uId uuid.UUID) ([]Workflow, error)
	GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (Workflow, error)
	UpdateWorkflowById(
		ctx context.Context,
		uId uuid.UUID,
		wId uuid.UUID,
		newW map[string]any,
	) (Workflow, error)
	SoftDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (Workflow, error)
	HardDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (Workflow, error)
}

type WorkflowDefinitionStore interface {
	CreateWorkflowDefinition(ctx context.Context, wD WorkflowDefinition) (WorkflowDefinition, error)
	GetWorkflowDefinitionByWorkflow(
		ctx context.Context,
		uId uuid.UUID,
		wId uuid.UUID,
	) ([]WorkflowDefinition, error)
	GetWorkflowDefinitionByUserAndId(
		ctx context.Context,
		uId uuid.UUID,
		wDId uuid.UUID,
	) (WorkflowDefinition, error)
	UpdateWorkflowDefinitionById(
		ctx context.Context,
		uId uuid.UUID,
		wDId uuid.UUID,
		newWD map[string]any,
	) (WorkflowDefinition, error)
	SoftDeleteWorkflowDefinition(
		ctx context.Context,
		uId uuid.UUID,
		wDId uuid.UUID,
	) (WorkflowDefinition, error)
	HardDeleteWorkflowDefinition(
		ctx context.Context,
		uId uuid.UUID,
		wDId uuid.UUID,
	) (WorkflowDefinition, error)
}

type WorkflowRunStore interface {
	CreateWorkflowRun(ctx context.Context, wR WorkflowRun) (WorkflowRun, error)
	GetWorkflowRunByWorkflowDefinition(
		ctx context.Context,
		uId uuid.UUID,
		wDId uuid.UUID,
	) ([]WorkflowRun, error)
	GetWorkflowRunByUserAndId(
		ctx context.Context,
		uId uuid.UUID,
		wRId uuid.UUID,
	) (WorkflowRun, error)
	UpdateWorkflowRunById(
		ctx context.Context,
		uId uuid.UUID,
		wRId uuid.UUID,
		newWR map[string]any,
	) (WorkflowRun, error)
	SoftDeleteWorkflowRun(ctx context.Context, uId uuid.UUID, wRId uuid.UUID) (WorkflowRun, error)
	HardDeleteWorkflowRun(ctx context.Context, uId uuid.UUID, wRId uuid.UUID) (WorkflowRun, error)
}

type WorkflowRepositoryInternal interface {
	WorkflowStoreInternal
	WorkflowDefinitionStoreInternal
	WorkflowRunStoreInternal
}

type WorkflowStoreInternal interface {
	GetWorkflowById(ctx context.Context, wId uuid.UUID) (Workflow, error)
}

type WorkflowDefinitionStoreInternal interface {
	GetWorkflowDefinitionById(ctx context.Context, wDId uuid.UUID) (WorkflowDefinition, error)
}

type WorkflowRunStoreInternal interface {
	GetWorkflowRunById(ctx context.Context, wRId uuid.UUID) (WorkflowRun, error)
}
