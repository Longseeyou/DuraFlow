package workflow

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowService struct {
	Repository         WorkflowRepository
	repositoryInternal WorkflowRepositoryInternal
}

func (wS WorkflowService) CreateWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wDto WorkflowRequestDto,
) (WorkflowResponseDto, error) {
	w := workflowRequestDtoToModel(wDto)
	w, err := wS.Repository.CreateWorkflow(ctx, w)
	if err != nil {
	}
	return workflowModelToResponseDto(w), nil
}

func (wS WorkflowService) GetWorkflowByUser(
	ctx context.Context,
	uId uuid.UUID,
) ([]WorkflowResponseDto, error) {
	w, err := wS.Repository.GetWorkflowByUser(ctx, uId)
	if err != nil {
	}
	wDto := []WorkflowResponseDto{}
	for _, v := range w {
		wDto = append(wDto, workflowModelToResponseDto(v))
	}
	return wDto, nil
}

func (wS WorkflowService) GetWorkflowByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (WorkflowResponseDto, error) {
	w, err := wS.Repository.GetWorkflowByUserAndId(ctx, uId, wId)
	if err != nil {
	}
	return workflowModelToResponseDto(w), nil
}

func (wS WorkflowService) UpdateWorkflowById(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
	wDto WorkflowRequestDto,
) (WorkflowResponseDto, error) {
	newW := map[string]any{}
	if wDto.Name != nil {
		newW["name"] = *wDto.Name
	}
	if wDto.Description != nil {
		newW["description"] = *wDto.Description
	}
	w, err := wS.Repository.UpdateWorkflowById(ctx, uId, wId, newW)
	if err != nil {
	}
	return workflowModelToResponseDto(w), nil
}

func (wS WorkflowService) SoftDeleteWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (WorkflowResponseDto, error) {
	w, err := wS.Repository.SoftDeleteWorkflow(ctx, uId, wId)
	if err != nil {
	}
	return workflowModelToResponseDto(w), nil
}

func (wS WorkflowService) HardDeleteWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (WorkflowResponseDto, error) {
	w, err := wS.Repository.HardDeleteWorkflow(ctx, uId, wId)
	if err != nil {
	}
	return workflowModelToResponseDto(w), nil
}

// WorkflowDefinition

func (wS WorkflowService) CreateWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	_, err := wS.GetWorkflowDefinitionByUserAndId(ctx, uId, wId)
	if err != nil {
	}

	wD := WorkflowDefinition{WorkflowID: wId}
	wD, err = wS.Repository.CreateWorkflowDefinition(ctx, wD)
	if err != nil {
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (wS WorkflowService) GetWorkflowDefinitionByUserAndWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) ([]WorkflowDefinitionResponseDto, error) {
	wD, err := wS.Repository.GetWorkflowDefinitionByWorkflow(ctx, uId, wId)
	if err != nil {
	}
	wDDto := []WorkflowDefinitionResponseDto{}
	for _, v := range wD {
		wDDto = append(wDDto, workflowDefinitionModelToResponseDto(v))
	}
	return wDDto, nil
}

func (wS WorkflowService) GetWorkflowDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	wD, err := wS.Repository.GetWorkflowDefinitionByUserAndId(ctx, uId, wDId)
	if err != nil {
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (wS WorkflowService) GetWorkflowDefinitionByUserAndWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (WorkflowDefinition, error) {
	wR, err := wS.repositoryInternal.GetWorkflowRunById(ctx, wRId)
	if err != nil {
	}

	wD, err := wS.repositoryInternal.GetWorkflowDefinitionById(ctx, wR.WorkflowDefinitionID)
	if err != nil {
	}

	_, err = wS.Repository.GetWorkflowByUserAndId(ctx, uId, wD.WorkflowID)
	if err != nil {
	}
	return wD, nil
}

func (wS WorkflowService) UpdateWorkflowDefinitionById(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
	wDDto WorkflowDefinitionRequestDto,
) (WorkflowDefinitionResponseDto, error) {
	newWD := map[string]any{}
	if wDDto.Version != nil {
		newWD["version"] = *wDDto.Version
	}
	if wDDto.Status != nil {
		newWD["status"] = *wDDto.Status
	}
	wD, err := wS.Repository.UpdateWorkflowDefinitionById(ctx, uId, wDId, newWD)
	if err != nil {
	}
	return workflowDefinitionModelToResponseDto(wD), nil
}

func (wS WorkflowService) SoftDeleteWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	wD, err := wS.Repository.SoftDeleteWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
	}
	return workflowDefinitionModelToResponseDto(wD), nil
}

func (wS WorkflowService) HardDeleteWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	wD, err := wS.Repository.HardDeleteWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
	}
	return workflowDefinitionModelToResponseDto(wD), nil
}

// WorkflowRun

func (wS WorkflowService) CreateWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (WorkflowRunResponseDto, error) {
	_, err := wS.GetWorkflowDefinitionByUserAndId(ctx, uId, wDId)
	if err != nil {
	}

	wR := WorkflowRun{WorkflowDefinitionID: wDId, Status: PENDING}
	wR, err = wS.Repository.CreateWorkflowRun(ctx, wR)
	if err != nil {
	}

	return workflowRunModelToResponseDto(wR), nil
}

func (wS WorkflowService) GetWorkflowRunByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]WorkflowRunResponseDto, error) {
	wR, err := wS.Repository.GetWorkflowRunByWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
	}

	wRDto := []WorkflowRunResponseDto{}
	for _, v := range wR {
		wRDto = append(wRDto, workflowRunModelToResponseDto(v))
	}
	return wRDto, nil
}

func (wS WorkflowService) GetWorkflowRunByWorkflowDefinitionAndId(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (WorkflowRunResponseDto, error) {
	wR, err := wS.Repository.GetWorkflowRunByUserAndId(ctx, uId, wRId)
	if err != nil {
	}

	return workflowRunModelToResponseDto(wR), nil
}

func (wS WorkflowService) UpdateWorkflowRunById(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
	wRDto WorkflowRunRequestDto,
) (WorkflowRunResponseDto, error) {
	newWR := map[string]any{}
	if wRDto.Status != nil {
		newWR["status"] = *wRDto.Status
	}

	wR, err := wS.Repository.UpdateWorkflowRunById(ctx, uId, wRId, newWR)
	if err != nil {
	}

	return workflowRunModelToResponseDto(wR), nil
}

func (wS WorkflowService) SoftDeleteWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (WorkflowRunResponseDto, error) {
	wR, err := wS.Repository.SoftDeleteWorkflowRun(ctx, uId, wRId)
	if err != nil {
	}

	return workflowRunModelToResponseDto(wR), nil
}

func (wS WorkflowService) HardDeleteWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (WorkflowRunResponseDto, error) {
	wR, err := wS.Repository.HardDeleteWorkflowRun(ctx, uId, wRId)
	if err != nil {
	}

	return workflowRunModelToResponseDto(wR), nil
}
