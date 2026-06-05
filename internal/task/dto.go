package task

import (
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/dto"
	"github.com/google/uuid"
)

type TaskDefinitionRequestDto struct {
	Name        *string
	Description *string
}

type TaskDefinitionResponseDto struct {
	dto.ResponseDto
	Name        string
	Description string
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
	CompletedAt      *time.Time
}

type TaskAttemptRequestDto struct {
}

type TaskAttemptResponseDto struct {
	dto.ResponseDto
	TaskRunID     uuid.UUID
	AttemptNumber uint
	WorkerID      string
	Status        TaskRunStatus
	StartedAt     *time.Time
	CompletedAt   *time.Time
	log           *string
}
