package task

import (
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/model"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type TaskType int

const (
	MOCK_TASK TaskType = iota
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
	WorkflowDefinition   workflow.WorkflowDefinition `json:"-"`
	Name                 string
	Description          string
	TaskType             TaskType
}

type TaskDependency struct {
	model.BaseModel
	TaskID         uuid.UUID      `gorm:"uniqueIndex:idx_task_dependency"`
	Task           TaskDefinition `json:"-"`
	DependOnTaskID uuid.UUID      `gorm:"uniqueIndex:idx_task_dependency"`
	DependOnTask   TaskDefinition `json:"-"`
}

type TaskRun struct {
	model.BaseModel
	WorkflowRunID    uuid.UUID            `gorm:"uniqueIndex:idx_workflow_task_run"`
	WorkflowRun      workflow.WorkflowRun `json:"-"`
	TaskDefinitionID uuid.UUID            `gorm:"uniqueIndex:idx_workflow_task_run"`
	TaskDefinition   TaskDefinition       `json:"-"`
	Status           TaskRunStatus
	RetryCount       uint
	MaxRetries       uint
	ScheduledAt      *time.Time
	StartedAt        *time.Time
	CompletedAt      *time.Time
	Input            *string
	Output           *string
}

type TaskAttempt struct {
	model.BaseModel
	TaskRunID     uuid.UUID `gorm:"uniqueIndex:idx_task_attempt"`
	TaskRun       TaskRun   `json:"-"`
	AttemptNumber uint      `gorm:"uniqueIndex:idx_task_attempt"`
	WorkerID      string
	Status        TaskRunStatus
	StartedAt     *time.Time
	CompletedAt   *time.Time
	Log           *string
}
