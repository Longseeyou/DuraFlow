package workflow

import (
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/model"
	"github.com/google/uuid"
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

type WorkflowDefinition struct {
	model.BaseModel
	Name   string
	Status WorkflowDefinitionStatus
}

type WorkflowRun struct {
	model.BaseModel
	WorkflowDefinitionId uuid.UUID
	WorkflowDefinition   WorkflowDefinition
	Status               WorkflowRunStatus
	StartedAt            time.Time
	CompletedAt          time.Time
	CancelledAt          time.Time
}
