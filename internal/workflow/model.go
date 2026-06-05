package workflow

import (
	"github.com/Longseeyou/DuraFlow/internal/shared/model"
	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/google/uuid"
	"time"
)

type WorkflowDefinitionStatus int

const (
	DRAFTED WorkflowDefinitionStatus = iota
	ACTIVATE
	ARCHIVED
)

type WorkflowRunStatus int

const (
	PENDING WorkflowRunStatus = iota
	RUNNING
	COMPLETED
	FAILED
	CANCELLED
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
	CompletedAt          *time.Time
	CancelledAt          *time.Time
}
