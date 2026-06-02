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
}
