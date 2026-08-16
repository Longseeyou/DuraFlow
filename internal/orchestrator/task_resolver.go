package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
)

type TaskResolver struct {
	orchestratorRepository OrchestratorRepository
	consumer               message.Consumer
}

func NewTaskResolver(
	orchestratorRepository OrchestratorRepository,
	consumer message.Consumer,
) *TaskResolver {
	return &TaskResolver{orchestratorRepository: orchestratorRepository, consumer: consumer}
}

func (taskResolver *TaskResolver) Run(ctx context.Context) {
	for {
		msg, err := taskResolver.consumer.ReceiveMessage(ctx)
		if err != nil {
			continue
		}

		switch string(msg.Key[:]) {
		case "TaskCommandResponse":
			err := taskResolver.ResolveTaskCommandResponse(ctx, msg)
			if err != nil {

			}
		}
	}
}

func (taskResolver *TaskResolver) ResolveTaskCommandResponse(
	ctx context.Context,
	msg *message.Message,
) error {
	var tCResponse task.TaskCommandResponse
	err := json.Unmarshal(msg.Value, &tCResponse)
	if err != nil {
		return err
	}

	// Idempotency
	ok, err := taskResolver.orchestratorRepository.TaskAttemptIdempotency(
		ctx,
		tCResponse.TaskAttemptID,
		task.TaskRunStatus(tCResponse.Status),
	)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("ResolveTaskCommandResponse Idempotency")
	}

	// Update TaskAttempt
	taskAttempt, err := taskResolver.orchestratorRepository.UpdateTaskAttemptById(
		ctx,
		tCResponse.TaskAttemptID,
		map[string]any{"Log": tCResponse.Log},
	)
	if err != nil {
		return err
	}

	// Update TaskRun
	newTR := map[string]any{"RetryCount": taskAttempt.AttemptNumber + 1}

	switch tCResponse.Status {
	case task.TASK_ATTEMPT_COMPLETED:
		newTR["Status"] = task.TASK_RUN_COMPLETED
		newTR["Output"] = tCResponse.Output
		newTR["EndedAt"] = tCResponse.EndedAt
	case task.TASK_ATTEMPT_FAILED:
		newTR["Status"] = task.TASK_RUN_FAILED
	}

	_, err = taskResolver.orchestratorRepository.UpdateTaskRunById(
		ctx,
		tCResponse.TaskAttemptID,
		newTR,
	)
	if err != nil {
		return err
	}

	return nil
}
