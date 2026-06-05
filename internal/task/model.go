package task

import (
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/model"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type TaskRunStatus int

const (
	PENDING TaskRunStatus = iota
	QUEUED
	RUNNING
	COMPLETED
	FAILED
	RETRYING
	CANCELLED
	DEAD_LETTERED
)

type TaskDefinition struct {
	model.BaseModel
	WorkflowDefinitionID uuid.UUID
	WorkflowDefinition   workflow.WorkflowDefinition
	Name                 string
	Description          string
}

type TaskDependency struct {
	model.BaseModel
	TaskID         uuid.UUID
	Task           TaskDefinition
	DependOnTaskID uuid.UUID
	DependOnTask   TaskDefinition
}

type TaskRun struct {
	model.BaseModel
	WorkflowRunID    uuid.UUID
	WorkflowRun      workflow.WorkflowRun
	TaskDefinitionID uuid.UUID
	TaskDefinition   TaskDefinition
	Status           TaskRunStatus
	RetryCount       uint
	MaxRetries       uint
	ScheduledAt      *time.Time
	StartedAt        *time.Time
	CompletedAt      *time.Time
}

type TaskAttempt struct {
	model.BaseModel
	TaskRunID     uuid.UUID
	TaskRun       TaskRun
	AttemptNumber uint
	WorkerID      string
	Status        TaskRunStatus
	StartedAt     *time.Time
	CompletedAt   *time.Time
	log           *string
}

type TaskEvent struct {
	model.BaseModel
	WorkflowRunID uuid.UUID
	WorkflowRun   workflow.WorkflowRun
}
