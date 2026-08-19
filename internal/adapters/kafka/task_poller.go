package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Longseeyou/DuraFlow/internal/message"
	"github.com/Longseeyou/DuraFlow/internal/orchestrator"
	"github.com/Longseeyou/DuraFlow/internal/shared/model"
	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/google/uuid"
)

type KafkaTaskPoller struct {
	consumer message.Consumer
}

func NewKafkaTaskPoller(consumer message.Consumer) orchestrator.TaskPoller {
	return KafkaTaskPoller{consumer: consumer}
}

type cdcTaskRunEvent struct {
	Payload struct {
		Op    string `json:"op"`
		After struct {
			ID               uuid.UUID          `json:"id"`
			WorkflowRunID    uuid.UUID          `json:"workflow_run_id"`
			TaskDefinitionID uuid.UUID          `json:"task_definition_id"`
			Status           task.TaskRunStatus `json:"status"`
			RetryCount       uint               `json:"retry_count"`
			MaxRetries       uint               `json:"max_retries"`
		} `json:"after"`
	} `json:"payload"`
}

func (poller KafkaTaskPoller) PollTaskRun(
	ctx context.Context,
	numberOfTasks uint,
) ([]task.TaskRun, error) {
	for {
		msg, err := poller.consumer.ReceiveMessage(ctx)
		if err != nil {
			return nil, err
		}

		var event cdcTaskRunEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			slog.Error("KafkaTaskPoller unmarshal CDC event", "error", err)
			msg.Ack()
			continue
		}

		op := event.Payload.Op
		after := event.Payload.After
		if op != "c" && op != "u" && op != "r" {
			msg.Ack()
			continue
		}
		if after.Status != task.TASK_RUN_PENDING && after.Status != task.TASK_RUN_FAILED {
			msg.Ack()
			continue
		}

		msg.Ack()
		return []task.TaskRun{
			{
				BaseModel:        model.BaseModel{ID: after.ID},
				WorkflowRunID:    after.WorkflowRunID,
				TaskDefinitionID: after.TaskDefinitionID,
				Status:           after.Status,
				RetryCount:       after.RetryCount,
				MaxRetries:       after.MaxRetries,
			},
		}, nil
	}
}
