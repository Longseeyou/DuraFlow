package orchestrator

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Longseeyou/DuraFlow/internal/message"
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
			msg.Ack()
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
	return nil
}
