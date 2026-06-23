package task

import (
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/dto"
	"github.com/google/uuid"
)

// HTTP

type TaskDefinitionRequestDto struct {
	Name        *string
	Description *string
	TaskType    *TaskType
}

type TaskDefinitionResponseDto struct {
	dto.ResponseDto
	Name        string
	Description string
	TaskType    TaskType
}

type TaskDependencyRequestDto struct {
	TaskID         uuid.UUID
	DependOnTaskID uuid.UUID
}

type TaskDependencyResponseDto struct {
	dto.ResponseDto
	TaskID         uuid.UUID
	DependOnTaskID uuid.UUID
}

type TaskRunRequestDto struct {
	MaxRetries *uint
	Input      *string
}

type TaskRunResponseDto struct {
	dto.ResponseDto
	WorkflowRunID    uuid.UUID
	TaskDefinitionID uuid.UUID
	Status           TaskRunStatus
	RetryCount       uint
	MaxRetries       uint
	ScheduledAt      *time.Time
	StartedAt        *time.Time
	EndedAt          *time.Time
	Input            *string
	Output           *string
	Error            *string
}

type TaskAttemptResponseDto struct {
	dto.ResponseDto
	TaskRunID     uuid.UUID
	AttemptNumber uint
	WorkerID      string
	Status        TaskRunStatus
	StartedAt     *time.Time
	EndedAt       *time.Time
	log           *string
}

// Internal

type TaskRunInternalDto struct {
	Status      *TaskRunStatus
	RetryCount  *uint
	ScheduledAt *time.Time
	StartedAt   *time.Time
	EndedAt     *time.Time
	Input       *string
	Output      *string
	Error       *string
}

type TaskAttemptInternalDto struct {
	AttemptNumber *uint
	WorkerID      *string
	Status        *TaskRunStatus
	StartedAt     *time.Time
	EndedAt       *time.Time
	log           *string
}

type TaskCommandRequest struct {
	TaskType TaskType
	Input    *string
}

type TaskCommandResponse struct {
	Status    *TaskRunStatus
	StartedAt *time.Time
	EndedAt   *time.Time
	Output    *string
	Log       *string
}
