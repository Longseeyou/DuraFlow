package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
)

type TaskResolver struct {
	repository OrchestratorRepository
	consumer   message.Consumer
}

func NewTaskResolver(
	repository OrchestratorRepository,
	consumer message.Consumer,
) *TaskResolver {
	return &TaskResolver{repository: repository, consumer: consumer}
}

func (taskResolver *TaskResolver) Run(ctx context.Context) {
	for {
		msg, err := taskResolver.consumer.ReceiveMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("TaskResolver ReceiveMessage", "error", err)
			continue
		}

		switch string(msg.Key[:]) {
		case TaskCommandResponseTopic:
			err := taskResolver.ResolveTaskCommandResponse(ctx, msg)
			if err != nil {
				slog.Error("TaskResolver ResolveTaskCommandResponse", "error", err)
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
	ok, err := taskResolver.repository.TaskAttemptIdempotency(
		ctx,
		tCResponse.TaskAttemptID,
		tCResponse.Status,
	)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("ResolveTaskCommandResponse Idempotency")
	}

	// Update TaskAttempt
	taskAttempt, err := taskResolver.repository.UpdateTaskAttemptByID(
		ctx,
		tCResponse.TaskAttemptID,
		map[string]any{
			"status":     tCResponse.Status,
			"started_at": tCResponse.StartedAt,
			"ended_at":   tCResponse.EndedAt,
			"log":        tCResponse.Log,
		},
	)
	if err != nil {
		return err
	}

	// Update TaskRun
	newTR := map[string]any{"retry_count": taskAttempt.AttemptNumber + 1}

	switch tCResponse.Status {
	case task.TASK_ATTEMPT_COMPLETED:
		newTR["status"] = task.TASK_RUN_COMPLETED
		newTR["output"] = tCResponse.Output
		newTR["ended_at"] = tCResponse.EndedAt
	case task.TASK_ATTEMPT_FAILED:
		newTR["status"] = task.TASK_RUN_FAILED
	}

	_, err = taskResolver.repository.UpdateTaskRunByID(
		ctx,
		tCResponse.TaskRunID,
		newTR,
	)
	if err != nil {
		return err
	}

	return nil
}
