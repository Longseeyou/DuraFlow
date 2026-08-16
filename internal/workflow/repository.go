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
	GetWorkflowsByUser(ctx context.Context, userID uuid.UUID) ([]Workflow, error)
	GetWorkflowByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		workflowID uuid.UUID,
	) (Workflow, error)
	UpdateWorkflowByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		workflowID uuid.UUID,
		newWorkflow map[string]any,
	) (Workflow, error)
	SoftDeleteWorkflow(
		ctx context.Context,
		userID uuid.UUID,
		workflowID uuid.UUID,
	) (Workflow, error)
}

type WorkflowDefinitionStore interface {
	CreateWorkflowDefinition(ctx context.Context, wD WorkflowDefinition) (WorkflowDefinition, error)
	GetWorkflowDefinitionsByWorkflow(
		ctx context.Context,
		userID uuid.UUID,
		workflowID uuid.UUID,
	) ([]WorkflowDefinition, error)
	GetWorkflowDefinitionByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		workflowDefinitionID uuid.UUID,
	) (WorkflowDefinition, error)
	UpdateWorkflowDefinitionByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		workflowDefinitionID uuid.UUID,
		newWorkflowD map[string]any,
	) (WorkflowDefinition, error)
	SoftDeleteWorkflowDefinition(
		ctx context.Context,
		userID uuid.UUID,
		workflowDefinitionID uuid.UUID,
	) (WorkflowDefinition, error)
}

type WorkflowRunStore interface {
	GetWorkflowRunsByWorkflowDefinition(
		ctx context.Context,
		userID uuid.UUID,
		workflowDefinitionID uuid.UUID,
	) ([]WorkflowRun, error)
	GetWorkflowRunByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		workflowRunID uuid.UUID,
	) (WorkflowRun, error)
	UpdateWorkflowRunByUserAndID(
		ctx context.Context,
		userID uuid.UUID,
		workflowRunID uuid.UUID,
		newWorkflowRun map[string]any,
	) (WorkflowRun, error)
	SoftDeleteWorkflowRun(
		ctx context.Context,
		userID uuid.UUID,
		workflowRunID uuid.UUID,
	) (WorkflowRun, error)

	CountActiveWorkflowRunByUserAndWorkflowDefinition(
		ctx context.Context,
		userID uuid.UUID,
		workflowDefinitionID uuid.UUID,
	) (int64, error)
}

type WorkflowRepositoryInternal interface {
	WorkflowStoreInternal
	WorkflowDefinitionStoreInternal
	WorkflowRunStoreInternal
}

type WorkflowStoreInternal interface {
	GetWorkflowByID(ctx context.Context, workflowID uuid.UUID) (Workflow, error)
	HardDeleteWorkflow(ctx context.Context, workflowID uuid.UUID) (Workflow, error)
}

type WorkflowDefinitionStoreInternal interface {
	GetWorkflowDefinitionByID(
		ctx context.Context,
		workflowDefinitionID uuid.UUID,
	) (WorkflowDefinition, error)
	HardDeleteWorkflowDefinition(
		ctx context.Context,
		workflowDefinitionID uuid.UUID,
	) (WorkflowDefinition, error)
}

type WorkflowRunStoreInternal interface {
	CreateWorkflowRun(
		ctx context.Context,
		workflowRun WorkflowRun,
	) (WorkflowRun, error)
	GetWorkflowRunByID(
		ctx context.Context,
		workflowRunID uuid.UUID,
	) (WorkflowRun, error)
	UpdateWorkflowRunByID(
		ctx context.Context,
		workflowRunID uuid.UUID,
		newWorkflowRun map[string]any,
	) (WorkflowRun, error)
	HardDeleteWorkflowRun(
		ctx context.Context,
		workflowRunID uuid.UUID,
	) (WorkflowRun, error)
}
