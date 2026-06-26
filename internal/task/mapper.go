package task

import (
	"fmt"
)

func taskDefinitionRequestDtoToModel(tDDto TaskDefinitionRequestDto) TaskDefinition {
	var tD TaskDefinition
	tD.Name = *tDDto.Name
	tD.Description = *tDDto.Description
	tD.TaskType = *tDDto.TaskType
	tD.Timeout = *tDDto.Timeout
	return tD
}

func taskDefinitionModelToResponseDto(tD TaskDefinition) TaskDefinitionResponseDto {
	var tDDto TaskDefinitionResponseDto
	tDDto.BaseModelToResponseDto(tD)
	tDDto.Name = tD.Name
	tDDto.Description = tD.Description
	tDDto.TaskType = tD.TaskType
	tDDto.Timeout = tD.Timeout
	return tDDto
}

func taskDependencyRequestDtoToModel(tDDto TaskDependencyRequestDto) TaskDependency {
	var tD TaskDependency
	tD.TaskID = tDDto.TaskID
	tD.DependOnTaskID = tDDto.DependOnTaskID
	return tD
}

func taskDependencyModelToResponseDto(tD TaskDependency) TaskDependencyResponseDto {
	var tDDto TaskDependencyResponseDto
	tDDto.BaseModelToResponseDto(tD)
	tDDto.TaskID = tD.TaskID
	tDDto.DependOnTaskID = tD.DependOnTaskID
	return tDDto
}

func taskRunRequestDtoToModel(tRDto TaskRunRequestDto) TaskRun {
	var tR TaskRun
	tR.MaxRetries = *tRDto.MaxRetries
	tR.Input = tRDto.Input
	return tR
}

func taskRunModelToResponseDto(tR TaskRun) TaskRunResponseDto {
	var tRDto TaskRunResponseDto
	tRDto.BaseModelToResponseDto(tR)
	tRDto.WorkflowRunID = tR.WorkflowRunID
	tRDto.TaskDefinitionID = tR.TaskDefinitionID
	tRDto.Status = tR.Status
	tRDto.RetryCount = tR.RetryCount
	tRDto.MaxRetries = tR.MaxRetries
	tRDto.Input = tR.Input
	tRDto.Output = tR.Output
	return tRDto
}

func taskAttemptModelToResponseDto(tA TaskAttempt) TaskAttemptResponseDto {
	var tADto TaskAttemptResponseDto
	tADto.BaseModelToResponseDto(tA)
	tADto.TaskRunID = tA.TaskRunID
	tADto.AttemptNumber = tA.AttemptNumber
	tADto.WorkerID = tA.WorkerID
	tADto.Status = tA.Status
	return tADto
}

// Internal

func taskRunInternalDtoToMap(tRDto TaskRunInternalDto) (map[string]any, error) {
	tR := map[string]any{}

	switch *tRDto.Status {
	case TASK_RUN_QUEUED:
		if tRDto.ScheduledAt == nil {
			return nil, fmt.Errorf("Missing ScheduledAt for status %s", *tRDto.Status)
		}
		tR["ScheduledAt"] = *tRDto.ScheduledAt
	case TASK_RUN_RUNNING:
		if tRDto.StartedAt == nil {
			return nil, fmt.Errorf("Missing StartedAt for status %s", *tRDto.Status)
		}
		tR["StartedAt"] = *tRDto.StartedAt
	case TASK_RUN_COMPLETED, TASK_RUN_FAILED:
		if tRDto.EndedAt == nil {
			return nil, fmt.Errorf("Missing EndedAt for status %s", *tRDto.Status)
		}
		tR["EndedAt"] = *tRDto.EndedAt
	}

	if tRDto.RetryCount != nil {
		tR["RetryCount"] = *tRDto.RetryCount
	}

	if tRDto.Input != nil {
		tR["Input"] = *tRDto.Input
	}

	if tRDto.Output != nil {
		tR["Output"] = *tRDto.Output
	}

	if tRDto.Error != nil {
		tR["Error"] = *tRDto.Error
	}

	return tR, nil
}

func taskAttemptInternalDtoToModel(tADto TaskAttemptInternalDto) TaskAttempt {
	var tA TaskAttempt
	tA.AttemptNumber = *tADto.AttemptNumber
	tA.WorkerID = *tADto.WorkerID
	tA.Status = *tADto.Status
	tA.StartedAt = tADto.StartedAt
	return tA
}

func taskAttemptInternalDtoToMap(tADto TaskAttemptInternalDto) (map[string]any, error) {
	tA := map[string]any{}
	return tA, nil
}
