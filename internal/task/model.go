package task

import (
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/model"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type TaskType string

const (
	MOCK_TASK TaskType = "MOCK_TASK"
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
	Timeout              time.Duration
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

func ValidPreviousTaskRunStatus(newStatus TaskRunStatus) []TaskRunStatus {
	switch newStatus {
	case TASK_RUN_PENDING:
		return []TaskRunStatus{TASK_RUN_FAILED}
	case TASK_RUN_QUEUED:
		return []TaskRunStatus{TASK_RUN_PENDING, TASK_RUN_FAILED}
	case TASK_RUN_RUNNING:
		return []TaskRunStatus{TASK_RUN_QUEUED}
	case TASK_RUN_COMPLETED:
		return []TaskRunStatus{TASK_RUN_RUNNING}
	case TASK_RUN_FAILED:
		return []TaskRunStatus{TASK_RUN_RUNNING}
	case TASK_RUN_CANCELLED:
		return []TaskRunStatus{
			TASK_RUN_PENDING,
			TASK_RUN_QUEUED,
			TASK_RUN_RUNNING,
			TASK_RUN_FAILED,
		}
	case TASK_RUN_DEAD_LETTERED:
		return []TaskRunStatus{TASK_RUN_RUNNING, TASK_RUN_FAILED}
	}
	return nil
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
	TimeoutAt     *time.Time
	Log           *string
}

func ValidPreviousTaskAttemptStatus(
	newStatus TaskAttemptStatus,
) []TaskAttemptStatus {
	validPrevious := []TaskAttemptStatus{}
	switch newStatus {
	case TASK_ATTEMPT_QUEUED:
	case TASK_ATTEMPT_RUNNING:
		validPrevious = append(validPrevious, TASK_ATTEMPT_QUEUED)
	case TASK_ATTEMPT_COMPLETED:
		validPrevious = append(validPrevious, TASK_ATTEMPT_RUNNING)
	case TASK_ATTEMPT_FAILED:
		validPrevious = append(validPrevious, TASK_ATTEMPT_RUNNING)
	case TASK_ATTEMPT_CANCELLED:
		validPrevious = append(validPrevious, TASK_ATTEMPT_QUEUED, TASK_ATTEMPT_RUNNING)
	case TASK_ATTEMPT_DEAD_LETTERED:
		validPrevious = append(validPrevious, TASK_ATTEMPT_RUNNING)
	default:
	}
	return validPrevious
}
