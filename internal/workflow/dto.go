package workflow

import (
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/dto"
)

type WorkflowRequestDto struct {
	Name        *string
	Description *string
}

type WorkflowResponseDto struct {
	dto.ResponseDto
	Name        string
	Description string
}

type WorkflowDefinitionRequestDto struct {
	Version *uint
	Status  *WorkflowDefinitionStatus
}

type WorkflowDefinitionResponseDto struct {
	dto.ResponseDto
	Version uint
	Status  WorkflowDefinitionStatus
}

type WorkflowRunRequestDto struct {
	Status *WorkflowRunStatus
}

type WorkflowRunResponseDto struct {
	dto.ResponseDto
	Status    WorkflowRunStatus
	StartedAt time.Time
	EndedAt   time.Time
}
