package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
	executorimpl "github.com/Longseeyou/DuraFlow/internal/task/executor_impl"
)

type Worker struct {
	consumer message.Consumer
	producer message.Producer
}

func NewWorker(consumer message.Consumer, producer message.Producer) Worker {
	return Worker{consumer: consumer, producer: producer}
}

func (w Worker) Run(ctx context.Context) {
	for {
		fmt.Println("for")
		msg, err := w.consumer.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println("ReceiveMessage:", err)
			continue
		}
		fmt.Println("ReceiveMessage")

		switch string(msg.Key[:]) {
		case "TaskCommandRequest":
			var tCRequest task.TaskCommandRequest
			err = json.Unmarshal(msg.Value, &tCRequest)
			if err != nil {
				fmt.Println("Decode:", err)
				continue
			}

			startedAt := time.Now()

			executor, err := executorimpl.NewTaskExecutor(tCRequest.TaskType)
			if err != nil {

			}

			taskOutput, taskLog, taskRunStatus, err := executor.Execute(ctx, *tCRequest.Input)

			if err != nil {
				taskRunStatus = task.FAILED
			}

			endedAt := time.Now()

			var tCResponse task.TaskCommandResponse
			tCResponse.Status = &taskRunStatus
			tCResponse.StartedAt = &startedAt
			tCResponse.EndedAt = &endedAt
			tCResponse.Output = &taskOutput
			tCResponse.Log = &taskLog

			value, err := json.Marshal(tCResponse)
			if err != nil {

			}

			err = w.producer.SendMessage(
				ctx,
				"TaskCommandResponse",
				"TaskCommandResponse",
				value,
			)
			if err != nil {

			}
		}
	}
}
