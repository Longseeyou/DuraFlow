package orchestrator

import (
	"context"
	"fmt"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type Orchestrator struct {
	orchestratorRepository OrchestratorRepository
}

func NewOrchestrator(orchestratorRepository OrchestratorRepository) *Orchestrator {
	return &Orchestrator{orchestratorRepository: orchestratorRepository}
}

func (oR *Orchestrator) RunWorkflow(ctx context.Context, workflowDefinitionID uuid.UUID) error {
	wD, err := oR.orchestratorRepository.GetWorkflowDefinitionByID(ctx, workflowDefinitionID)
	if err != nil {
		return err
	}
	if wD.Status != workflow.WORKFLOW_DEFINITION_ACTIVATED {
		return fmt.Errorf(
			"workflow definition %s is not activated",
			workflowDefinitionID,
		)
	}

	workflowRun := workflow.WorkflowRun{
		WorkflowDefinitionID: workflowDefinitionID,
		Status:               workflow.WORKFLOW_RUN_PENDING,
	}
	workflowRun, err = oR.orchestratorRepository.CreateWorkflowRun(ctx, workflowRun)
	if err != nil {
		return err
	}

	tDs, err := oR.orchestratorRepository.GetTaskDefinitionsByWorkflowDefinition(
		ctx,
		workflowDefinitionID,
	)
	if err != nil {
		return err
	}
	for _, tD := range tDs {
		_, err := oR.orchestratorRepository.CreateTaskRun(
			ctx,
			task.TaskRun{
				WorkflowRunID:    workflowRun.ID,
				TaskDefinitionID: tD.ID,
				Status:           task.TASK_RUN_PENDING,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (oR *Orchestrator) ResolveTimeoutTaskRun() {

}

func (oR *Orchestrator) ResolveDeadLetter() {

}
