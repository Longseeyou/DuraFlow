package workflow

import (
	"context"

	"github.com/google/uuid"
)

type WorkflowService struct {
	repository WorkflowRepository
}

func NewWorkflowService(
	workflowRepository WorkflowRepository,
) *WorkflowService {
	return &WorkflowService{
		repository: workflowRepository,
	}
}

func (wS *WorkflowService) CreateWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wDto WorkflowRequestDto,
) (WorkflowResponseDto, error) {
	w := workflowRequestDtoToModel(wDto)
	w.UserID = uId

	w, err := wS.repository.CreateWorkflow(ctx, w)
	if err != nil {
		return WorkflowResponseDto{}, err
	}

	return workflowModelToResponseDto(w), nil
}

func (wS *WorkflowService) GetWorkflowByUser(
	ctx context.Context,
	uId uuid.UUID,
) ([]WorkflowResponseDto, error) {
	w, err := wS.repository.GetWorkflowByUser(ctx, uId)
	if err != nil {
		return nil, err
	}

	wDto := []WorkflowResponseDto{}
	for _, v := range w {
		wDto = append(wDto, workflowModelToResponseDto(v))
	}

	return wDto, nil
}

func (wS *WorkflowService) GetWorkflowByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (WorkflowResponseDto, error) {
	w, err := wS.repository.GetWorkflowByUserAndId(ctx, uId, wId)
	if err != nil {
		return WorkflowResponseDto{}, err
	}

	return workflowModelToResponseDto(w), nil
}

func (wS *WorkflowService) UpdateWorkflowById(
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

	w, err := wS.repository.UpdateWorkflowById(ctx, uId, wId, newW)
	if err != nil {
		return WorkflowResponseDto{}, err
	}

	return workflowModelToResponseDto(w), nil
}

func (wS *WorkflowService) SoftDeleteWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (WorkflowResponseDto, error) {
	w, err := wS.repository.SoftDeleteWorkflow(ctx, uId, wId)
	if err != nil {
		return WorkflowResponseDto{}, err
	}

	return workflowModelToResponseDto(w), nil
}

func (wS *WorkflowService) HardDeleteWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (WorkflowResponseDto, error) {
	w, err := wS.repository.HardDeleteWorkflow(ctx, uId, wId)
	if err != nil {
		return WorkflowResponseDto{}, err
	}

	return workflowModelToResponseDto(w), nil
}

// WorkflowDefinition

func (wS *WorkflowService) CreateWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	_, err := wS.GetWorkflowDefinitionByUserAndId(ctx, uId, wId)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	wD := WorkflowDefinition{WorkflowID: wId}
	wD.WorkflowID = wId
	wD.Status = DRAFTED

	wD, err = wS.repository.CreateWorkflowDefinition(ctx, wD)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (wS *WorkflowService) GetWorkflowDefinitionByUserAndWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) ([]WorkflowDefinitionResponseDto, error) {
	wD, err := wS.repository.GetWorkflowDefinitionByWorkflow(ctx, uId, wId)
	if err != nil {
		return nil, err
	}

	wDDto := []WorkflowDefinitionResponseDto{}
	for _, v := range wD {
		wDDto = append(wDDto, workflowDefinitionModelToResponseDto(v))
	}

	return wDDto, nil
}

func (wS *WorkflowService) GetWorkflowDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	wD, err := wS.repository.GetWorkflowDefinitionByUserAndId(ctx, uId, wDId)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (wS *WorkflowService) UpdateWorkflowDefinitionById(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
	wDDto WorkflowDefinitionRequestDto,
) (WorkflowDefinitionResponseDto, error) {
	newWD := map[string]any{}
	if wDDto.Status != nil {
		newWD["status"] = *wDDto.Status
	}

	wD, err := wS.repository.UpdateWorkflowDefinitionById(ctx, uId, wDId, newWD)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (wS *WorkflowService) SoftDeleteWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	wD, err := wS.repository.SoftDeleteWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (wS *WorkflowService) HardDeleteWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	wD, err := wS.repository.HardDeleteWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

// WorkflowRun

func (wS *WorkflowService) CreateWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (WorkflowRunResponseDto, error) {
	_, err := wS.GetWorkflowDefinitionByUserAndId(ctx, uId, wDId)
	if err != nil {
		return WorkflowRunResponseDto{}, err
	}

	wR := WorkflowRun{WorkflowDefinitionID: wDId, Status: PENDING}
	wR.WorkflowDefinitionID = wDId

	wR, err = wS.repository.CreateWorkflowRun(ctx, wR)
	if err != nil {
		return WorkflowRunResponseDto{}, err
	}

	return workflowRunModelToResponseDto(wR), nil
}

func (wS *WorkflowService) GetWorkflowRunByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]WorkflowRunResponseDto, error) {
	wR, err := wS.repository.GetWorkflowRunByWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
		return nil, err
	}

	wRDto := []WorkflowRunResponseDto{}
	for _, v := range wR {
		wRDto = append(wRDto, workflowRunModelToResponseDto(v))
	}
	return wRDto, nil
}

func (wS *WorkflowService) GetWorkflowRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (WorkflowRunResponseDto, error) {
	wR, err := wS.repository.GetWorkflowRunByUserAndId(ctx, uId, wRId)
	if err != nil {
		return WorkflowRunResponseDto{}, err
	}

	return workflowRunModelToResponseDto(wR), nil
}

func (wS *WorkflowService) UpdateWorkflowRunById(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
	wRDto WorkflowRunRequestDto,
) (WorkflowRunResponseDto, error) {
	newWR := map[string]any{}
	if wRDto.Status != nil {
		newWR["status"] = *wRDto.Status
	}

	wR, err := wS.repository.UpdateWorkflowRunById(ctx, uId, wRId, newWR)
	if err != nil {
		return WorkflowRunResponseDto{}, err
	}

	return workflowRunModelToResponseDto(wR), nil
}

func (wS *WorkflowService) SoftDeleteWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (WorkflowRunResponseDto, error) {
	wR, err := wS.repository.SoftDeleteWorkflowRun(ctx, uId, wRId)
	if err != nil {
		return WorkflowRunResponseDto{}, err
	}

	return workflowRunModelToResponseDto(wR), nil
}

func (wS *WorkflowService) HardDeleteWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (WorkflowRunResponseDto, error) {
	wR, err := wS.repository.HardDeleteWorkflowRun(ctx, uId, wRId)
	if err != nil {
		return WorkflowRunResponseDto{}, err
	}

	return workflowRunModelToResponseDto(wR), nil
}
