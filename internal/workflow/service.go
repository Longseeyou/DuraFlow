package workflow

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowService struct {
	Repository         WorkflowRepository
	repositoryInternal WorkflowRepositoryInternal
}

func (workflowService WorkflowService) CreateWorkflow(ctx context.Context, uId uuid.UUID, wDto WorkflowRequestDto) (WorkflowResponseDto, error) {
	w := workflowRequestDtoToModel(wDto)
	w, error := workflowService.Repository.CreateWorkflow(ctx, w)
	if error != nil {

	}
	return workflowModelToResponseDto(w), nil
}

func (workflowService WorkflowService) GetWorkflowByUser(ctx context.Context, uId uuid.UUID) ([]WorkflowResponseDto, error) {
	w, error := workflowService.Repository.GetWorkflowByUser(ctx, uId)
	if error != nil {

	}
	wDto := []WorkflowResponseDto{}
	for _, v := range w {
		wDto = append(wDto, workflowModelToResponseDto(v))
	}
	return wDto, nil
}

func (workflowService WorkflowService) GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (WorkflowResponseDto, error) {
	w, error := workflowService.Repository.GetWorkflowByUserAndId(ctx, uId, wId)
	if error != nil {

	}
	return workflowModelToResponseDto(w), nil
}

func (workflowService WorkflowService) UpdateWorkflowById(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDto WorkflowRequestDto) (WorkflowResponseDto, error) {
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
	return workflowModelToResponseDto(w), nil
}

func (workflowService WorkflowService) SoftDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (WorkflowResponseDto, error) {
	w, error := workflowService.Repository.SoftDeleteWorkflow(ctx, uId, wId)
	if error != nil {

	}
	return workflowModelToResponseDto(w), nil
}

func (workflowService WorkflowService) HardDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (WorkflowResponseDto, error) {
	w, error := workflowService.Repository.HardDeleteWorkflow(ctx, uId, wId)
	if error != nil {

	}
	return workflowModelToResponseDto(w), nil
}

// WorkflowDefinition

func (workflowService WorkflowService) isUserAndWorkflowExists(ctx context.Context, uId uuid.UUID, wId uuid.UUID) bool {
	_, error := workflowService.Repository.GetWorkflowByUserAndId(ctx, uId, wId)
	return error != nil
}

func (workflowService WorkflowService) CreateWorkflowDefinition(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (WorkflowDefinitionResponseDto, error) {
	if !workflowService.isUserAndWorkflowExists(ctx, uId, wId) {

	}

	wD := WorkflowDefinition{WorkflowID: wId}
	wD, error := workflowService.Repository.CreateWorkflowDefinition(ctx, wD)
	if error != nil {

	}
	return workflowDefinitionModelToResponseDto(wD), nil
}

func (workflowService WorkflowService) getWorkflowDefinitionById(ctx context.Context, wDId uuid.UUID) (WorkflowDefinition, error) {
	wD, error := workflowService.repositoryInternal.GetWorkflowDefinitionById(ctx, wDId)
	if error != nil {

	}
	return wD, nil
}

func (workflowService WorkflowService) GetWorkflowDefinitionByUserAndWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) ([]WorkflowDefinitionResponseDto, error) {
	if !workflowService.isUserAndWorkflowExists(ctx, uId, wId) {

	}

	wD, error := workflowService.Repository.GetWorkflowDefinitionByWorkflow(ctx, wId)
	if error != nil {

	}
	wDDto := []WorkflowDefinitionResponseDto{}
	for _, v := range wD {
		wDDto = append(wDDto, workflowDefinitionModelToResponseDto(v))
	}
	return wDDto, nil
}

func (workflowService WorkflowService) GetWorkflowDefinitionByUserAndWorkflowAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDId uuid.UUID) (WorkflowDefinitionResponseDto, error) {
	if !workflowService.isUserAndWorkflowExists(ctx, uId, wId) {

	}

	wD, error := workflowService.Repository.GetWorkflowDefinitionByWorkflowAndId(ctx, wId, wDId)
	if error != nil {

	}
	return workflowDefinitionModelToResponseDto(wD), nil
}

func (workflowService WorkflowService) UpdateWorkflowDefinitionById(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDId uuid.UUID, wDDto WorkflowDefinitionRequestDto) (WorkflowDefinitionResponseDto, error) {
	if !workflowService.isUserAndWorkflowExists(ctx, uId, wId) {

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
	return workflowDefinitionModelToResponseDto(wD), nil
}

func (workflowService WorkflowService) SoftDeleteWorkflowDefinition(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDId uuid.UUID) (WorkflowDefinitionResponseDto, error) {
	if !workflowService.isUserAndWorkflowExists(ctx, uId, wId) {

	}

	wD, error := workflowService.Repository.SoftDeleteWorkflowDefinition(ctx, wId, wDId)
	if error != nil {

	}
	return workflowDefinitionModelToResponseDto(wD), nil
}

func (workflowService WorkflowService) HardDeleteWorkflowDefinition(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDId uuid.UUID) (WorkflowDefinitionResponseDto, error) {
	if !workflowService.isUserAndWorkflowExists(ctx, uId, wId) {

	}

	wD, error := workflowService.Repository.HardDeleteWorkflowDefinition(ctx, wId, wDId)
	if error != nil {

	}
	return workflowDefinitionModelToResponseDto(wD), nil
}

// WorkflowRun

func (workflowService WorkflowService) IsUserAndWorkflowDefinitionExists(ctx context.Context, uId uuid.UUID, wDId uuid.UUID) bool {
	wD, error := workflowService.getWorkflowDefinitionById(ctx, wDId)
	if error != nil {
		return false
	}

	_, error = workflowService.GetWorkflowByUserAndId(ctx, uId, wD.WorkflowID)
	return error != nil
}

func (workflowService WorkflowService) CreateWorkflowRun(ctx context.Context, uId uuid.UUID, wDId uuid.UUID) (WorkflowRunResponseDto, error) {
	if workflowService.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	wR := WorkflowRun{WorkflowDefinitionID: wDId, Status: PENDING}
	wR, error := workflowService.Repository.CreateWorkflowRun(ctx, wR)
	if error != nil {

	}
	return workflowRunModelToResponseDto(wR), nil
}

func (workflowService WorkflowService) GetWorkflowRunByWorkflowDefinition(ctx context.Context, uId uuid.UUID, wDId uuid.UUID) ([]WorkflowRunResponseDto, error) {
	if workflowService.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	wR, error := workflowService.Repository.GetWorkflowRunByWorkflowDefinition(ctx, wDId)
	if error != nil {

	}
	wRDto := []WorkflowRunResponseDto{}
	for _, v := range wR {
		wRDto = append(wRDto, workflowRunModelToResponseDto(v))
	}
	return wRDto, nil
}

func (workflowService WorkflowService) GetWorkflowRunByWorkflowDefinitionAndId(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, wRId uuid.UUID) (WorkflowRunResponseDto, error) {
	if workflowService.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	wR, error := workflowService.Repository.GetWorkflowRunByWorkflowDefinitionAndId(ctx, wDId, wRId)
	if error != nil {

	}
	return workflowRunModelToResponseDto(wR), nil
}

func (workflowService WorkflowService) UpdateWorkflowRunById(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, wRId uuid.UUID, wRDto WorkflowRunRequestDto) (WorkflowRunResponseDto, error) {
	if workflowService.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	newWR := map[string]any{}
	if wRDto.Status != nil {
		newWR["status"] = *wRDto.Status
	}
	wR, error := workflowService.Repository.UpdateWorkflowRunById(ctx, wDId, wRId, newWR)
	if error != nil {

	}
	return workflowRunModelToResponseDto(wR), nil
}

func (workflowService WorkflowService) SoftDeleteWorkflowRun(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, wRId uuid.UUID) (WorkflowRunResponseDto, error) {
	if workflowService.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	wR, error := workflowService.Repository.SoftDeleteWorkflowRun(ctx, wDId, wRId)
	if error != nil {

	}
	return workflowRunModelToResponseDto(wR), nil
}

func (workflowService WorkflowService) HardDeleteWorkflowRun(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, wRId uuid.UUID) (WorkflowRunResponseDto, error) {
	if workflowService.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	wR, error := workflowService.Repository.HardDeleteWorkflowRun(ctx, wDId, wRId)
	if error != nil {

	}
	return workflowRunModelToResponseDto(wR), nil
}
