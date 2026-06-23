package workflow

import (
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/model"
	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/google/uuid"
)

type WorkflowDefinitionStatus string

const (
	DRAFTED  WorkflowDefinitionStatus = "WORKFLOW_DEFINITION_DRAFTED"
	ACTIVATE WorkflowDefinitionStatus = "WORKFLOW_DEFINITION_ACTIVATE"
	ARCHIVED WorkflowDefinitionStatus = "WORKFLOW_DEFINITION_ARCHIVED"
)

type WorkflowRunStatus string

const (
	PENDING   WorkflowRunStatus = "WORKFLOW_RUN_PENDING"
	RUNNING   WorkflowRunStatus = "WORKFLOW_RUN_RUNNING"
	COMPLETED WorkflowRunStatus = "WORKFLOW_RUN_COMPLETED"
	FAILED    WorkflowRunStatus = "WORKFLOW_RUN_FAILED"
	// CANCELLED WorkflowRunStatus = "WORKFLOW_RUN_CANCELLED"
)

type Workflow struct {
	model.BaseModel
	UserID      uuid.UUID
	User        user.User
	Name        string
	Description string
}

type WorkflowDefinition struct {
	model.BaseModel
	WorkflowID uuid.UUID
	Workflow   Workflow
	Version    uint
	Status     WorkflowDefinitionStatus
}

type WorkflowRun struct {
	model.BaseModel
	WorkflowDefinitionID uuid.UUID
	WorkflowDefinition   WorkflowDefinition
	Status               WorkflowRunStatus
	StartedAt            *time.Time
	EndedAt              *time.Time
}
