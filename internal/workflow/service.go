package workflow

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowService struct {
	Repository WorkflowRepository
}

func (workflowService WorkflowService) CreateWorkflow(ctx context.Context, uId uuid.UUID, wDto WorkflowRequestDto) WorkflowResponseDto {
	w := workflowRequestDtoToWorkflowModel(wDto)
	w, error := workflowService.Repository.CreateWorkflow(ctx, w)
	if error != nil {

	}
	return workflowModelToResponseDto(w)
}

func (workflowService WorkflowService) GetWorkflowByUser(ctx context.Context, uId uuid.UUID) []WorkflowResponseDto {
	w, error := workflowService.Repository.GetWorkflowByUser(ctx, uId)
	if error != nil {

	}
	wDto := []WorkflowResponseDto{}
	for _, v := range w {
		wDto = append(wDto, workflowModelToResponseDto(v))
	}
	return wDto
}

func (workflowService WorkflowService) GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID) WorkflowResponseDto {
	w, error := workflowService.Repository.GetWorkflowByUserAndId(ctx, uId, wId)
	if error != nil {

	}
	return workflowModelToResponseDto(w)
}

func (workflowService WorkflowService) UpdateWorkflowById(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDto WorkflowRequestDto) WorkflowResponseDto {
	newW := map[string]any{}
	if wDto.Name != nil {
		newW["name"] = *wDto.Name
	}
	if wDto.Description != nil {
		newW["description"] = *wDto.Description
	}
	w, error := workflowService.Repository.UpdateWorkflowById(ctx, uId, wId, newW)
	if error != nil {

	}
	return workflowModelToResponseDto(w)
}

func (workflowService WorkflowService) SoftDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) WorkflowResponseDto {
	w, error := workflowService.Repository.SoftDeleteWorkflow(ctx, uId, wId)
	if error != nil {

	}
	return workflowModelToResponseDto(w)
}

func (workflowService WorkflowService) HardDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) WorkflowResponseDto {
	w, error := workflowService.Repository.HardDeleteWorkflow(ctx, uId, wId)
	if error != nil {

	}
	return workflowModelToResponseDto(w)
}

// WorkflowDefinition

func (workflowService WorkflowService) CreateWorkflowDefinition(ctx context.Context, uId uuid.UUID, wId uuid.UUID) WorkflowDefinitionResponseDto {
	wD := WorkflowDefinition{WorkflowId: wId}
	wD, error := workflowService.Repository.CreateWorkflowDefinition(ctx, wD)
	if error != nil {

	}
	return workflowDefinitionModelToResponseDto(wD)
}

func (workflowService WorkflowService) IsUserAndWorkflowExists(ctx context.Context, uId uuid.UUID, wId uuid.UUID) bool {
	_, error := workflowService.Repository.GetWorkflowByUserAndId(ctx, uId, wId)
	return error != nil
}

func (workflowService WorkflowService) GetWorkflowDefinitionByUserAndWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) []WorkflowDefinitionResponseDto {
	if !workflowService.IsUserAndWorkflowExists(ctx, uId, wId) {

	}

	wD, error := workflowService.Repository.GetWorkflowDefinitionByWorkflow(ctx, wId)
	if error != nil {

	}
	wDDto := []WorkflowDefinitionResponseDto{}
	for _, v := range wD {
		wDDto = append(wDDto, workflowDefinitionModelToResponseDto(v))
	}
	return wDDto
}

func (workflowService WorkflowService) GetWorkflowDefinitionByUserAndWorkflowAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDId uuid.UUID) WorkflowDefinitionResponseDto {
	if !workflowService.IsUserAndWorkflowExists(ctx, uId, wId) {

	}

	wD, error := workflowService.Repository.GetWorkflowDefinitionByWorkflowAndId(ctx, wId, wDId)
	if error != nil {

	}
	return workflowDefinitionModelToResponseDto(wD)
}

func (workflowService WorkflowService) UpdateWorkflowDefinitionById(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDId uuid.UUID, wDDto WorkflowDefinitionRequestDto) WorkflowDefinitionResponseDto {
	if !workflowService.IsUserAndWorkflowExists(ctx, uId, wId) {

	}

	newWD := map[string]any{}
	if wDDto.Version != nil {
		newWD["version"] = *wDDto.Version
	}
	if wDDto.Status != nil {
		newWD["status"] = *wDDto.Status
	}
	wD, error := workflowService.Repository.UpdateWorkflowDefinitionById(ctx, wId, wDId, newWD)
	if error != nil {

	}
	return workflowDefinitionModelToResponseDto(wD)
}

func (workflowService WorkflowService) SoftDeleteWorkflowDefinition(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDId uuid.UUID) WorkflowDefinitionResponseDto {
	if !workflowService.IsUserAndWorkflowExists(ctx, uId, wId) {

	}

	wD, error := workflowService.Repository.SoftDeleteWorkflowDefinition(ctx, wId, wDId)
	if error != nil {

	}
	return workflowDefinitionModelToResponseDto(wD)
}

func (workflowService WorkflowService) HardDeleteWorkflowDefinition(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDId uuid.UUID) WorkflowDefinitionResponseDto {
	if !workflowService.IsUserAndWorkflowExists(ctx, uId, wId) {

	}

	wD, error := workflowService.Repository.HardDeleteWorkflowDefinition(ctx, wId, wDId)
	if error != nil {

	}
	return workflowDefinitionModelToResponseDto(wD)
}
