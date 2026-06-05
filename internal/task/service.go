package task

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type TaskService struct {
	Repository         TaskRepository
	repositoryInternal TaskRepositoryInternal
	WS                 workflow.WorkflowService
}

// TaskDefinition

func (s TaskService) CreateTaskDefinition(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, tDDto TaskDefinitionRequestDto) (TaskDefinitionResponseDto, error) {
	if !s.WS.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}
	tD := taskDefinitionRequestDtoToModel(tDDto)
	tD.WorkflowDefinitionID = wDId
	tD, err := s.Repository.CreateTaskDefinition(ctx, tD)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}
	return taskDefinitionModelToResponseDto(tD), nil
}

func (s TaskService) GetTaskDefinitionByWorkflowDefinition(ctx context.Context, uId uuid.UUID, wDId uuid.UUID) ([]TaskDefinitionResponseDto, error) {
	if !s.WS.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	tDs, err := s.Repository.GetTaskDefinitionByWorkflowDefinition(ctx, wDId)
	if err != nil {
		return nil, err
	}
	dtos := []TaskDefinitionResponseDto{}
	for _, v := range tDs {
		dtos = append(dtos, taskDefinitionModelToResponseDto(v))
	}
	return dtos, nil
}

func (s TaskService) GetTaskDefinitionByWorkflowDefinitionAndId(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, tDId uuid.UUID) (TaskDefinitionResponseDto, error) {
	if !s.WS.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	tD, err := s.Repository.GetTaskDefinitionByWorkflowDefinitionAndId(ctx, wDId, tDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}
	return taskDefinitionModelToResponseDto(tD), nil
}

func (s TaskService) isUserAndTaskDefinitionExists(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, tDId uuid.UUID) bool {
	if !s.WS.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {
		return false
	}
	_, error := s.Repository.GetTaskDefinitionByWorkflowDefinitionAndId(ctx, wDId, tDId)
	return error != nil
}

func (s TaskService) UpdateTaskDefinitionById(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, tDId uuid.UUID, tDDto TaskDefinitionRequestDto) (TaskDefinitionResponseDto, error) {
	if !s.isUserAndTaskDefinitionExists(ctx, uId, wDId, tDId) {

	}

	newTD := map[string]any{}
	if tDDto.Name != nil {
		newTD["name"] = *tDDto.Name
	}
	if tDDto.Description != nil {
		newTD["description"] = *tDDto.Description
	}
	tD, err := s.Repository.UpdateTaskDefinitionById(ctx, wDId, tDId, newTD)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}
	return taskDefinitionModelToResponseDto(tD), nil
}

func (s TaskService) SoftDeleteTaskDefinition(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, tDId uuid.UUID) (TaskDefinitionResponseDto, error) {
	if !s.isUserAndTaskDefinitionExists(ctx, uId, wDId, tDId) {

	}

	tD, err := s.Repository.SoftDeleteTaskDefinition(ctx, wDId, tDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}
	return taskDefinitionModelToResponseDto(tD), nil
}

func (s TaskService) HardDeleteTaskDefinition(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, tDId uuid.UUID) (TaskDefinitionResponseDto, error) {
	if !s.isUserAndTaskDefinitionExists(ctx, uId, wDId, tDId) {

	}

	tD, err := s.Repository.HardDeleteTaskDefinition(ctx, wDId, tDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}
	return taskDefinitionModelToResponseDto(tD), nil
}

// TaskDependency

func (s TaskService) isTaskDependencyValid(ctx context.Context, uId uuid.UUID, tId uuid.UUID, dpTId uuid.UUID) bool {
	t, error := s.repositoryInternal.GetTaskDefinitionById(ctx, tId)
	if error != nil {
		return false
	}

	dpT, error := s.repositoryInternal.GetTaskDefinitionById(ctx, dpTId)
	if error != nil {
		return false
	}

	if t.WorkflowDefinitionID != dpT.WorkflowDefinitionID {
		return false
	}

	if !s.WS.IsUserAndWorkflowDefinitionExists(ctx, uId, t.WorkflowDefinitionID) {
		return false
	}

	return true
}

func (s TaskService) isUserAndTaskDepedencyExists(ctx context.Context, uId uuid.UUID, tDpId uuid.UUID) bool {
	tDp, err := s.repositoryInternal.GetTaskDependencyById(ctx, tDpId)
	if err != nil {
		return false
	}
	return s.isTaskDependencyValid(ctx, uId, tDp.TaskID, tDp.DependOnTaskID)
}

func (s TaskService) CreateTaskDependency(ctx context.Context, uId uuid.UUID, tDpDto TaskDependencyRequestDto) (TaskDependencyResponseDto, error) {
	if !s.isTaskDependencyValid(ctx, uId, tDpDto.TaskID, tDpDto.DependOnTaskID) {

	}

	tDp := taskDependencyRequestDtoToModel(tDpDto)
	tDp, err := s.Repository.CreateTaskDependency(ctx, tDp)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	return taskDependencyModelToResponseDto(tDp), nil
}

func (s TaskService) GetTaskDependencyByWorkflowDefinition(ctx context.Context, uId uuid.UUID, wDId uuid.UUID) ([]TaskDependencyResponseDto, error) {
	if !s.WS.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	tDp, err := s.Repository.GetTaskDependencyByWorkflowDefinition(ctx, wDId)
	if err != nil {
		return nil, err
	}
	tDpDto := []TaskDependencyResponseDto{}
	for _, v := range tDp {
		tDpDto = append(tDpDto, taskDependencyModelToResponseDto(v))
	}
	return tDpDto, nil
}

func (s TaskService) GetTaskDependencyByWorkflowDefinitionAndId(ctx context.Context, uId uuid.UUID, wDId uuid.UUID, tDpId uuid.UUID) (TaskDependencyResponseDto, error) {
	if !s.WS.IsUserAndWorkflowDefinitionExists(ctx, uId, wDId) {

	}

	tDp, err := s.Repository.GetTaskDependencyByWorkflowDefinitionAndId(ctx, wDId, tDpId)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	return taskDependencyModelToResponseDto(tDp), nil
}

func (s TaskService) UpdateTaskDependencyById(ctx context.Context, uId uuid.UUID, tDpId uuid.UUID, tDpDto TaskDependencyRequestDto) (TaskDependencyResponseDto, error) {
	if !s.isTaskDependencyValid(ctx, uId, tDpDto.TaskID, tDpDto.DependOnTaskID) {

	}

	newTDp := map[string]any{}
	newTDp["task_id"] = tDpDto.TaskID
	newTDp["depend_on_task_id"] = tDpDto.DependOnTaskID
	tDp, err := s.Repository.UpdateTaskDependencyById(ctx, tDpId, newTDp)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	return taskDependencyModelToResponseDto(tDp), nil
}

func (s TaskService) SoftDeleteTaskDependency(ctx context.Context, uId uuid.UUID, tDpId uuid.UUID) (TaskDependencyResponseDto, error) {
	if !s.isUserAndTaskDepedencyExists(ctx, uId, tDpId) {

	}

	tDp, err := s.Repository.SoftDeleteTaskDependency(ctx, tDpId)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	return taskDependencyModelToResponseDto(tDp), nil
}

func (s TaskService) HardDeleteTaskDependency(ctx context.Context, uId uuid.UUID, tDpId uuid.UUID) (TaskDependencyResponseDto, error) {
	if !s.isUserAndTaskDepedencyExists(ctx, uId, tDpId) {

	}

	tDp, err := s.Repository.HardDeleteTaskDependency(ctx, tDpId)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	return taskDependencyModelToResponseDto(tDp), nil
}

// TaskRun

func (s TaskService) CreateTaskRun(ctx context.Context, tRDto TaskRunRequestDto) (TaskRunResponseDto, error) {
	tR := taskRunRequestDtoToModel(tRDto)
	tR, err := s.Repository.CreateTaskRun(ctx, tR)
	if err != nil {
		return TaskRunResponseDto{}, err
	}
	return taskRunModelToResponseDto(tR), nil
}

func (s TaskService) GetTaskRunByWorkflowRun(ctx context.Context, wRId uuid.UUID) ([]TaskRunResponseDto, error) {
	tRs, err := s.Repository.GetTaskRunByWorkflowRun(ctx, wRId)
	if err != nil {
		return nil, err
	}
	dtos := []TaskRunResponseDto{}
	for _, v := range tRs {
		dtos = append(dtos, taskRunModelToResponseDto(v))
	}
	return dtos, nil
}

func (s TaskService) GetTaskRunByWorkflowRunAndId(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID) (TaskRunResponseDto, error) {
	tR, err := s.Repository.GetTaskRunByWorkflowRunAndId(ctx, wRId, tRId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}
	return taskRunModelToResponseDto(tR), nil
}

func (s TaskService) UpdateTaskRunById(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID, tRDto TaskRunRequestDto) (TaskRunResponseDto, error) {
	newTR := map[string]any{}
	if tRDto.MaxRetries != nil {
		newTR["max_retries"] = *tRDto.MaxRetries
	}
	tR, err := s.Repository.UpdateTaskRunById(ctx, wRId, tRId, newTR)
	if err != nil {
		return TaskRunResponseDto{}, err
	}
	return taskRunModelToResponseDto(tR), nil
}

func (s TaskService) SoftDeleteTaskRun(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID) (TaskRunResponseDto, error) {
	tR, err := s.Repository.SoftDeleteTaskRun(ctx, wRId, tRId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}
	return taskRunModelToResponseDto(tR), nil
}

func (s TaskService) HardDeleteTaskRun(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID) (TaskRunResponseDto, error) {
	tR, err := s.Repository.HardDeleteTaskRun(ctx, wRId, tRId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}
	return taskRunModelToResponseDto(tR), nil
}

// TaskAttempt

// func (s TaskService) CreateTaskAttempt(ctx context.Context, tADto TaskAttemptRequestDto) (TaskAttemptResponseDto, error) {
// 	tA := taskAttemptRequestDtoToModel(tADto)
// 	tA, err := s.Repository.CreateTaskAttempt(ctx, tA)
// 	if err != nil {
// 		return TaskAttemptResponseDto{}, err
// 	}
// 	return taskAttemptModelToResponseDto(tA), nil
// }

func (s TaskService) GetTaskAttemptByTaskRun(ctx context.Context, tRId uuid.UUID) ([]TaskAttemptResponseDto, error) {
	tAs, err := s.Repository.GetTaskAttemptByTaskRun(ctx, tRId)
	if err != nil {
		return nil, err
	}
	dtos := []TaskAttemptResponseDto{}
	for _, v := range tAs {
		dtos = append(dtos, taskAttemptModelToResponseDto(v))
	}
	return dtos, nil
}

func (s TaskService) GetTaskAttemptByTaskRunAndId(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID) (TaskAttemptResponseDto, error) {
	tA, err := s.Repository.GetTaskAttemptByTaskRunAndId(ctx, tRId, tAId)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}
	return taskAttemptModelToResponseDto(tA), nil
}

// func (s TaskService) UpdateTaskAttemptById(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID, tADto TaskAttemptRequestDto) (TaskAttemptResponseDto, error) {
// 	newTA := map[string]any{}
// 	if tADto.WorkerID != nil {
// 		newTA["worker_id"] = *tADto.WorkerID
// 	}
// 	if tADto.Status != nil {
// 		newTA["status"] = *tADto.Status
// 	}
// 	if tADto.StartedAt != nil {
// 		newTA["started_at"] = tADto.StartedAt
// 	}
// 	if tADto.CompletedAt != nil {
// 		newTA["completed_at"] = tADto.CompletedAt
// 	}
// 	if tADto.Log != nil {
// 		newTA["log"] = tADto.Log
// 	}
// 	tA, err := s.Repository.UpdateTaskAttemptById(ctx, tRId, tAId, newTA)
// 	if err != nil {
// 		return TaskAttemptResponseDto{}, err
// 	}
// 	return taskAttemptModelToResponseDto(tA), nil
// }

func (s TaskService) SoftDeleteTaskAttempt(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID) (TaskAttemptResponseDto, error) {
	tA, err := s.Repository.SoftDeleteTaskAttempt(ctx, tRId, tAId)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}
	return taskAttemptModelToResponseDto(tA), nil
}

func (s TaskService) HardDeleteTaskAttempt(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID) (TaskAttemptResponseDto, error) {
	tA, err := s.Repository.HardDeleteTaskAttempt(ctx, tRId, tAId)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}
	return taskAttemptModelToResponseDto(tA), nil
}
