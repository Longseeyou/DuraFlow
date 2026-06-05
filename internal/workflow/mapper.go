package workflow

func workflowRequestDtoToModel(wDto WorkflowRequestDto) Workflow {
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

// func workflowDefinitionRequestDtoToModel(wDDto WorkflowDefinitionRequestDto) WorkflowDefinition {
// 	wD := WorkflowDefinition{}
// 	wD.Version = *wDDto.Version
// 	wD.Status = *wDDto.Status
// 	return wD
// }

func workflowDefinitionModelToResponseDto(wD WorkflowDefinition) WorkflowDefinitionResponseDto {
	wDDto := WorkflowDefinitionResponseDto{}
	wDDto.BaseModelToResponseDto(wD)
	wDDto.Version = wD.Version
	wDDto.Status = wD.Status
	return wDDto
}

// func workflowRunRequestDtoToModel(wRDto WorkflowRunRequestDto) WorkflowRun {
// 	wR := WorkflowRun{}
// 	wR.Status = *wRDto.Status
// 	return wR
// }

func workflowRunModelToResponseDto(wR WorkflowRun) WorkflowRunResponseDto {
	wRDto := WorkflowRunResponseDto{}
	wRDto.BaseModelToResponseDto(wR)
	wRDto.Status = wR.Status
	wRDto.StartedAt = *wR.StartedAt
	wRDto.CompletedAt = *wR.CompletedAt
	wRDto.CancelledAt = *wR.CancelledAt
	return wRDto
}
