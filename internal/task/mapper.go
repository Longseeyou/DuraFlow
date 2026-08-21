package task

func taskDefinitionModelToResponseDto(tD TaskDefinition) TaskDefinitionResponseDto {
	var tDDto TaskDefinitionResponseDto
	tDDto.BaseModelToResponseDto(tD)
	tDDto.Name = tD.Name
	tDDto.Description = tD.Description
	tDDto.TaskType = tD.TaskType
	tDDto.Timeout = tD.Timeout
	return tDDto
}

func taskDependencyModelToResponseDto(tD TaskDependency) TaskDependencyResponseDto {
	var tDDto TaskDependencyResponseDto
	tDDto.BaseModelToResponseDto(tD)
	tDDto.TaskID = tD.TaskID
	tDDto.DependOnTaskID = tD.DependOnTaskID
	return tDDto
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
