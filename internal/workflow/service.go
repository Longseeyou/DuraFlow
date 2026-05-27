package workflow

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowService struct {
	Repository WorkflowRepository
}

func (workflowService WorkflowService) CreateWorkflow(ctx context.Context, uId uuid.UUID, workflowName string, workflowDescription string) {
	w := Workflow{UserId: uId, Name: workflowName, Description: workflowDescription}
	workflowService.Repository.CreateWorkflow(ctx, w)
}

func (workflowService WorkflowService) GetWorkflowByUser(ctx context.Context, uId uuid.UUID) []Workflow {
	return workflowService.Repository.GetWorkflowByUser(ctx, uId)
}

func (workflowService WorkflowService) GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, w Workflow) Workflow {
	return workflowService.Repository.GetWorkflowByUserAndId(ctx, uId, w)
}

func (workflowService WorkflowService) UpdateWorkflowById(ctx context.Context, uId uuid.UUID, w Workflow) Workflow {
	return Workflow{}
}
