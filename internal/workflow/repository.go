package workflow

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowRepository interface {
	CreateWorkflow(ctx context.Context, w Workflow) (Workflow, error)
	GetWorkflowByUser(ctx context.Context, uId uuid.UUID) ([]Workflow, error)
	GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (Workflow, error)
	UpdateWorkflowById(ctx context.Context, uId uuid.UUID, wId uuid.UUID, newW map[string]any) (Workflow, error)
	SoftDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (Workflow, error)
	HardDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (Workflow, error)

	CreateWorkflowDefinition(ctx context.Context, wD WorkflowDefinition) (WorkflowDefinition, error)
	GetWorkflowDefinitionByWorkflow(ctx context.Context, wId uuid.UUID) ([]WorkflowDefinition, error)
	GetWorkflowDefinitionByWorkflowAndId(ctx context.Context, wId uuid.UUID, wDId uuid.UUID) (WorkflowDefinition, error)
	UpdateWorkflowDefinitionById(ctx context.Context, wId uuid.UUID, wDId uuid.UUID, newWD map[string]any) (WorkflowDefinition, error)
	SoftDeleteWorkflowDefinition(ctx context.Context, wId uuid.UUID, wDId uuid.UUID) (WorkflowDefinition, error)
	HardDeleteWorkflowDefinition(ctx context.Context, wId uuid.UUID, wDId uuid.UUID) (WorkflowDefinition, error)

	CreateWorkflowRun(ctx context.Context, wR WorkflowRun) (WorkflowRun, error)
	GetWorkflowRunByWorkflowDefinition(ctx context.Context, wDId uuid.UUID) ([]WorkflowRun, error)
	GetWorkflowRunByWorkflowDefinitionAndId(ctx context.Context, wDId uuid.UUID, wRId uuid.UUID) (WorkflowRun, error)
	UpdateWorkflowRunById(ctx context.Context, wDId uuid.UUID, wRId uuid.UUID, newWR map[string]any) (WorkflowRun, error)
	SoftDeleteWorkflowRun(ctx context.Context, wDId uuid.UUID, wRId uuid.UUID) (WorkflowRun, error)
	HardDeleteWorkflowRun(ctx context.Context, wDId uuid.UUID, wRId uuid.UUID) (WorkflowRun, error)
}

type WorkflowRepositoryInternal interface {
	GetWorkflowById(ctx context.Context, wId uuid.UUID) (Workflow, error)
	GetWorkflowDefinitionById(ctx context.Context, wDId uuid.UUID) (WorkflowDefinition, error)
	GetWorkflowRunById(ctx context.Context, wRId uuid.UUID) (WorkflowRun, error)
}
