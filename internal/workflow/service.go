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

func (workflowService WorkflowService) GetWorkflowDefinitionByUserAndWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) []WorkflowDefinitionResponseDto {
	wD, error := workflowService.Repository.GetWorkflowDefinitionByUserAndWorkflow(ctx, uId, wId)
	if error != nil {

	}
	wDDto := []WorkflowDefinitionResponseDto{}
	for _, v := range wD {
		wDDto = append(wDDto, workflowDefinitionModelToResponseDto(v))
	}
	return wDDto
}

//
// func (workflowService WorkflowService) GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID) WorkflowResponseDto {
// 	w, error := workflowService.Repository.GetWorkflowByUserAndId(ctx, uId, wId)
// 	if error != nil {
//
// 	}
// 	return workflowModelToResponseDto(w)
// }
//
// func (workflowService WorkflowService) UpdateWorkflowById(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDto WorkflowRequestDto) WorkflowResponseDto {
// 	newW := map[string]any{}
// 	if wDto.Name != nil {
// 		newW["name"] = *wDto.Name
// 	}
// 	if wDto.Description != nil {
// 		newW["description"] = *wDto.Description
// 	}
// 	w, error := workflowService.Repository.UpdateWorkflowById(ctx, uId, wId, newW)
// 	if error != nil {
//
// 	}
// 	return workflowModelToResponseDto(w)
// }
//
// func (workflowService WorkflowService) SoftDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) WorkflowResponseDto {
// 	w, error := workflowService.Repository.SoftDeleteWorkflow(ctx, uId, wId)
// 	if error != nil {
//
// 	}
// 	return workflowModelToResponseDto(w)
// }
//
// func (workflowService WorkflowService) HardDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) WorkflowResponseDto {
// 	w, error := workflowService.Repository.HardDeleteWorkflow(ctx, uId, wId)
// 	if error != nil {
//
// 	}
// 	return workflowModelToResponseDto(w)
// }
