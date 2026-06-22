package worker

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
)

type Worker struct {
	consumer message.Consumer
	producer message.Producer
}

func (w Worker) Run(ctx context.Context) {
	var buf bytes.Buffer
	encode := gob.NewEncoder(&buf)
	decode := gob.NewDecoder(&buf)

	for {
		msg, err := w.consumer.ReceiveMessage(ctx)
		if err != nil {

		}
		switch string(msg.Key[:]) {
		case "Task":
			var tCRequest task.TaskCommandRequest
			err = decode.Decode(tCRequest)
			if err != nil {

			}

			executor, err := task.NewTaskExecutor(tCRequest.TaskType)
			if err != nil {

			}
			taskOutput, taskLog, err := executor.Execute(ctx, *tCRequest.Input)

			var tCResponse task.TaskCommandResponse
			if err != nil {

			}
			tCResponse.Output = &taskOutput
		}
	}
}
