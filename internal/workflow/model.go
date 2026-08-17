package workflow

import (
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/model"
	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/google/uuid"
)

type Workflow struct {
	model.BaseModel
	UserID      uuid.UUID
	User        user.User `json:"-"`
	Name        string
	Description string
}

type WorkflowDefinitionStatus string

const (
	WORKFLOW_DEFINITION_EDITING   WorkflowDefinitionStatus = "WORKFLOW_DEFINITION_EDITING"
	WORKFLOW_DEFINITION_ACTIVATED WorkflowDefinitionStatus = "WORKFLOW_DEFINITION_ACTIVATED"
	WORKFLOW_DEFINITION_RUNNING   WorkflowDefinitionStatus = "WORKFLOW_DEFINITION_RUNNING"
)

func ValidPreviousWorkflowDefinitionStatus(
	newStatus WorkflowDefinitionStatus,
) []WorkflowDefinitionStatus {
	switch newStatus {
	case WORKFLOW_DEFINITION_EDITING:
		return []WorkflowDefinitionStatus{newStatus, WORKFLOW_DEFINITION_ACTIVATED}
	case WORKFLOW_DEFINITION_ACTIVATED:
		return []WorkflowDefinitionStatus{
			newStatus,
			WORKFLOW_DEFINITION_EDITING,
			WORKFLOW_DEFINITION_RUNNING,
		}
	case WORKFLOW_DEFINITION_RUNNING:
		return []WorkflowDefinitionStatus{newStatus, WORKFLOW_DEFINITION_ACTIVATED}
	default:
		return []WorkflowDefinitionStatus{}
	}
}

type WorkflowDefinition struct {
	model.BaseModel
	WorkflowID uuid.UUID
	Workflow   Workflow `json:"-"`
	Version    uint     `         gorm:"uniqueIndex:idx_workflow_definition;autoIncrement"`
	Status     WorkflowDefinitionStatus
}

type WorkflowRunStatus string

const (
	WORKFLOW_RUN_PENDING   WorkflowRunStatus = "WORKFLOW_RUN_PENDING"
	WORKFLOW_RUN_RUNNING   WorkflowRunStatus = "WORKFLOW_RUN_RUNNING"
	WORKFLOW_RUN_COMPLETED WorkflowRunStatus = "WORKFLOW_RUN_COMPLETED"
	WORKFLOW_RUN_FAILED    WorkflowRunStatus = "WORKFLOW_RUN_FAILED"
	WORKFLOW_RUN_CANCELLED WorkflowRunStatus = "WORKFLOW_RUN_CANCELLED"
)

type WorkflowRun struct {
	model.BaseModel
	WorkflowDefinitionID uuid.UUID
	WorkflowDefinition   WorkflowDefinition `json:"-"`
	Status               WorkflowRunStatus
	StartedAt            *time.Time
	EndedAt              *time.Time
	CancelledAt          *time.Time
}
