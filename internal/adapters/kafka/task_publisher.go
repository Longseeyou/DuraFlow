package kafka

import (
	"context"
	"encoding/json"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
)

type TaskPublisher struct {
	producer message.Producer
	topic    string
}

func NewTaskPublisher(producer message.Producer, topic string) *TaskPublisher {
	return &TaskPublisher{producer: producer, topic: topic}
}

func (publisher *TaskPublisher) PublishTask(ctx context.Context, command task.Command) error {
	payload, err := json.Marshal(command)
	if err != nil {
		return err
	}
	return publisher.producer.SendMessage(
		ctx,
		publisher.topic,
		[]byte(command.WorkflowRunID.String()),
		payload,
	)
}
