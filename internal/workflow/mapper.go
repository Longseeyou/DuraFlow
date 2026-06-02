package workflow

func workflowRequestDtoToWorkflowModel(wDto WorkflowRequestDto) Workflow {
	w := Workflow{}
	w.Name = *wDto.Name
	w.Description = *wDto.Description
	return w
}

func workflowModelToResponseDto(w Workflow) WorkflowResponseDto {
	wDto := WorkflowResponseDto{}
	wDto.BaseModelToResponseDto(w)
	wDto.Name = w.Name
	wDto.Description = w.Description
	return wDto
}

func workflowDefinitionRequestDtoToWorkflowModel(wDDto WorkflowDefinitionRequestDto) WorkflowDefinition {
	wD := WorkflowDefinition{}
	wD.Version = *wDDto.Version
	wD.Status = *wDDto.Status
	return wD
}

func workflowDefinitionModelToResponseDto(wD WorkflowDefinition) WorkflowDefinitionResponseDto {
	wDDto := WorkflowDefinitionResponseDto{}
	wDDto.BaseModelToResponseDto(wD)
	wDDto.Version = wD.Version
	wDDto.Status = wD.Status
	return wDDto
}
