package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
)

const (
	TaskCommandRequestTopic  = "TaskCommandRequest"
	TaskCommandResponseTopic = "TaskCommandResponse"

	defaultPollInterval = 500 * time.Millisecond
)

type TaskScheduler struct {
	repository   OrchestratorRepository
	taskPoller   TaskPoller
	producer     message.Producer
	pollInterval time.Duration
}

func NewTaskScheduler(
	repository OrchestratorRepository,
	taskPoller TaskPoller,
	producer message.Producer,
) *TaskScheduler {
	return &TaskScheduler{
		repository:   repository,
		taskPoller:   taskPoller,
		producer:     producer,
		pollInterval: defaultPollInterval,
	}
}

func (taskScheduler *TaskScheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(taskScheduler.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		taskRuns, err := taskScheduler.taskPoller.PollTaskRun(ctx, 10)
		if err != nil {
			slog.Error("TaskScheduler PollTaskRun", "error", err)
			continue
		}

		for _, taskRun := range taskRuns {
			err := taskScheduler.ScheduleTask(ctx, &taskRun)
			if err != nil {
				slog.Error(
					"TaskScheduler ScheduleTask",
					"taskRunID",
					taskRun.ID,
					"error",
					err,
				)
			}
		}
	}
}

func (taskScheduler *TaskScheduler) ScheduleTask(ctx context.Context, taskRun *task.TaskRun) error {
	// Dead letter when retries are exhausted
	if taskRun.Status == task.TASK_RUN_FAILED && taskRun.RetryCount >= taskRun.MaxRetries {
		ok, err := taskScheduler.repository.TaskRunIdempotency(
			ctx,
			taskRun.ID,
			task.TASK_RUN_DEAD_LETTERED,
		)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("task run %s cannot be dead lettered", taskRun.ID)
		}
		slog.Warn("ScheduleTask dead lettered", "taskRunID", taskRun.ID)
		return nil
	}

	// Idempotency
	ok, err := taskScheduler.repository.TaskRunIdempotency(
		ctx,
		taskRun.ID,
		task.TASK_RUN_QUEUED,
	)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("task run %s cannot be queued", taskRun.ID)
	}

	// Load TaskDefinition
	tD, err := taskScheduler.repository.GetTaskDefinitionByID(ctx, taskRun.TaskDefinitionID)
	if err != nil {
		return err
	}

	// Update TaskRun
	newTaskRun := map[string]any{}
	if taskRun.Status == task.TASK_RUN_PENDING {
		newTaskRun["started_at"] = time.Now()
	}

	// Retrieve task's predecessors' output for input
	// TODO: Only pull task's predecessors with new output
	tRs, err := taskScheduler.repository.GetPredecessorTaskRuns(ctx, taskRun.ID)
	if err != nil {
		return err
	}

	input := map[string]string{}
	for _, tR := range tRs {
		if tR.Output == nil {
			continue
		}

		tDp, err := taskScheduler.repository.GetTaskDefinitionByID(
			ctx,
			tR.TaskDefinitionID,
		)
		if err != nil {
			return err
		}

		input[tDp.Name] = *tR.Output
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return err
	}
	newTaskRun["input"] = string(inputJSON)

	tR, err := taskScheduler.repository.UpdateTaskRunByID(ctx, taskRun.ID, newTaskRun)
	if err != nil {
		return err
	}

	// Create TaskCommand
	taskCommandRequest := task.TaskCommandRequest{
		TaskRunID:     tR.ID,
		TaskType:      tD.TaskType,
		Timeout:       tD.Timeout,
		AttemptNumber: tR.RetryCount + 1,
	}

	value, err := json.Marshal(taskCommandRequest)
	if err != nil {
		return err
	}

	err = taskScheduler.producer.SendMessage(
		ctx,
		TaskCommandRequestTopic,
		TaskCommandRequestTopic,
		value,
	)
	if err != nil {
		return err
	}

	slog.Info("ScheduleTask", "taskRunID", tR.ID, "attemptNumber", tR.RetryCount+1, "Name", tD.Name)

	return nil
}
