package workflow

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowRepository interface {
	CreateWorkflow(ctx context.Context, w Workflow) Workflow
	GetWorkflowByUser(ctx context.Context, uId uuid.UUID) []Workflow
	GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, w Workflow) Workflow
	UpdateWorkflowById(ctx context.Context, w Workflow) Workflow
}
