package orchestrator

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type Orchestrator struct {
	orchestratorRepository OrchestratorRepository
}

func (oR *Orchestrator) RunWorkflow(ctx context.Context, workflowDefinitionID uuid.UUID) error {
	wD, err := oR.orchestratorRepository.GetWorkflowDefinitionByID(ctx, workflowDefinitionID)
	if err != nil {

	}
	if wD.Status != workflow.WORKFLOW_DEFINITION_ACTIVATED {

	}

	workflowRun := workflow.WorkflowRun{WorkflowDefinitionID: workflowDefinitionID, Status: workflow.WORKFLOW_RUN_PENDING}
	workflowRun, err = oR.orchestratorRepository.CreateWorkflowRun(ctx, workflowRun)
	if err != nil {

	}

	tDs, err := oR.orchestratorRepository.GetTaskDefinitionsByWorkflowDefinition(ctx, workflowDefinitionID)
	if err != nil {

	}
	for _, tD := range tDs {
		oR.orchestratorRepository.CreateTaskRun(
			ctx,
			task.TaskRun{
				WorkflowRunID:    workflowRun.ID,
				TaskDefinitionID: tD.ID,
				Status:           task.TASK_RUN_PENDING,
			},
		)
	}

	return nil
}

func (oR *Orchestrator) ResolveTimeoutTaskRun() {

}

func (oR *Orchestrator) ResolveDeadLetter() {

}
