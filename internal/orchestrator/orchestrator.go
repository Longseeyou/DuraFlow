package orchestrator

import (
	"context"
	"log/slog"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

type Orchestrator struct {
	repository             OrchestratorRepository
	resolveTimeoutInterval time.Duration
}

const defaultResolveTimeoutInterval = 5 * time.Second

func NewOrchestrator(repository OrchestratorRepository) *Orchestrator {
	return &Orchestrator{
		repository:             repository,
		resolveTimeoutInterval: defaultResolveTimeoutInterval,
	}
}

func (orchestrator *Orchestrator) RunWorkflow(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) error {
	wD, err := orchestrator.repository.UpdateWorkflowDefinitionByUserAndID(
		ctx,
		userID,
		workflowDefinitionID,
		map[string]any{"status": workflow.WORKFLOW_DEFINITION_RUNNING},
	)
	if err != nil {
		return err
	}

	workflowRun := workflow.WorkflowRun{
		WorkflowDefinitionID: wD.ID,
		Status:               workflow.WORKFLOW_RUN_PENDING,
	}
	workflowRun, err = orchestrator.repository.CreateWorkflowRun(ctx, workflowRun)
	if err != nil {
		return err
	}

	tDs, err := orchestrator.repository.GetTaskDefinitionsByWorkflowDefinition(
		ctx,
		wD.ID,
	)
	if err != nil {
		return err
	}
	predecessorCounts, err := orchestrator.repository.
		GetPredecessorTaskCountsByWorkflowDefinition(ctx, wD.ID)
	if err != nil {
		return err
	}

	for _, tD := range tDs {
		_, err := orchestrator.repository.CreateTaskRun(
			ctx,
			task.TaskRun{
				WorkflowRunID:                      workflowRun.ID,
				TaskDefinitionID:                   tD.ID,
				Status:                             task.TASK_RUN_PENDING,
				MaxRetries:                         10,
				NumberOfIncompletePredecessorTasks: predecessorCounts[tD.ID],
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (orchestrator *Orchestrator) ResolveTimeoutTaskRun(ctx context.Context) {
	ticker := time.NewTicker(orchestrator.resolveTimeoutInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		tRs, err := orchestrator.repository.GetTimedOutTaskRuns(ctx)
		if err != nil {
			slog.Error("Orchestrator ResolveTimeoutTaskRun GetTimedOutTaskRuns", "error", err)
			continue
		}

		for _, tR := range tRs {
			ok, err := orchestrator.repository.TaskRunIdempotency(ctx, tR.ID, task.TASK_RUN_FAILED)
			if err != nil {
				slog.Error(
					"Orchestrator ResolveTimeoutTaskRun",
					"taskRunID",
					tR.ID,
					"error",
					err,
				)
				continue
			}
			if !ok {
				slog.Error(
					"Orchestrator ResolveTimeoutTaskRun",
					"taskRunID",
					tR.ID,
					"error",
					"Idempotency",
				)
				continue
			}
			slog.Info("Orchestrator ResolveTimeoutTaskRun", "taskRunID", tR.ID)
		}
	}
}

func (oR *Orchestrator) ResolveDeadLetter() {

}
