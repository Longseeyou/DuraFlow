package task

import (
	"context"
	"fmt"

	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type TaskService struct {
	repository TaskRepository
	wS         *workflow.WorkflowService
}

func NewTaskService(
	taskRepository TaskRepository,
	workflowService *workflow.WorkflowService,
) TaskService {
	return TaskService{
		repository: taskRepository,
		wS:         workflowService,
	}
}

// TaskDefinition

func (tS *TaskService) CreateTaskDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
	tDDto TaskDefinitionRequestDto,
) (TaskDefinitionResponseDto, error) {
	_, err := tS.wS.GetWorkflowDefinitionByUserAndId(ctx, uId, wDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	tD := taskDefinitionRequestDtoToModel(tDDto)
	tD.WorkflowDefinitionID = wDId

	tD, err = tS.repository.CreateTaskDefinition(ctx, tD)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (tS *TaskService) GetTaskDefinitionByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]TaskDefinitionResponseDto, error) {
	tD, err := tS.repository.GetTaskDefinitionByWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
		return nil, err
	}

	tDDto := []TaskDefinitionResponseDto{}
	for _, v := range tD {
		tDDto = append(tDDto, taskDefinitionModelToResponseDto(v))
	}

	return tDDto, nil
}

func (tS *TaskService) getTaskDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (TaskDefinition, error) {
	tD, err := tS.repository.GetTaskDefinitionByUserAndId(ctx, uId, tDId)
	if err != nil {
		return TaskDefinition{}, err
	}

	return tD, nil
}

func (tS *TaskService) GetTaskDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (TaskDefinitionResponseDto, error) {
	tD, err := tS.getTaskDefinitionByUserAndId(ctx, uId, tDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (tS *TaskService) UpdateTaskDefinitionById(
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

	tD, err := tS.repository.UpdateTaskDefinitionByUserAndId(ctx, uId, tDId, newTD)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

func (tS *TaskService) SoftDeleteTaskDefinition(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (TaskDefinitionResponseDto, error) {
	tD, err := tS.repository.SoftDeleteTaskDefinition(ctx, uId, tDId)
	if err != nil {
		return TaskDefinitionResponseDto{}, err
	}

	return taskDefinitionModelToResponseDto(tD), nil
}

// TaskDependency

func (tS *TaskService) isTaskDependencyValid(
	ctx context.Context,
	uId uuid.UUID,
	tId uuid.UUID,
	dpTId uuid.UUID,
) (bool, error) {
	t, err := tS.getTaskDefinitionByUserAndId(ctx, uId, tId)
	if err != nil {
		return false, err
	}

	dpT, err := tS.getTaskDefinitionByUserAndId(ctx, uId, dpTId)
	if err != nil {
		return false, err
	}

	return t.WorkflowDefinitionID == dpT.WorkflowDefinitionID, nil
}

func (tS *TaskService) CreateTaskDependency(
	ctx context.Context,
	uId uuid.UUID,
	tDpDto TaskDependencyRequestDto,
) (TaskDependencyResponseDto, error) {
	valid, err := tS.isTaskDependencyValid(ctx, uId, tDpDto.TaskID, tDpDto.DependOnTaskID)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	if !valid {
		return TaskDependencyResponseDto{}, fmt.Errorf("")
	}

	tDp := taskDependencyRequestDtoToModel(tDpDto)

	tDp, err = tS.repository.CreateTaskDependency(ctx, tDp)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

func (tS *TaskService) GetTaskDependencyByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]TaskDependencyResponseDto, error) {
	tDp, err := tS.repository.GetTaskDependencyByWorkflowDefinition(ctx, uId, wDId)
	if err != nil {
		return nil, err
	}

	tDpDto := []TaskDependencyResponseDto{}
	for _, v := range tDp {
		tDpDto = append(tDpDto, taskDependencyModelToResponseDto(v))
	}

	return tDpDto, nil
}

func (tS *TaskService) GetTaskDependencyByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) (TaskDependencyResponseDto, error) {
	tDp, err := tS.repository.GetTaskDependencyByUserAndId(ctx, uId, tDpId)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

func (tS *TaskService) UpdateTaskDependencyById(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
	tDpDto TaskDependencyRequestDto,
) (TaskDependencyResponseDto, error) {
	valid, err := tS.isTaskDependencyValid(ctx, uId, tDpDto.TaskID, tDpDto.DependOnTaskID)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	if !valid {
		return TaskDependencyResponseDto{}, fmt.Errorf("")
	}

	newTDp := map[string]any{}
	newTDp["task_id"] = tDpDto.TaskID
	newTDp["depend_on_task_id"] = tDpDto.DependOnTaskID

	tDp, err := tS.repository.UpdateTaskDependencyByUserAndId(ctx, uId, tDpId, newTDp)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}
	return taskDependencyModelToResponseDto(tDp), nil
}

func (tS *TaskService) SoftDeleteTaskDependency(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) (TaskDependencyResponseDto, error) {
	tDp, err := tS.repository.SoftDeleteTaskDependency(ctx, uId, tDpId)
	if err != nil {
		return TaskDependencyResponseDto{}, err
	}

	return taskDependencyModelToResponseDto(tDp), nil
}

func (tS *TaskService) IsTaskDependencyDag(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (bool, error) {
	tDps, err := tS.repository.GetTaskDependencyByWorkflowDefinition(ctx, uId, wDId)
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
	visit = func(tDId uuid.UUID) error {
		switch mark[tDId] {
		case visiting:
			return fmt.Errorf("graph has at least one cycle")
		case visited:
			return nil
		}

		mark[tDId] = visiting

		for _, nextTDId := range graph[tDId] {
			err := visit(nextTDId)
			if err != nil {
				return err
			}
		}

		mark[tDId] = visited
		return nil
	}

	for tDId := range graph {
		err = visit(tDId)
		if err != nil {
			return false, nil
		}
	}

	return true, nil
}

// TaskRun

// func (tS *TaskService) CreateTaskRun(
// 	ctx context.Context,
// 	uId uuid.UUID,
// 	wRId uuid.UUID,
// 	tDId uuid.UUID,
// 	tRDto TaskRunRequestDto,
// ) (TaskRunResponseDto, error) {
// 	_, err := tS.wS.GetWorkflowRunByUserAndId(ctx, uId, wRId)
// 	if err != nil {
// 		return TaskRunResponseDto{}, err
// 	}
//
// 	_, err = tS.GetTaskDefinitionByUserAndId(ctx, uId, tDId)
// 	if err != nil {
// 		return TaskRunResponseDto{}, err
// 	}
//
// 	tR := taskRunRequestDtoToModel(tRDto)
// 	tR.WorkflowRunID = wRId
// 	tR.TaskDefinitionID = tDId
//
// 	tR, err = tS.repository.CreateTaskRun(ctx, tR)
// 	if err != nil {
// 		return TaskRunResponseDto{}, err
// 	}
// 	return taskRunModelToResponseDto(tR), nil
// }

func (tS *TaskService) GetTaskRunByWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) ([]TaskRunResponseDto, error) {
	tRs, err := tS.repository.GetTaskRunByWorkflowRun(ctx, uId, wRId)
	if err != nil {
		return nil, err
	}

	dtos := []TaskRunResponseDto{}
	for _, v := range tRs {
		dtos = append(dtos, taskRunModelToResponseDto(v))
	}

	return dtos, nil
}

func (tS *TaskService) GetTaskRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (TaskRunResponseDto, error) {
	tR, err := tS.repository.GetTaskRunByUserAndId(ctx, uId, tRId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

func (tS *TaskService) UpdateTaskRunByUserAndId(
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

	tR, err := tS.repository.UpdateTaskRunByUserAndId(ctx, uId, tRId, newTR)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

func (tS *TaskService) SoftDeleteTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (TaskRunResponseDto, error) {
	tR, err := tS.repository.SoftDeleteTaskRun(ctx, uId, tRId)
	if err != nil {
		return TaskRunResponseDto{}, err
	}

	return taskRunModelToResponseDto(tR), nil
}

// TaskAttempt

func (tS *TaskService) GetTaskAttemptByTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) ([]TaskAttemptResponseDto, error) {
	tAs, err := tS.repository.GetTaskAttemptByTaskRun(ctx, uId, tRId)
	if err != nil {
		return nil, err
	}

	dtos := []TaskAttemptResponseDto{}
	for _, v := range tAs {
		dtos = append(dtos, taskAttemptModelToResponseDto(v))
	}

	return dtos, nil
}

func (tS *TaskService) GetTaskAttemptByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (TaskAttemptResponseDto, error) {
	tA, err := tS.repository.GetTaskAttemptByUserAndId(ctx, uId, tAId)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}

	return taskAttemptModelToResponseDto(tA), nil
}

func (tS *TaskService) SoftDeleteTaskAttempt(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (TaskAttemptResponseDto, error) {
	tA, err := tS.repository.SoftDeleteTaskAttempt(ctx, uId, tAId)
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

func (tS *TaskServiceInternal) UpdateTaskRunById(
	ctx context.Context,
	tRId uuid.UUID,
	tRDto TaskRunInternalDto,
) (TaskRun, error) {
	newTR, err := taskRunInternalDtoToMap(tRDto)
	if err != nil {
		return TaskRun{}, err
	}

	tR, err := tS.repository.UpdateTaskRunById(ctx, tRId, newTR)
	if err != nil {
		return TaskRun{}, err
	}

	return tR, nil
}

func (tS *TaskServiceInternal) CreateTaskAttempt(
	ctx context.Context,
	tRId uuid.UUID,
	tADto TaskAttemptInternalDto,
) (TaskAttempt, error) {
	tA := taskAttemptInternalDtoToModel(tADto)
	tA.TaskRunID = tRId

	tA, err := tS.repository.CreateTaskAttempt(ctx, tA)
	if err != nil {
		return TaskAttempt{}, err
	}

	return tA, nil
}

func (tS *TaskServiceInternal) UpdateTaskAttemptById(
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
	if tADto.Log != nil {
		newTA["log"] = tADto.Log
	}

	tA, err := tS.repository.UpdateTaskAttemptById(ctx, tAId, newTA)
	if err != nil {
		return TaskAttemptResponseDto{}, err
	}

	return taskAttemptModelToResponseDto(tA), nil
}
