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

type TaskRunStatus string

const (
	TASK_RUN_PENDING       TaskRunStatus = "TASK_RUN_PENDING"
	TASK_RUN_QUEUED        TaskRunStatus = "TASK_RUN_QUEUED"
	TASK_RUN_RUNNING       TaskRunStatus = "TASK_RUN_RUNNING"
	TASK_RUN_COMPLETED     TaskRunStatus = "TASK_RUN_COMPLETED"
	TASK_RUN_FAILED        TaskRunStatus = "TASK_RUN_FAILED"
	TASK_RUN_CANCELLED     TaskRunStatus = "TASK_RUN_CANCELLED"
	TASK_RUN_DEAD_LETTERED TaskRunStatus = "TASK_RUN_DEAD_LETTERED"
)

type TaskAttemptStatus string

const (
	TASK_ATTEMPT_QUEUED        TaskAttemptStatus = "TASK_ATTEMPT_QUEUED"
	TASK_ATTEMPT_RUNNING       TaskAttemptStatus = "TASK_ATTEMPT_RUNNING"
	TASK_ATTEMPT_COMPLETED     TaskAttemptStatus = "TASK_ATTEMPT_COMPLETED"
	TASK_ATTEMPT_FAILED        TaskAttemptStatus = "TASK_ATTEMPT_FAILED"
	TASK_ATTEMPT_CANCELLED     TaskAttemptStatus = "TASK_ATTEMPT_CANCELLED"
	TASK_ATTEMPT_DEAD_LETTERED TaskAttemptStatus = "TASK_ATTEMPT_DEAD_LETTERED"
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
	EndedAt          *time.Time
	Input            *string
	Output           *string
}

type TaskAttempt struct {
	model.BaseModel
	TaskRunID     uuid.UUID `gorm:"uniqueIndex:idx_task_attempt"`
	TaskRun       TaskRun   `json:"-"`
	AttemptNumber uint      `gorm:"uniqueIndex:idx_task_attempt"`
	WorkerID      string
	Status        TaskAttemptStatus
	StartedAt     *time.Time
	EndedAt       *time.Time
	Log           *string
}
