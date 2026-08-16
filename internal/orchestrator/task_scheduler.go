package orchestrator

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
)

type TaskScheduler struct {
	orchestratorRepository OrchestratorRepository
	taskPoller             TaskPoller
	producer               message.Producer
}

func NewTaskScheduler(taskPoller TaskPoller, producer message.Producer) *TaskScheduler {
	return &TaskScheduler{taskPoller: taskPoller, producer: producer}
}

func (taskScheduler *TaskScheduler) Run(ctx context.Context) {
	for {
		taskRuns, err := taskScheduler.taskPoller.PollTaskRun(ctx)
		if err != nil {

		}

		for _, taskRun := range taskRuns {
			err := taskScheduler.ScheduleTask(ctx, &taskRun)
			if err != nil {

			}
		}
	}
}

func (taskScheduler *TaskScheduler) ScheduleTask(ctx context.Context, taskRun *task.TaskRun) error {
	// Idempotency
	ok, err := taskScheduler.orchestratorRepository.TaskRunIdempotency(
		ctx,
		taskRun.ID,
		task.TASK_RUN_QUEUED,
	)
	if err != nil {

	}
	if !ok {

	}

	// Update TaskRun
	input := map[string]string{}
	if taskRun.Status == task.TASK_RUN_PENDING {
		tRs, err := taskScheduler.orchestratorRepository.GetPredecessorTaskRuns(ctx, taskRun.ID)
		if err != nil {

		}
		for _, tR := range tRs {
			input[tR.TaskDefinition.Name] = *tR.Output
		}
	}
	input_json, err := json.Marshal(input)
	if err != nil {
	}

	_, err = taskScheduler.orchestratorRepository.UpdateTaskRunById(
		ctx,
		taskRun.ID,
		map[string]any{"ScheduledAt": time.Now(), "Input": input_json},
	)
	if err != nil {

	}

	// Create TaskAttempt

	taskAttempt, err := taskScheduler.orchestratorRepository.CreateTaskAttempt(
		ctx,
		task.TaskAttempt{TaskRunID: taskRun.ID, AttemptNumber: taskRun.RetryCount},
	)
	if err != nil {

	}

	// Create TaskCommand
	taskCommandRequest := task.TaskCommandRequest{
		TaskAttemptID: taskAttempt.ID,
		TaskType:      taskRun.TaskDefinition.TaskType,
		Input:         taskRun.Input,
		Timeout:       &taskRun.TaskDefinition.Timeout,
	}

	value, err := json.Marshal(taskCommandRequest)
	if err != nil {
	}

	err = taskScheduler.producer.SendMessage(
		ctx,
		"TaskCommandRequest",
		"TaskCommandRequest",
		value,
	)
	if err != nil {
	}

	return nil
}
