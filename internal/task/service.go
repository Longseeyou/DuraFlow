package task

import (
	"context"
	"fmt"

	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type TaskService struct {
	Repository TaskRepository
	wS         workflow.WorkflowService
}

func NewTaskService(
	taskRepository TaskRepository,
	workflowService workflow.WorkflowService,
) TaskService {
	return TaskService{
		Repository: taskRepository,
		wS:         workflowService,
	}
}

// TaskDefinition

func (s TaskService) CreateTaskDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
	tDDto TaskDefinitionRequestDto,
) (TaskDefinitionResponseDto, error) {
	_, err := s.wS.GetWorkflowDefinitionByUserAndId(ctx, uId, wDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	tD := taskDefinitionRequestDtoToModel(tDDto)
	tD.WorkflowDefinitionID = wDId

	tD, err = s.Repository.CreateTaskDefinition(ctx, tD)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (s TaskService) GetTaskDefinitionByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]TaskDefinitionResponseDto, error) {
	tDs, err := s.Repository.GetTaskDefinitionByWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
		return nil, err
	}

	dtos := []TaskDefinitionResponseDto{}
	for _, v := range tDs {
		dtos = append(dtos, taskDefinitionModelToResponseDto(v))
	}

	return dtos, nil
}

func (s TaskService) getTaskDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (TaskDefinition, error) {
	tD, err := s.Repository.GetTaskDefinitionByUserAndId(ctx, uId, tDId)
	if err != nil {
		return TaskDefinition{}, err
	}

	return tD, nil
}

func (s TaskService) GetTaskDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (TaskDefinitionResponseDto, error) {
	tD, err := s.getTaskDefinitionByUserAndId(ctx, uId, tDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (s TaskService) UpdateTaskDefinitionById(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
	tDDto TaskDefinitionRequestDto,
) (TaskDefinitionResponseDto, error) {
	newTD := map[string]any{}
	if tDDto.Name != nil {
		newTD["name"] = *tDDto.Name
	}
	if tDDto.Description != nil {
		newTD["description"] = *tDDto.Description
	}

	tD, err := s.Repository.UpdateTaskDefinitionByUserAndId(ctx, uId, tDId, newTD)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (s TaskService) SoftDeleteTaskDefinition(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (TaskDefinitionResponseDto, error) {
	tD, err := s.Repository.SoftDeleteTaskDefinition(ctx, uId, tDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (s TaskService) HardDeleteTaskDefinition(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (TaskDefinitionResponseDto, error) {
	tD, err := s.Repository.HardDeleteTaskDefinition(ctx, uId, tDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

// TaskDependency

func (s TaskService) isTaskDependencyValid(
	ctx context.Context,
	uId uuid.UUID,
	tId uuid.UUID,
	dpTId uuid.UUID,
) (bool, error) {
	t, err := s.getTaskDefinitionByUserAndId(ctx, uId, tId)
	if err != nil {
		return false, err
	}

	dpT, err := s.getTaskDefinitionByUserAndId(ctx, uId, dpTId)
	if err != nil {
		return false, err
	}

	return t.WorkflowDefinitionID == dpT.WorkflowDefinitionID, nil
}

func (s TaskService) CreateTaskDependency(
	ctx context.Context,
	uId uuid.UUID,
	tDpDto TaskDependencyRequestDto,
) (TaskDependencyResponseDto, error) {
	valid, err := s.isTaskDependencyValid(ctx, uId, tDpDto.TaskID, tDpDto.DependOnTaskID)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	if !valid {
		return TaskDependencyResponseDto{}, fmt.Errorf("")
	}

	tDp := taskDependencyRequestDtoToModel(tDpDto)

	tDp, err = s.Repository.CreateTaskDependency(ctx, tDp)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

func (s TaskService) GetTaskDependencyByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]TaskDependencyResponseDto, error) {
	tDp, err := s.Repository.GetTaskDependencyByWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
		return nil, err
	}

	tDpDto := []TaskDependencyResponseDto{}
	for _, v := range tDp {
		tDpDto = append(tDpDto, taskDependencyModelToResponseDto(v))
	}

	return tDpDto, nil
}

func (s TaskService) GetTaskDependencyByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) (TaskDependencyResponseDto, error) {
	tDp, err := s.Repository.GetTaskDependencyByUserAndId(ctx, uId, tDpId)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

func (s TaskService) UpdateTaskDependencyById(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
	tDpDto TaskDependencyRequestDto,
) (TaskDependencyResponseDto, error) {
	valid, err := s.isTaskDependencyValid(ctx, uId, tDpDto.TaskID, tDpDto.DependOnTaskID)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	if !valid {
		return TaskDependencyResponseDto{}, fmt.Errorf("")
	}

	newTDp := map[string]any{}
	newTDp["task_id"] = tDpDto.TaskID
	newTDp["depend_on_task_id"] = tDpDto.DependOnTaskID

	tDp, err := s.Repository.UpdateTaskDependencyByUserAndId(ctx, uId, tDpId, newTDp)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	return taskDependencyModelToResponseDto(tDp), nil
}

func (s TaskService) SoftDeleteTaskDependency(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) (TaskDependencyResponseDto, error) {
	tDp, err := s.Repository.SoftDeleteTaskDependency(ctx, uId, tDpId)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

func (s TaskService) HardDeleteTaskDependency(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) (TaskDependencyResponseDto, error) {
	tDp, err := s.Repository.HardDeleteTaskDependency(ctx, uId, tDpId)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

// TaskRun

func (s TaskService) CreateTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
	tDId uuid.UUID,
	tRDto TaskRunRequestDto,
) (TaskRunResponseDto, error) {
	_, err := s.wS.GetWorkflowRunByUserAndId(ctx, uId, wRId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	_, err = s.GetTaskDefinitionByUserAndId(ctx, uId, tDId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	tR := taskRunRequestDtoToModel(tRDto)
	tR.WorkflowRunID = wRId
	tR.TaskDefinitionID = tDId

	tR, err = s.Repository.CreateTaskRun(ctx, tR)
	if err != nil {
		return TaskRunResponseDto{}, err
	}
	return taskRunModelToResponseDto(tR), nil
}

func (s TaskService) GetTaskRunByWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) ([]TaskRunResponseDto, error) {
	tRs, err := s.Repository.GetTaskRunByWorkflowRun(ctx, uId, wRId)
	if err != nil {
		return nil, err
	}

	dtos := []TaskRunResponseDto{}
	for _, v := range tRs {
		dtos = append(dtos, taskRunModelToResponseDto(v))
	}

	return dtos, nil
}

func (s TaskService) GetTaskRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (TaskRunResponseDto, error) {
	tR, err := s.Repository.GetTaskRunByUserAndId(ctx, uId, tRId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

func (s TaskService) UpdateTaskRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
	tRDto TaskRunRequestDto,
) (TaskRunResponseDto, error) {
	newTR := map[string]any{}
	if tRDto.MaxRetries != nil {
		newTR["max_retries"] = *tRDto.MaxRetries
	}
	if tRDto.Input != nil {
		newTR["input"] = *tRDto.Input
	}

	tR, err := s.Repository.UpdateTaskRunByUserAndId(ctx, uId, tRId, newTR)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

func (s TaskService) SoftDeleteTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (TaskRunResponseDto, error) {
	tR, err := s.Repository.SoftDeleteTaskRun(ctx, uId, tRId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

func (s TaskService) HardDeleteTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (TaskRunResponseDto, error) {
	tR, err := s.Repository.HardDeleteTaskRun(ctx, uId, tRId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

// TaskAttempt

func (s TaskService) GetTaskAttemptByTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) ([]TaskAttemptResponseDto, error) {
	tAs, err := s.Repository.GetTaskAttemptByTaskRun(ctx, uId, tRId)
	if err != nil {
		return nil, err
	}

	dtos := []TaskAttemptResponseDto{}
	for _, v := range tAs {
		dtos = append(dtos, taskAttemptModelToResponseDto(v))
	}

	return dtos, nil
}

func (s TaskService) GetTaskAttemptByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (TaskAttemptResponseDto, error) {
	tA, err := s.Repository.GetTaskAttemptByUserAndId(ctx, uId, tAId)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}

	return taskAttemptModelToResponseDto(tA), nil
}

func (s TaskService) SoftDeleteTaskAttempt(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (TaskAttemptResponseDto, error) {
	tA, err := s.Repository.SoftDeleteTaskAttempt(ctx, uId, tAId)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}

	return taskAttemptModelToResponseDto(tA), nil
}

func (s TaskService) HardDeleteTaskAttempt(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (TaskAttemptResponseDto, error) {
	tA, err := s.Repository.HardDeleteTaskAttempt(ctx, uId, tAId)
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

func (s TaskServiceInternal) UpdateTaskRunById(
	ctx context.Context,
	tRId uuid.UUID,
	tRDto TaskRunInternalDto,
) (TaskRun, error) {
	newTR, err := taskRunInternalDtoToMap(tRDto)
	if err != nil {
		return TaskRun{}, err
	}

	tR, err := s.repository.UpdateTaskRunById(ctx, tRId, newTR)
	if err != nil {
		return TaskRun{}, err
	}

	return tR, nil
}

func (s TaskServiceInternal) CreateTaskAttempt(
	ctx context.Context,
	tRId uuid.UUID,
	tADto TaskAttemptInternalDto,
) (TaskAttempt, error) {
	tA := taskAttemptInternalDtoToModel(tADto)
	tA.TaskRunID = tRId

	tA, err := s.repository.CreateTaskAttempt(ctx, tA)
	if err != nil {
		return TaskAttempt{}, err
	}

	return tA, nil
}

func (s TaskServiceInternal) UpdateTaskAttemptById(
	ctx context.Context,
	tAId uuid.UUID,
	tADto TaskAttemptInternalDto,
) (TaskAttemptResponseDto, error) {
	newTA := map[string]any{}
	if tADto.WorkerID != nil {
		newTA["worker_id"] = *tADto.WorkerID
	}
	if tADto.Status != nil {
		newTA["status"] = *tADto.Status
	}
	if tADto.StartedAt != nil {
		newTA["started_at"] = *tADto.StartedAt
	}
	if tADto.EndedAt != nil {
		newTA["ended_at"] = *tADto.EndedAt
	}
	if tADto.log != nil {
		newTA["log"] = tADto.log
	}

	tA, err := s.repository.UpdateTaskAttemptById(ctx, tAId, newTA)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}

	return taskAttemptModelToResponseDto(tA), nil
}
