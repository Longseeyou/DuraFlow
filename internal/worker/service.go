package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
	executorimpl "github.com/Longseeyou/DuraFlow/internal/task/executor_impl"
)

type Worker struct {
	workerID string
	consumer message.Consumer
	producer message.Producer
}

func NewWorker(workerID string, consumer message.Consumer, producer message.Producer) Worker {
	return Worker{workerID: workerID, consumer: consumer, producer: producer}
}

func (w Worker) Run(ctx context.Context) {
	slog.Info("Worker Run started", "workerID", w.workerID)
	for {
		msg, err := w.consumer.ReceiveMessage(ctx)
		if err != nil {
			slog.Error(
				"Worker Run consumer.ReceiveMessage",
				"workerID",
				w.workerID,
				"error",
				err,
			)
			continue
		}

		switch string(msg.Key[:]) {
		case "TaskCommandRequest":
			var tCRequest task.TaskCommandRequest
			err = json.Unmarshal(msg.Value, &tCRequest)
			if err != nil {
				slog.Error("Worker Run json.Unmarshal", "workerID", w.workerID, "error", err)
				continue
			}

			startedAt := time.Now()

			executor, err := executorimpl.NewTaskExecutor(tCRequest.TaskType)
			if err != nil {
				slog.Error(
					"Worker Run executorimpl.NewTaskExecutor",
					"workerID",
					w.workerID,
					"error",
					err,
				)
				continue
			}

			taskOutput, taskLog, taskRunStatus, err := executor.Execute(ctx, *tCRequest.Input)

			if err != nil {
				slog.Error("Worker Run executor.Execute", "workerID", w.workerID, "error", err)
				taskRunStatus = task.TASK_ATTEMPT_FAILED
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
