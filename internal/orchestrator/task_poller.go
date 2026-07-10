package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/google/uuid"
)

type TaskPollerOption func(*TaskPoller)

func WithTaskPollerClock(now func() time.Time) TaskPollerOption {
	return func(poller *TaskPoller) {
		if now != nil {
			poller.now = now
		}
	}
}

// TaskPoller owns task claiming, attempt creation, and command dispatch.
// The conditional task-run update makes the database the source of truth when
// multiple orchestrator instances discover the same runnable task.
type TaskPoller struct {
	tasks     TaskPollingRepository
	publisher TaskPublisher
	now       func() time.Time
}

func NewTaskPoller(
	tasks TaskPollingRepository,
	publisher TaskPublisher,
	options ...TaskPollerOption,
) *TaskPoller {
	poller := &TaskPoller{
		tasks:     tasks,
		publisher: publisher,
		now:       time.Now,
	}
	for _, option := range options {
		option(poller)
	}
	return poller
}

func (poller *TaskPoller) QueueTask(
	ctx context.Context,
	taskRun task.TaskRun,
	rollback task.TaskRun,
) error {
	if taskRun.Status != task.TASK_RUN_PENDING && taskRun.Status != task.TASK_RUN_FAILED {
		return fmt.Errorf(
			"%w: task %s cannot be queued from %v",
			ErrInvalidStateTransition,
			taskRun.ID,
			taskRun.Status,
		)
	}
	if poller.tasks == nil {
		return errors.New("task polling repository is required")
	}
	if poller.publisher == nil {
		return errors.New("task publisher is required")
	}

	definition := taskRun.TaskDefinition
	if definition.ID == uuid.Nil {
		var err error
		definition, err = poller.tasks.GetTaskDefinitionById(ctx, taskRun.TaskDefinitionID)
		if err != nil {
			return fmt.Errorf("get task definition: %w", err)
		}
	}

	now := poller.now().UTC()
	attemptNumber := taskRun.RetryCount + 1
	queued, changed, err := poller.tasks.UpdateTaskRunByIdAndStatus(
		ctx,
		taskRun.ID,
		[]task.TaskRunStatus{task.TASK_RUN_PENDING, task.TASK_RUN_FAILED},
		map[string]any{
			"status":       task.TASK_RUN_QUEUED,
			"retry_count":  taskRun.RetryCount,
			"scheduled_at": now,
			"started_at":   nil,
			"ended_at":     nil,
		},
	)
	if err != nil {
		return fmt.Errorf("queue task run: %w", err)
	}
	if !changed {
		return nil
	}
	taskRun = mergeTaskRunSnapshot(queued, rollback)
	taskRun.TaskDefinition = definition

	attempt, err := poller.tasks.CreateTaskAttempt(ctx, task.TaskAttempt{
		TaskRunID:     taskRun.ID,
		AttemptNumber: attemptNumber,
		Status:        task.TASK_ATTEMPT_QUEUED,
	})
	if err != nil {
		return errors.Join(
			fmt.Errorf("create task attempt: %w", err),
			poller.restoreTaskRunIfQueued(ctx, rollback),
		)
	}

	command := task.Command{
		CommandID:     uuid.New(),
		WorkflowRunID: taskRun.WorkflowRunID,
		TaskRunID:     taskRun.ID,
		TaskKey:       definition.Name,
		TaskType:      definition.TaskType,
		Payload:       taskRun.Input,
		Attempt:       attemptNumber,
		CreatedAt:     now,
	}
	if err := poller.publisher.PublishTask(ctx, command); err != nil {
		return errors.Join(
			fmt.Errorf("publish task command: %w", err),
			poller.deleteTaskAttempt(ctx, attempt.ID),
			poller.restoreTaskRunIfQueued(ctx, rollback),
		)
	}
	return nil
}

func (poller *TaskPoller) deleteTaskAttempt(ctx context.Context, attemptID uuid.UUID) error {
	_, err := poller.tasks.HardDeleteTaskAttempt(ctx, attemptID)
	return wrapIfError("delete unpublished task attempt", err)
}

func (poller *TaskPoller) restoreTaskRunIfQueued(
	ctx context.Context,
	taskRun task.TaskRun,
) error {
	_, _, err := poller.tasks.UpdateTaskRunByIdAndStatus(
		ctx,
		taskRun.ID,
		[]task.TaskRunStatus{task.TASK_RUN_QUEUED},
		map[string]any{
			"status":       taskRun.Status,
			"retry_count":  taskRun.RetryCount,
			"scheduled_at": taskRun.ScheduledAt,
			"started_at":   taskRun.StartedAt,
			"ended_at":     taskRun.EndedAt,
			"output":       taskRun.Output,
		},
	)
	return wrapIfError("restore task run after dispatch failure", err)
}
