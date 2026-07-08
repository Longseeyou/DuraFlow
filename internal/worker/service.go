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
	workerID       string
	consumer       message.Consumer
	producer       message.Producer
	taskRepository task.TaskRepositoryInternal
}

func NewWorker(
	workerID string,
	consumer message.Consumer,
	producer message.Producer,
	taskRepository task.TaskRepositoryInternal,
) Worker {
	return Worker{
		workerID:       workerID,
		consumer:       consumer,
		producer:       producer,
		taskRepository: taskRepository,
	}
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

			valid, err := w.taskRepository.TaskAttemptIdempotency(
				ctx,
				tCRequest.TaskAttemptID,
				task.TASK_ATTEMPT_RUNNING,
			)
			if err != nil {
				continue
			}
			if !valid {
				slog.Error(
					"Worker Run taskRepository.TaskAttemptIdempotency",
					"workerID",
					w.workerID,
					"error",
					"Idempotency",
				)
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

				err = w.producer.SendMessage(
					ctx,
					"TaskCommandResponse",
					"Error",
					[]byte(err.Error()),
				)
				if err != nil {
					slog.Error(
						"Worker Run producer.SendMessage", "workerID", w.workerID, "error", err,
					)
				}
				continue
			}

			slog.Info("Worker Run executor.Execute start", "workerID", w.workerID)
			taskOutput, taskLog, taskAttemptStatus, err := executor.Execute(ctx, *tCRequest.Input)
			slog.Info("Worker Run executor.Execute end", "workerID", w.workerID)

			if err != nil {
				slog.Error("Worker Run executor.Execute", "workerID", w.workerID, "error", err)
				taskAttemptStatus = task.TASK_ATTEMPT_FAILED
			}

			endedAt := time.Now()

			var tCResponse task.TaskCommandResponse
			tCResponse.Status = taskAttemptStatus
			tCResponse.StartedAt = startedAt
			tCResponse.EndedAt = endedAt
			tCResponse.Output = &taskOutput
			tCResponse.Log = &taskLog

			value, err := json.Marshal(tCResponse)
			if err != nil {
				slog.Error("Worker Run json.Marshal", "workerID", w.workerID, "error", err)
				continue
			}

			err = w.producer.SendMessage(
				ctx,
				"TaskCommandResponse",
				"TaskCommandResponse",
				value,
			)
			if err != nil {
				slog.Error(
					"Worker Run producer.SendMessage", "workerID", w.workerID, "error", err,
				)
			}
		}
	}
}
