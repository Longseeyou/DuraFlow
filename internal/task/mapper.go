package task

import "fmt"

func taskDefinitionRequestDtoToModel(tDDto TaskDefinitionRequestDto) TaskDefinition {
	var tD TaskDefinition
	tD.Name = *tDDto.Name
	tD.Description = *tDDto.Description
	return tD
}

func taskDefinitionModelToResponseDto(tD TaskDefinition) TaskDefinitionResponseDto {
	var tDDto TaskDefinitionResponseDto
	tDDto.BaseModelToResponseDto(tD)
	tDDto.Name = tD.Name
	tDDto.Description = tD.Description
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
	fmt.Println("missing ScheduledAt, StartedAt, CompletedAt")
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
