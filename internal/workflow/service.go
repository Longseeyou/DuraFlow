package workflow

import (
	"context"
	"errors"

	"github.com/Longseeyou/DuraFlow/internal/shared/custom_error"
	"github.com/Longseeyou/DuraFlow/internal/shared/mapper"
	"github.com/google/uuid"
)

var (
	ErrWorkflowNameRequired     = errors.New("workflow name is required")
	ErrWorkflowDefinitionNotDAG = errors.New(
		"workflow definition task dependencies contain a cycle",
	)
	ErrWorkflowDefinitionHasActiveRuns = errors.New(
		"workflow definition has active workflow runs",
	)
)

type TaskService interface {
	IsTaskDependencyDag(
		ctx context.Context,
		userID uuid.UUID,
		workflowDefinitionID uuid.UUID,
	) (bool, error)
}

type WorkflowService struct {
	repository  WorkflowRepository
	taskService TaskService
}

func NewWorkflowService(
	repository WorkflowRepository,
	repositoryInternal WorkflowRepositoryInternal,
	taskService TaskService,
) *WorkflowService {
	return &WorkflowService{
		repository:  repository,
		taskService: taskService,
	}
}

func (s *WorkflowService) CreateWorkflow(
	ctx context.Context,
	userID uuid.UUID,
	dto WorkflowRequestDto,
) (WorkflowResponseDto, error) {
	if dto.Name == nil {
		return WorkflowResponseDto{}, ErrWorkflowNameRequired
	}

	w := Workflow{
		Name:        *dto.Name,
		Description: mapper.StringValue(dto.Description),
		UserID:      userID,
	}

	w, err := s.repository.CreateWorkflow(ctx, w)
	if err != nil {
		return WorkflowResponseDto{}, err
	}

	return workflowModelToResponseDto(w), nil
}

func (s *WorkflowService) GetWorkflowByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]WorkflowResponseDto, error) {
	ws, err := s.repository.GetWorkflowsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	wDtos := make([]WorkflowResponseDto, len(ws))
	for i, w := range ws {
		wDtos[i] = workflowModelToResponseDto(w)
	}

	return wDtos, nil
}

func (s *WorkflowService) GetWorkflowByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowID uuid.UUID,
) (WorkflowResponseDto, error) {
	w, err := s.repository.GetWorkflowByUserAndID(ctx, userID, workflowID)
	if err != nil {
		return WorkflowResponseDto{}, err
	}

	return workflowModelToResponseDto(w), nil
}

func (s *WorkflowService) UpdateWorkflowByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowID uuid.UUID,
	dto WorkflowRequestDto,
) (WorkflowResponseDto, error) {
	newWorkflow := map[string]any{}
	if dto.Name != nil {
		newWorkflow["name"] = *dto.Name
	}
	if dto.Description != nil {
		newWorkflow["description"] = *dto.Description
	}

	if len(newWorkflow) == 0 {
		return WorkflowResponseDto{}, custom_error.ErrUpdateInvalidRequest
	}

	w, err := s.repository.UpdateWorkflowByUserAndID(ctx, userID, workflowID, newWorkflow)
	if err != nil {
		return WorkflowResponseDto{}, err
	}

	return workflowModelToResponseDto(w), nil
}

func (s *WorkflowService) SoftDeleteWorkflow(
	ctx context.Context,
	userID uuid.UUID,
	workflowID uuid.UUID,
) (WorkflowResponseDto, error) {
	w, err := s.repository.SoftDeleteWorkflow(ctx, userID, workflowID)
	if err != nil {
		return WorkflowResponseDto{}, err
	}

	return workflowModelToResponseDto(w), nil
}

// WorkflowDefinition

func (s *WorkflowService) CreateWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowID uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	_, err := s.GetWorkflowDefinitionByUserAndID(ctx, userID, workflowID)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	wD := WorkflowDefinition{WorkflowID: workflowID, Status: WORKFLOW_DEFINITION_EDITING}

	wD, err = s.repository.CreateWorkflowDefinition(ctx, wD)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (s *WorkflowService) GetWorkflowDefinitionByUserAndWorkflow(
	ctx context.Context,
	userID uuid.UUID,
	workflowID uuid.UUID,
) ([]WorkflowDefinitionResponseDto, error) {
	wDs, err := s.repository.GetWorkflowDefinitionsByWorkflow(ctx, userID, workflowID)
	if err != nil {
		return nil, err
	}

	wDDtos := make([]WorkflowDefinitionResponseDto, len(wDs))
	for i, wD := range wDs {
		wDDtos[i] = workflowDefinitionModelToResponseDto(wD)
	}

	return wDDtos, nil
}

func (s *WorkflowService) GetWorkflowDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	wD, err := s.repository.GetWorkflowDefinitionByUserAndID(ctx, userID, workflowDefinitionID)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (s *WorkflowService) UpdateWorkflowDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
	dto WorkflowDefinitionRequestDto,
) (WorkflowDefinitionResponseDto, error) {
	newWorkflowDefinition := map[string]any{}
	if dto.Status != nil {
		newWorkflowDefinition["status"] = *dto.Status
		if newWorkflowDefinition["status"] == WORKFLOW_DEFINITION_ACTIVATED {
			ok, err := s.taskService.IsTaskDependencyDag(ctx, userID, workflowDefinitionID)
			if err != nil {
				return WorkflowDefinitionResponseDto{}, err
			}
			if !ok {
				return WorkflowDefinitionResponseDto{}, ErrWorkflowDefinitionNotDAG
			}
		} else {
			count, err := s.repository.CountActiveWorkflowRunByUserAndWorkflowDefinition(
				ctx,
				userID,
				workflowDefinitionID,
			)
			if err != nil {
				return WorkflowDefinitionResponseDto{}, err
			}
			if count != 0 {
				return WorkflowDefinitionResponseDto{}, ErrWorkflowDefinitionHasActiveRuns
			}

		}
	}

	if len(newWorkflowDefinition) == 0 {
		return WorkflowDefinitionResponseDto{}, custom_error.ErrUpdateInvalidRequest
	}

	wD, err := s.repository.UpdateWorkflowDefinitionByUserAndID(
		ctx,
		userID,
		workflowDefinitionID,
		newWorkflowDefinition,
	)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (s *WorkflowService) SoftDeleteWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) (WorkflowDefinitionResponseDto, error) {
	wD, err := s.repository.SoftDeleteWorkflowDefinition(ctx, userID, workflowDefinitionID)
	if err != nil {
		return WorkflowDefinitionResponseDto{}, err
	}

	return workflowDefinitionModelToResponseDto(wD), nil
}

func (s *WorkflowService) ExistsWorkflowDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) (bool, error) {
	_, err := s.repository.GetWorkflowDefinitionByUserAndID(ctx, userID, workflowDefinitionID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// WorkflowRun

func (s *WorkflowService) GetWorkflowRunByWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) ([]WorkflowRunResponseDto, error) {
	wRs, err := s.repository.GetWorkflowRunsByWorkflowDefinition(
		ctx,
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		return nil, err
	}

	wRDtos := make([]WorkflowRunResponseDto, len(wRs))
	for i, wR := range wRs {
		wRDtos[i] = workfloworkflowRununModelToResponseDto(wR)
	}
	return wRDtos, nil
}

func (s *WorkflowService) GetWorkflowRunByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowRunID uuid.UUID,
) (WorkflowRunResponseDto, error) {
	wR, err := s.repository.GetWorkflowRunByUserAndID(
		ctx,
		userID,
		workflowRunID,
	)
	if err != nil {
		return WorkflowRunResponseDto{}, err
	}

	return workfloworkflowRununModelToResponseDto(wR), nil
}

func (s *WorkflowService) UpdateWorkflowRunByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowRunID uuid.UUID,
	dto WorkflowRunRequestDto,
) (WorkflowRunResponseDto, error) {
	newWorkflowRun := map[string]any{}
	if dto.Status != nil {
		newWorkflowRun["status"] = *dto.Status
	}

	if len(newWorkflowRun) == 0 {
		return WorkflowRunResponseDto{}, custom_error.ErrUpdateInvalidRequest
	}

	workflowRun, err := s.repository.UpdateWorkflowRunByUserAndID(
		ctx,
		userID,
		workflowRunID,
		newWorkflowRun,
	)
	if err != nil {
		return WorkflowRunResponseDto{}, err
	}

	return workfloworkflowRununModelToResponseDto(workflowRun), nil
}

func (s *WorkflowService) SoftDeleteWorkflowRun(
	ctx context.Context,
	userID uuid.UUID,
	workflowRunID uuid.UUID,
) (WorkflowRunResponseDto, error) {
	workflowRun, err := s.repository.SoftDeleteWorkflowRun(
		ctx,
		userID,
		workflowRunID,
	)
	if err != nil {
		return WorkflowRunResponseDto{}, err
	}

	return workfloworkflowRununModelToResponseDto(workflowRun), nil
}
