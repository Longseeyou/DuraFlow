package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
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

		taskRuns, err := taskScheduler.taskPoller.PollTaskRun(ctx)
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
	newTaskRun := map[string]any{"scheduled_at": time.Now()}
	if taskRun.Status == task.TASK_RUN_PENDING {
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
	}

	tR, err := taskScheduler.repository.UpdateTaskRunByID(ctx, taskRun.ID, newTaskRun)
	if err != nil {
		return err
	}

	// Create TaskAttempt
	taskAttempt, err := taskScheduler.repository.CreateTaskAttempt(
		ctx,
		task.TaskAttempt{
			TaskRunID:     taskRun.ID,
			AttemptNumber: taskRun.RetryCount,
			Status:        task.TASK_ATTEMPT_QUEUED,
		},
	)
	if err != nil {
		return err
	}

	// Create TaskCommand
	taskCommandRequest := task.TaskCommandRequest{
		TaskRunID:     taskRun.ID,
		TaskAttemptID: taskAttempt.ID,
		TaskType:      tD.TaskType,
		Input:         tR.Input,
		Timeout:       &tD.Timeout,
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

	return nil
}
