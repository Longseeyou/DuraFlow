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
	TaskCommandRequestTopic = "TaskCommandRequest"

	defaultRescuePollInterval = 5 * time.Second
)

type TaskScheduler struct {
	repository         OrchestratorRepository
	taskPoller         TaskPoller
	rescueTaskPoller   TaskPoller
	producer           message.Producer
	rescuePollInterval time.Duration
}

func NewTaskScheduler(
	repository OrchestratorRepository,
	taskPoller TaskPoller,
	rescueTaskPoller TaskPoller,
	producer message.Producer,
) *TaskScheduler {
	return &TaskScheduler{
		repository:         repository,
		taskPoller:         taskPoller,
		rescueTaskPoller:   rescueTaskPoller,
		producer:           producer,
		rescuePollInterval: defaultRescuePollInterval,
	}
}

func (taskScheduler *TaskScheduler) Run(ctx context.Context) {
	go taskScheduler.runRescuePoller(ctx)

	for {
		taskRuns, err := taskScheduler.taskPoller.PollTaskRun(ctx, 1)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("TaskScheduler PollTaskRun", "error", err)
			continue
		}

		taskScheduler.scheduleTaskRuns(ctx, taskRuns)
	}
}

func (taskScheduler *TaskScheduler) runRescuePoller(ctx context.Context) {
	ticker := time.NewTicker(taskScheduler.rescuePollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		taskRuns, err := taskScheduler.rescueTaskPoller.PollTaskRun(ctx, 10)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			slog.Error("TaskScheduler rescue PollTaskRun", "error", err)
			continue
		}

		taskScheduler.scheduleTaskRuns(ctx, taskRuns)
	}
}

func (taskScheduler *TaskScheduler) scheduleTaskRuns(ctx context.Context, taskRuns []task.TaskRun) {
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

	// Expose the executing task's own name so executors can self-identify
	// (e.g. a DDIM trainer derives its shard/round from "train-r1-s2-4").
	input["__self"] = tD.Name

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
