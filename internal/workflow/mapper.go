package workflow

func workflowModelToResponseDto(w Workflow) WorkflowResponseDto {
	dto := WorkflowResponseDto{}
	dto.BaseModelToResponseDto(w)
	dto.Name = w.Name
	dto.Description = w.Description
	return dto
}

func workflowDefinitionModelToResponseDto(wD WorkflowDefinition) WorkflowDefinitionResponseDto {
	dto := WorkflowDefinitionResponseDto{}
	dto.BaseModelToResponseDto(wD)
	dto.Version = wD.Version
	dto.Status = wD.Status
	return dto
}

func workfloworkflowRununModelToResponseDto(
	workflowRun WorkflowRun,
) WorkflowRunResponseDto {
	dto := WorkflowRunResponseDto{}
	dto.BaseModelToResponseDto(workflowRun)
	dto.Status = workflowRun.Status
	dto.StartedAt = workflowRun.StartedAt
	dto.EndedAt = workflowRun.EndedAt
	return dto
}
