package task

import (
	"context"
	"errors"
	"fmt"

	"github.com/Longseeyou/DuraFlow/internal/shared/custom_error"
	"github.com/Longseeyou/DuraFlow/internal/shared/mapper"
	"github.com/google/uuid"
)

var (
	ErrWorkflowDefinitionNotEditing = errors.New("workflow definition is not in editing status")
)

type TaskService struct {
	repository TaskRepository
}

func NewTaskService(
	repository TaskRepository,
) TaskService {
	return TaskService{
		repository: repository,
	}
}

// TaskDefinition

func (s *TaskService) CreateTaskDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
	dto TaskDefinitionRequestDto,
) (TaskDefinitionResponseDto, error) {
	if dto.Name == nil || dto.TaskType == nil || dto.Timeout == nil {
		return TaskDefinitionResponseDto{}, custom_error.ErrCreateInvalidRequest
	}

	tD := TaskDefinition{
		WorkflowDefinitionID: workflowDefinitionID,
		Name:                 *dto.Name,
		Description:          mapper.StringValue(dto.Description),
		TaskType:             *dto.TaskType,
		Timeout:              *dto.Timeout,
	}

	tD, err := s.repository.CreateTaskDefinition(ctx, userID, workflowDefinitionID, tD)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (s *TaskService) GetTaskDefinitionsByWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) ([]TaskDefinitionResponseDto, error) {
	tDs, err := s.repository.GetTaskDefinitionsByWorkflowDefinition(
		ctx,
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		return nil, err
	}

	tDDtos := make([]TaskDefinitionResponseDto, len(tDs))
	for i, tD := range tDs {
		tDDtos[i] = taskDefinitionModelToResponseDto(tD)
	}

	return tDDtos, nil
}

func (s *TaskService) getTaskDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskDefinitionID uuid.UUID,
) (TaskDefinition, error) {
	tD, err := s.repository.GetTaskDefinitionByUserAndID(ctx, userID, taskDefinitionID)
	if err != nil {
		return TaskDefinition{}, err
	}

	return tD, nil
}

func (s *TaskService) GetTaskDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskDefinitionID uuid.UUID,
) (TaskDefinitionResponseDto, error) {
	tD, err := s.getTaskDefinitionByUserAndID(ctx, userID, taskDefinitionID)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (s *TaskService) UpdateTaskDefinitionByID(
	ctx context.Context,
	userID uuid.UUID,
	taskDefinitionID uuid.UUID,
	dto TaskDefinitionRequestDto,
) (TaskDefinitionResponseDto, error) {
	newTaskDefinition := map[string]any{}
	if dto.Name != nil {
		newTaskDefinition["name"] = *dto.Name
	}
	if dto.Description != nil {
		newTaskDefinition["description"] = *dto.Description
	}

	if len(newTaskDefinition) == 0 {
		return TaskDefinitionResponseDto{}, custom_error.ErrUpdateInvalidRequest
	}

	tD, err := s.repository.UpdateTaskDefinitionByUserAndID(
		ctx,
		userID,
		taskDefinitionID,
		newTaskDefinition,
	)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (s *TaskService) SoftDeleteTaskDefinition(
	ctx context.Context,
	userID uuid.UUID,
	taskDefinitionID uuid.UUID,
) (TaskDefinitionResponseDto, error) {
	tD, err := s.repository.SoftDeleteTaskDefinition(ctx, userID, taskDefinitionID)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

// TaskDependency

func (s *TaskService) isTaskDependencyValid(
	ctx context.Context,
	userID uuid.UUID,
	taskID uuid.UUID,
	dependOnTaskID uuid.UUID,
) (bool, error) {
	t, err := s.getTaskDefinitionByUserAndID(ctx, userID, taskID)
	if err != nil {
		return false, err
	}

	dependOnT, err := s.getTaskDefinitionByUserAndID(ctx, userID, dependOnTaskID)
	if err != nil {
		return false, err
	}

	return t.WorkflowDefinitionID == dependOnT.WorkflowDefinitionID, nil
}

func (s *TaskService) CreateTaskDependency(
	ctx context.Context,
	userID uuid.UUID,
	dto TaskDependencyRequestDto,
) (TaskDependencyResponseDto, error) {
	ok, err := s.isTaskDependencyValid(ctx, userID, dto.TaskID, dto.DependOnTaskID)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	if !ok {
		return TaskDependencyResponseDto{}, fmt.Errorf("")
	}

	tDp := TaskDependency{TaskID: dto.TaskID, DependOnTaskID: dto.DependOnTaskID}

	tDp, err = s.repository.CreateTaskDependency(ctx, tDp)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

func (s *TaskService) GetTaskDependencyByWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) ([]TaskDependencyResponseDto, error) {
	tDps, err := s.repository.GetTaskDependenciesByWorkflowDefinition(
		ctx,
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		return nil, err
	}

	tDpDtos := make([]TaskDependencyResponseDto, len(tDps))
	for i, tDp := range tDps {
		tDpDtos[i] = taskDependencyModelToResponseDto(tDp)
	}

	return tDpDtos, nil
}

func (s *TaskService) GetTaskDependencyByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskDependencyID uuid.UUID,
) (TaskDependencyResponseDto, error) {
	tDp, err := s.repository.GetTaskDependencyByUserAndID(ctx, userID, taskDependencyID)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

func (s *TaskService) UpdateTaskDependencyByID(
	ctx context.Context,
	userID uuid.UUID,
	taskDependencyID uuid.UUID,
	dto TaskDependencyRequestDto,
) (TaskDependencyResponseDto, error) {
	ok, err := s.isTaskDependencyValid(ctx, userID, dto.TaskID, dto.DependOnTaskID)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	if !ok {
		return TaskDependencyResponseDto{}, fmt.Errorf("")
	}

	newTaskDependency := map[string]any{}
	newTaskDependency["task_id"] = dto.TaskID
	newTaskDependency["depend_on_task_id"] = dto.DependOnTaskID

	tDp, err := s.repository.UpdateTaskDependencyByUserAndID(
		ctx,
		userID,
		taskDependencyID,
		newTaskDependency,
	)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	return taskDependencyModelToResponseDto(tDp), nil
}

func (s *TaskService) SoftDeleteTaskDependency(
	ctx context.Context,
	userID uuid.UUID,
	taskDependencyID uuid.UUID,
) (TaskDependencyResponseDto, error) {
	tDp, err := s.repository.SoftDeleteTaskDependency(ctx, userID, taskDependencyID)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

func (s *TaskService) IsTaskDependencyDag(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) (bool, error) {
	tDps, err := s.repository.GetTaskDependenciesByWorkflowDefinition(
		ctx,
		userID,
		workflowDefinitionID,
	)
	if err != nil {
		return false, err
	}

	graph := make(map[uuid.UUID][]uuid.UUID)

	for _, tDp := range tDps {
		graph[tDp.DependOnTaskID] = append(graph[tDp.DependOnTaskID], tDp.TaskID)
	}

	// https://en.wikipedia.org/wiki/Topological_sorting
	const (
		unvisited = iota
		visiting
		visited
	)

	mark := make(map[uuid.UUID]int, len(graph))

	var visit func(uuid.UUID) error
	visit = func(taskDefinitionID uuid.UUID) error {
		switch mark[taskDefinitionID] {
		case visiting:
			return fmt.Errorf("graph has at least one cycle")
		case visited:
			return nil
		}

		mark[taskDefinitionID] = visiting

		for _, nextTaskDefinitionID := range graph[taskDefinitionID] {
			err := visit(nextTaskDefinitionID)
			if err != nil {
				return err
			}
		}

		mark[taskDefinitionID] = visited
		return nil
	}

	for taskDefinitionID := range graph {
		err = visit(taskDefinitionID)
		if err != nil {
			return false, nil
		}
	}

	return true, nil
}

// TaskRun

func (s *TaskService) GetTaskRunByWorkflowRun(
	ctx context.Context,
	userID uuid.UUID,
	workflowRunID uuid.UUID,
) ([]TaskRunResponseDto, error) {
	tRs, err := s.repository.GetTaskRunsByWorkflowRun(ctx, userID, workflowRunID)
	if err != nil {
		return nil, err
	}

	tRDtos := make([]TaskRunResponseDto, len(tRs))
	for i, tR := range tRs {
		tRDtos[i] = taskRunModelToResponseDto(tR)
	}

	return tRDtos, nil
}

func (s *TaskService) GetTaskRunByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskRunID uuid.UUID,
) (TaskRunResponseDto, error) {
	tR, err := s.repository.GetTaskRunByUserAndID(ctx, userID, taskRunID)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

func (s *TaskService) UpdateTaskRunByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskRunID uuid.UUID,
	dto TaskRunRequestDto,
) (TaskRunResponseDto, error) {
	newTaskRun := map[string]any{}
	if dto.MaxRetries != nil {
		newTaskRun["max_retries"] = *dto.MaxRetries
	}
	if dto.Input != nil {
		newTaskRun["input"] = *dto.Input
	}

	if len(newTaskRun) == 0 {
		return TaskRunResponseDto{}, custom_error.ErrUpdateInvalidRequest
	}

	tR, err := s.repository.UpdateTaskRunByUserAndID(ctx, userID, taskRunID, newTaskRun)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

func (s *TaskService) SoftDeleteTaskRun(
	ctx context.Context,
	userID uuid.UUID,
	taskRunID uuid.UUID,
) (TaskRunResponseDto, error) {
	tR, err := s.repository.SoftDeleteTaskRun(ctx, userID, taskRunID)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

// TaskAttempt

func (s *TaskService) GetTaskAttemptByTaskRun(
	ctx context.Context,
	userID uuid.UUID,
	taskRunID uuid.UUID,
) ([]TaskAttemptResponseDto, error) {
	tAs, err := s.repository.GetTaskAttemptsByTaskRun(ctx, userID, taskRunID)
	if err != nil {
		return nil, err
	}

	tADtos := make([]TaskAttemptResponseDto, len(tAs))
	for i, tA := range tAs {
		tADtos[i] = taskAttemptModelToResponseDto(tA)
	}

	return tADtos, nil
}

func (s *TaskService) GetTaskAttemptByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskAttemptID uuid.UUID,
) (TaskAttemptResponseDto, error) {
	tA, err := s.repository.GetTaskAttemptByUserAndID(ctx, userID, taskAttemptID)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}

	return taskAttemptModelToResponseDto(tA), nil
}

func (s *TaskService) SoftDeleteTaskAttempt(
	ctx context.Context,
	userID uuid.UUID,
	taskAttemptID uuid.UUID,
) (TaskAttemptResponseDto, error) {
	tA, err := s.repository.SoftDeleteTaskAttempt(ctx, userID, taskAttemptID)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}

	return taskAttemptModelToResponseDto(tA), nil
}

// Internal

type TaskServiceInternal struct {
	repository TaskRepositoryInternal
}

func NewTaskServiceInternal(
	taskRepositoryInternal TaskRepositoryInternal,
) TaskServiceInternal {
	return TaskServiceInternal{
		repository: taskRepositoryInternal,
	}
}
