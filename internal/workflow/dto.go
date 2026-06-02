package workflow

import (
	"github.com/Longseeyou/DuraFlow/internal/shared/dto"
	"github.com/google/uuid"
)

type WorkflowRequestDto struct {
	ID          *uuid.UUID
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
