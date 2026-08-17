package worker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
	executorimpl "github.com/Longseeyou/DuraFlow/internal/task/executor_impl"
)

type Worker struct {
	workerID   string
	consumer   message.Consumer
	producer   message.Producer
	repository WorkerRepository
}

func NewWorker(
	workerID string,
	consumer message.Consumer,
	producer message.Producer,
	repository WorkerRepository,
) Worker {
	return Worker{
		workerID:   workerID,
		consumer:   consumer,
		producer:   producer,
		repository: repository,
	}
}

func (w Worker) Run(ctx context.Context) {
	slog.Info("Worker Run started", "workerID", w.workerID)
	for {
		msg, err := w.consumer.ReceiveMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
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
			err = w.ExecuteTaskCommandRequest(ctx, msg)
			if err != nil {
				slog.Error("Worker Run ExecuteTaskCommandRequest", "error", err)
			}
			msg.Ack()
		}
	}
}

func (w Worker) ExecuteTaskCommandRequest(
	ctx context.Context,
	msg *message.Message,
) error {
	// Decode message
	var tCRequest task.TaskCommandRequest
	err := json.Unmarshal(msg.Value, &tCRequest)
	if err != nil {
		return err
	}

	// Idempotency: only start if the run is queued for this exact attempt number
	ok, err := w.repository.MarkTaskRunRunning(
		ctx,
		tCRequest.TaskRunID,
		tCRequest.AttemptNumber,
	)
	if err != nil {
		return err

	}
	if !ok {
		return errors.New("Idempotency")
	}

	executor, err := executorimpl.NewTaskExecutor(tCRequest.TaskType)
	if err != nil {
		return err
	}

	// Execute
	startedAt := time.Now()

	tR, err := w.repository.UpdateTaskRunByID(
		ctx,
		tCRequest.TaskRunID,
		map[string]any{"timeout_at": startedAt.Add(tCRequest.Timeout)},
	)
	if err != nil {
		return err
	}

	slog.Info("Worker ExecuteTaskCommandRequest executor.Execute start", "workerID", w.workerID)

	executor_ctx, cancel := context.WithTimeout(ctx, tCRequest.Timeout)
	defer cancel()
	taskOutput, taskLog, err := executor.Execute(
		executor_ctx,
		*tR.Input,
	)

	endedAt := time.Now()
	slog.Info("Worker ExecuteTaskCommandRequest executor.Execute end", "workerID", w.workerID)

	// Update TaskRun and create TaskAttempt
	newTaskRun := map[string]any{"retry_count": tCRequest.AttemptNumber}

	tA := task.TaskAttempt{
		TaskRunID:     tR.ID,
		AttemptNumber: tCRequest.AttemptNumber,
		WorkerID:      w.workerID,
		StartedAt:     &startedAt,
		EndedAt:       &endedAt,
		Log:           &taskLog,
	}

	if err == nil {
		newTaskRun["status"] = task.TASK_RUN_COMPLETED
		newTaskRun["ended_at"] = endedAt
		newTaskRun["output"] = taskOutput

		tA.Status = task.TASK_ATTEMPT_COMPLETED
	} else {
		newTaskRun["status"] = task.TASK_RUN_FAILED

		tA.Status = task.TASK_ATTEMPT_FAILED
	}

	tR, err = w.repository.UpdateTaskRunByID(ctx, tCRequest.TaskRunID, newTaskRun)
	if err != nil {
		return err
	}

	tA, err = w.repository.CreateTaskAttempt(ctx, tA)
	if err != nil {
		return err
	}

	return nil
}
