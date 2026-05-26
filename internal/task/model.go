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
	WorkflowDefinitionId uuid.UUID
	WorkflowDefinition   workflow.WorkflowDefinition
}

type TaskRun struct {
	model.BaseModel
	WorkflowRunId    uuid.UUID
	WorkflowRun      workflow.WorkflowRun
	TaskDefinitionId uuid.UUID
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
	TaskRunId     uuid.UUID
	TaskRun       TaskRun
	AttemptNumber uint
	WorkerId      string
	Status        TaskRunStatus
	StartedAt     *time.Time
	CompletedAt   *time.Time
	log           *string
}

type TaskEvent struct {
	model.BaseModel
	WorkflowRunId uuid.UUID
	WorkflowRun   workflow.WorkflowRun
}
