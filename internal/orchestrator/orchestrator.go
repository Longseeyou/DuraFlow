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

func (oR *Orchestrator) RunWorkflow(ctx context.Context, wDId uuid.UUID) error {
	wD, err := oR.orchestratorRepository.GetWorkflowDefinitionById(ctx, wDId)
	if err != nil {

	}
	if wD.Status != workflow.WORKFLOW_DEFINITION_ACTIVATED {

	}

	wR := workflow.WorkflowRun{WorkflowDefinitionID: wDId, Status: workflow.WORKFLOW_RUN_PENDING}
	wR, err = oR.orchestratorRepository.CreateWorkflowRun(ctx, wR)
	if err != nil {

	}

	tDs, err := oR.orchestratorRepository.GetTaskDefinitionsByWorkflowDefinition(ctx, wDId)
	if err != nil {

	}
	for _, tD := range tDs {
		oR.orchestratorRepository.CreateTaskRun(
			ctx,
			task.TaskRun{
				WorkflowRunID:    wR.ID,
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
