package orchestrator

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

var (
	ErrInvalidStateTransition = errors.New("invalid state transition")
	ErrMismatchedTaskEvent    = errors.New("task event does not match persisted task")
	ErrMismatchedTaskAttempt  = errors.New("task event does not match current attempt")
	ErrInvalidWorkflowGraph   = errors.New("invalid workflow dependency graph")
	ErrUnsupportedTaskEvent   = errors.New("unsupported task event")
)

type Option func(*Service)

func WithClock(now func() time.Time) Option {
	return func(service *Service) {
		service.now = now
	}
}

type Service struct {
	workflows  WorkflowRepository
	tasks      TaskRepository
	publisher  TaskPublisher
	now        func() time.Time
	runLocksMu sync.Mutex
	runLocks   map[uuid.UUID]*runLock
}

type runLock struct {
	mutex sync.Mutex
	refs  int
}

func NewService(
	workflows WorkflowRepository,
	tasks TaskRepository,
	publisher TaskPublisher,
	options ...Option,
) *Service {
	service := &Service{
		workflows: workflows,
		tasks:     tasks,
		publisher: publisher,
		now:       time.Now,
		runLocks:  make(map[uuid.UUID]*runLock),
	}
	for _, option := range options {
		option(service)
	}
	return service
}

// ScheduleWorkflow starts a pending workflow and queues every task whose
// dependencies have completed. Calling it again is safe: already queued,
// running, and completed tasks are not republished.
func (service *Service) ScheduleWorkflow(ctx context.Context, workflowRunID uuid.UUID) error {
	unlock := service.lockWorkflowRun(workflowRunID)
	defer unlock()
	return service.scheduleWorkflow(ctx, workflowRunID)
}

// ValidateWorkflowDefinition checks the workflow definition's task graph
// before any workflow run is started or task command is published.
func (service *Service) ValidateWorkflowDefinition(
	ctx context.Context,
	workflowDefinitionID uuid.UUID,
) error {
	if workflowDefinitionID == uuid.Nil {
		return fmt.Errorf("%w: workflow definition ID is required", ErrInvalidWorkflowGraph)
	}
	if _, err := service.workflows.GetWorkflowDefinitionById(ctx, workflowDefinitionID); err != nil {
		return fmt.Errorf("get workflow definition: %w", err)
	}
	definitions, dependencies, err := service.loadWorkflowDefinitionGraph(ctx, workflowDefinitionID)
	if err != nil {
		return err
	}
	if _, err := validateDependencyGraph(workflowDefinitionID, definitions, dependencies); err != nil {
		return err
	}
	return nil
}

// CancelWorkflow stops further scheduling, cancels every non-terminal task,
// and marks the workflow run cancelled. Workers still need their own
// cooperative cancellation mechanism for work already executing.
func (service *Service) CancelWorkflow(
	ctx context.Context,
	workflowRunID uuid.UUID,
	reason string,
) error {
	unlock := service.lockWorkflowRun(workflowRunID)
	defer unlock()

	workflowRun, err := service.workflows.GetWorkflowRunById(ctx, workflowRunID)
	if err != nil {
		return fmt.Errorf("get workflow run: %w", err)
	}
	if workflowRun.Status == workflow.CANCELLED {
		return nil
	}
	if workflowRun.Status == workflow.COMPLETED || workflowRun.Status == workflow.FAILED {
		return fmt.Errorf(
			"%w: workflow %s cannot be cancelled from %v",
			ErrInvalidStateTransition,
			workflowRunID,
			workflowRun.Status,
		)
	}

	taskRuns, err := service.tasks.GetTaskRunsByWorkflowRunId(ctx, workflowRunID)
	if err != nil {
		return fmt.Errorf("list task runs for cancellation: %w", err)
	}
	now := service.now().UTC()
	var reasonPointer *string
	if reason != "" {
		reasonPointer = &reason
	}

	for _, taskRun := range taskRuns {
		switch taskRun.Status {
		case task.TASK_RUN_COMPLETED, task.TASK_RUN_CANCELLED, task.TASK_RUN_DEAD_LETTERED:
			continue
		case task.TASK_RUN_QUEUED, task.TASK_RUN_RUNNING:
			attempt, err := service.currentAttempt(ctx, taskRun)
			if err != nil {
				return err
			}
			if _, err := service.tasks.UpdateTaskAttemptById(ctx, attempt.ID, map[string]any{
				"status":   task.TASK_ATTEMPT_CANCELLED,
				"ended_at": now,
				"log":      reasonPointer,
			}); err != nil {
				return fmt.Errorf("cancel task attempt: %w", err)
			}
		}

		if _, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
			"status":   task.TASK_RUN_CANCELLED,
			"ended_at": now,
		}); err != nil {
			return fmt.Errorf("cancel task run: %w", err)
		}
	}

	if _, err := service.workflows.UpdateWorkflowRunByIdInternal(ctx, workflowRunID, map[string]any{
		"status":       workflow.CANCELLED,
		"cancelled_at": now,
	}); err != nil {
		return fmt.Errorf("cancel workflow run: %w", err)
	}
	return nil
}

// HandleTaskEvent applies one worker event, records its attempt, schedules any
// newly unblocked tasks, and resolves the workflow's final state.
func (service *Service) HandleTaskEvent(ctx context.Context, event task.TaskEvent) error {
	if event.ID == uuid.Nil {
		return errors.New("task event ID is required")
	}
	if event.WorkflowRunID == uuid.Nil || event.TaskRunID == uuid.Nil {
		return errors.New("workflow run ID and task run ID are required")
	}

	unlock := service.lockWorkflowRun(event.WorkflowRunID)
	defer unlock()

	processed, err := service.tasks.TaskEventExists(ctx, event.ID)
	if err != nil {
		return fmt.Errorf("check task event: %w", err)
	}
	if processed {
		return nil
	}

	workflowRun, err := service.workflows.GetWorkflowRunById(ctx, event.WorkflowRunID)
	if err != nil {
		return fmt.Errorf("get workflow run: %w", err)
	}

	taskRun, err := service.tasks.GetTaskRunById(ctx, event.TaskRunID)
	if err != nil {
		return fmt.Errorf("get task run: %w", err)
	}
	if taskRun.WorkflowRunID != event.WorkflowRunID {
		return ErrMismatchedTaskEvent
	}
	currentAttempt := taskRun.RetryCount + 1
	if event.Attempt > 0 && event.Attempt != currentAttempt {
		previousFailureAlreadyApplied := event.EventType == task.TaskFailed &&
			event.Attempt < currentAttempt &&
			(taskRun.Status == task.TASK_RUN_QUEUED ||
				taskRun.Status == task.TASK_RUN_DEAD_LETTERED)
		if !previousFailureAlreadyApplied {
			return fmt.Errorf(
				"%w: event attempt %d, current attempt %d",
				ErrMismatchedTaskAttempt,
				event.Attempt,
				currentAttempt,
			)
		}
	}
	if isTerminalWorkflow(workflowRun.Status) {
		if eventEffectIsPersisted(event.EventType, taskRun.Status, workflowRun.Status) {
			if _, err := service.tasks.CreateTaskEvent(ctx, event); err != nil {
				return fmt.Errorf("record task event: %w", err)
			}
			return nil
		}
		return fmt.Errorf("%w: workflow is %v", ErrInvalidStateTransition, workflowRun.Status)
	}

	// A failure may have moved the task to its next attempt before the event
	// record was persisted. Redelivery must finish idempotency bookkeeping,
	// not spend another retry.
	if event.EventType == task.TaskFailed &&
		event.Attempt > 0 &&
		event.Attempt < currentAttempt &&
		(taskRun.Status == task.TASK_RUN_QUEUED ||
			taskRun.Status == task.TASK_RUN_DEAD_LETTERED) {
		if _, err := service.tasks.CreateTaskEvent(ctx, event); err != nil {
			return fmt.Errorf("record task event: %w", err)
		}
		return nil
	}

	switch event.EventType {
	case task.TaskStarted:
		err = service.handleTaskStarted(ctx, taskRun, event)
	case task.TaskCompleted:
		err = service.handleTaskCompleted(ctx, taskRun, event)
	case task.TaskFailed:
		err = service.handleTaskFailed(ctx, taskRun, event)
	case task.TaskDeadLettered:
		err = service.handleTaskDeadLettered(ctx, taskRun, event)
	case task.TaskCancelled:
		err = service.handleTaskCancelled(ctx, taskRun, event)
	default:
		err = fmt.Errorf("%w: %s", ErrUnsupportedTaskEvent, event.EventType)
	}
	if err != nil {
		return err
	}

	if _, err := service.tasks.CreateTaskEvent(ctx, event); err != nil {
		return fmt.Errorf("record task event: %w", err)
	}
	return nil
}

func (service *Service) handleTaskDeadLettered(
	ctx context.Context,
	taskRun task.TaskRun,
	event task.TaskEvent,
) error {
	if taskRun.Status == task.TASK_RUN_DEAD_LETTERED {
		return nil
	}
	if taskRun.Status != task.TASK_RUN_RUNNING && taskRun.Status != task.TASK_RUN_FAILED {
		return fmt.Errorf(
			"%w: task %s cannot be dead-lettered from %v",
			ErrInvalidStateTransition,
			taskRun.ID,
			taskRun.Status,
		)
	}

	now := event.CreatedAt
	if now.IsZero() {
		now = service.now().UTC()
	}
	if taskRun.Status == task.TASK_RUN_RUNNING {
		if err := service.finishCurrentAttempt(ctx, taskRun, task.TASK_ATTEMPT_DEAD_LETTERED, now, event); err != nil {
			return err
		}
	}
	if _, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
		"status":   task.TASK_RUN_DEAD_LETTERED,
		"ended_at": now,
	}); err != nil {
		return fmt.Errorf("dead-letter task run: %w", err)
	}
	if _, err := service.workflows.UpdateWorkflowRunByIdInternal(ctx, taskRun.WorkflowRunID, map[string]any{
		"status":       workflow.FAILED,
		"completed_at": now,
	}); err != nil {
		return fmt.Errorf("fail workflow run: %w", err)
	}
	return nil
}

func (service *Service) scheduleWorkflow(ctx context.Context, workflowRunID uuid.UUID) error {
	workflowRun, err := service.workflows.GetWorkflowRunById(ctx, workflowRunID)
	if err != nil {
		return fmt.Errorf("get workflow run: %w", err)
	}
	if isTerminalWorkflow(workflowRun.Status) {
		return fmt.Errorf("%w: workflow is %v", ErrInvalidStateTransition, workflowRun.Status)
	}

	definitions, dependencies, err := service.loadWorkflowDefinitionGraph(
		ctx,
		workflowRun.WorkflowDefinitionID,
	)
	if err != nil {
		return err
	}
	definitionsByID, err := validateDependencyGraph(
		workflowRun.WorkflowDefinitionID,
		definitions,
		dependencies,
	)
	if err != nil {
		return err
	}

	taskRuns, err := service.tasks.GetTaskRunsByWorkflowRunId(ctx, workflowRunID)
	if err != nil {
		return fmt.Errorf("list task runs: %w", err)
	}

	now := service.now().UTC()
	if len(taskRuns) == 0 {
		if len(definitionsByID) > 0 {
			return fmt.Errorf(
				"%w: workflow definition %s has tasks but workflow run %s has no task runs",
				ErrInvalidWorkflowGraph,
				workflowRun.WorkflowDefinitionID,
				workflowRunID,
			)
		}
		updates := map[string]any{
			"status":       workflow.COMPLETED,
			"completed_at": now,
		}
		if workflowRun.Status == workflow.PENDING {
			updates["started_at"] = now
		}
		_, err := service.workflows.UpdateWorkflowRunByIdInternal(ctx, workflowRunID, updates)
		if err != nil {
			return fmt.Errorf("complete empty workflow run: %w", err)
		}
		return nil
	}

	byDefinition, err := validateWorkflowRunMaterialization(
		workflowRunID,
		taskRuns,
		definitionsByID,
	)
	if err != nil {
		return err
	}
	if workflowRun.Status == workflow.PENDING {
		if _, err := service.workflows.UpdateWorkflowRunByIdInternal(ctx, workflowRunID, map[string]any{
			"status":     workflow.RUNNING,
			"started_at": now,
		}); err != nil {
			return fmt.Errorf("start workflow run: %w", err)
		}
	}

	required := make(map[uuid.UUID][]uuid.UUID)
	for _, dependency := range dependencies {
		required[dependency.TaskID] = append(required[dependency.TaskID], dependency.DependOnTaskID)
	}

	for _, taskRun := range taskRuns {
		if taskRun.Status != task.TASK_RUN_PENDING {
			continue
		}
		if !dependenciesCompleted(required[taskRun.TaskDefinitionID], byDefinition) {
			continue
		}
		if err := service.queueTask(ctx, taskRun, taskRun); err != nil {
			return err
		}
	}
	return service.resolveWorkflow(ctx, workflowRunID)
}

func (service *Service) queueTask(
	ctx context.Context,
	taskRun task.TaskRun,
	rollback task.TaskRun,
) error {
	if taskRun.Status != task.TASK_RUN_PENDING && taskRun.Status != task.TASK_RUN_FAILED {
		return fmt.Errorf("%w: task %s cannot be queued from %v", ErrInvalidStateTransition, taskRun.ID, taskRun.Status)
	}

	definition := taskRun.TaskDefinition
	if definition.ID == uuid.Nil {
		var err error
		definition, err = service.tasks.GetTaskDefinitionById(ctx, taskRun.TaskDefinitionID)
		if err != nil {
			return fmt.Errorf("get task definition: %w", err)
		}
	}

	now := service.now().UTC()
	attemptNumber := taskRun.RetryCount + 1
	if _, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
		"status":       task.TASK_RUN_QUEUED,
		"retry_count":  taskRun.RetryCount,
		"scheduled_at": now,
		"started_at":   nil,
		"ended_at":     nil,
	}); err != nil {
		return fmt.Errorf("queue task run: %w", err)
	}

	attempt, err := service.tasks.CreateTaskAttempt(ctx, task.TaskAttempt{
		TaskRunID:     taskRun.ID,
		AttemptNumber: attemptNumber,
		Status:        task.TASK_ATTEMPT_QUEUED,
	})
	if err != nil {
		return errors.Join(
			fmt.Errorf("create task attempt: %w", err),
			service.restoreTaskRun(ctx, rollback),
		)
	}

	command := task.Command{
		WorkflowRunID: taskRun.WorkflowRunID,
		TaskRunID:     taskRun.ID,
		TaskKey:       definition.Name,
		TaskType:      definition.TaskType,
		Payload:       taskRun.Input,
		Attempt:       attemptNumber,
		CreatedAt:     now,
	}
	if err := service.publisher.PublishTask(ctx, command); err != nil {
		return errors.Join(
			fmt.Errorf("publish task command: %w", err),
			wrapIfError(
				"delete unpublished task attempt",
				service.tasks.DeleteTaskAttemptById(ctx, attempt.ID),
			),
			service.restoreTaskRun(ctx, rollback),
		)
	}
	return nil
}

func (service *Service) restoreTaskRun(ctx context.Context, taskRun task.TaskRun) error {
	_, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
		"status":       taskRun.Status,
		"retry_count":  taskRun.RetryCount,
		"scheduled_at": taskRun.ScheduledAt,
		"started_at":   taskRun.StartedAt,
		"ended_at":     taskRun.EndedAt,
		"output":       taskRun.Output,
	})
	return wrapIfError("restore task run after dispatch failure", err)
}

func (service *Service) handleTaskStarted(
	ctx context.Context,
	taskRun task.TaskRun,
	event task.TaskEvent,
) error {
	if taskRun.Status == task.TASK_RUN_RUNNING {
		return nil
	}
	if taskRun.Status != task.TASK_RUN_QUEUED {
		return fmt.Errorf("%w: task %s cannot start from %v", ErrInvalidStateTransition, taskRun.ID, taskRun.Status)
	}

	now := event.CreatedAt
	if now.IsZero() {
		now = service.now().UTC()
	}
	attempt, err := service.currentAttempt(ctx, taskRun)
	if err != nil {
		return err
	}
	if _, err := service.tasks.UpdateTaskAttemptById(ctx, attempt.ID, map[string]any{
		"status":     task.TASK_ATTEMPT_RUNNING,
		"worker_id":  event.WorkerID,
		"started_at": now,
	}); err != nil {
		return fmt.Errorf("start task attempt: %w", err)
	}
	if _, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
		"status":     task.TASK_RUN_RUNNING,
		"started_at": now,
	}); err != nil {
		return fmt.Errorf("mark task running: %w", err)
	}
	return nil
}

func (service *Service) handleTaskCompleted(
	ctx context.Context,
	taskRun task.TaskRun,
	event task.TaskEvent,
) error {
	if taskRun.Status != task.TASK_RUN_RUNNING && taskRun.Status != task.TASK_RUN_COMPLETED {
		return fmt.Errorf("%w: task %s cannot complete from %v", ErrInvalidStateTransition, taskRun.ID, taskRun.Status)
	}

	if taskRun.Status != task.TASK_RUN_COMPLETED {
		now := event.CreatedAt
		if now.IsZero() {
			now = service.now().UTC()
		}
		if err := service.finishCurrentAttempt(ctx, taskRun, task.TASK_ATTEMPT_COMPLETED, now, event); err != nil {
			return err
		}
		if _, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
			"status":   task.TASK_RUN_COMPLETED,
			"output":   event.Result,
			"ended_at": now,
		}); err != nil {
			return fmt.Errorf("complete task run: %w", err)
		}
	}

	return service.scheduleWorkflow(ctx, taskRun.WorkflowRunID)
}

func (service *Service) handleTaskFailed(
	ctx context.Context,
	taskRun task.TaskRun,
	event task.TaskEvent,
) error {
	if taskRun.Status != task.TASK_RUN_RUNNING &&
		taskRun.Status != task.TASK_RUN_FAILED &&
		taskRun.Status != task.TASK_RUN_DEAD_LETTERED {
		return fmt.Errorf("%w: task %s cannot fail from %v", ErrInvalidStateTransition, taskRun.ID, taskRun.Status)
	}
	if taskRun.Status == task.TASK_RUN_DEAD_LETTERED {
		return nil
	}

	now := event.CreatedAt
	if now.IsZero() {
		now = service.now().UTC()
	}
	if taskRun.Status == task.TASK_RUN_RUNNING {
		if err := service.finishCurrentAttempt(ctx, taskRun, task.TASK_ATTEMPT_FAILED, now, event); err != nil {
			return err
		}
		if _, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
			"status":   task.TASK_RUN_FAILED,
			"ended_at": now,
		}); err != nil {
			return fmt.Errorf("fail task run: %w", err)
		}
		taskRun.Status = task.TASK_RUN_FAILED
		taskRun.EndedAt = &now
	}

	if taskRun.RetryCount < taskRun.MaxRetries {
		rollback := taskRun
		retryCount := taskRun.RetryCount + 1
		updated, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
			"status":      task.TASK_RUN_FAILED,
			"retry_count": retryCount,
		})
		if err != nil {
			return fmt.Errorf("schedule task retry: %w", err)
		}
		updated.RetryCount = retryCount
		updated.Status = task.TASK_RUN_FAILED
		if updated.WorkflowRunID == uuid.Nil {
			updated.WorkflowRunID = taskRun.WorkflowRunID
		}
		if updated.TaskDefinitionID == uuid.Nil {
			updated.TaskDefinitionID = taskRun.TaskDefinitionID
			updated.TaskDefinition = taskRun.TaskDefinition
			updated.Input = taskRun.Input
		}
		return service.queueTask(ctx, updated, rollback)
	}

	if _, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
		"status":   task.TASK_RUN_DEAD_LETTERED,
		"ended_at": now,
	}); err != nil {
		return fmt.Errorf("dead-letter task run: %w", err)
	}
	if _, err := service.workflows.UpdateWorkflowRunByIdInternal(ctx, taskRun.WorkflowRunID, map[string]any{
		"status":       workflow.FAILED,
		"completed_at": now,
	}); err != nil {
		return fmt.Errorf("fail workflow run: %w", err)
	}
	return nil
}

func (service *Service) handleTaskCancelled(
	ctx context.Context,
	taskRun task.TaskRun,
	event task.TaskEvent,
) error {
	if taskRun.Status != task.TASK_RUN_PENDING &&
		taskRun.Status != task.TASK_RUN_QUEUED &&
		taskRun.Status != task.TASK_RUN_RUNNING &&
		taskRun.Status != task.TASK_RUN_CANCELLED {
		return fmt.Errorf("%w: task %s cannot be cancelled from %v", ErrInvalidStateTransition, taskRun.ID, taskRun.Status)
	}
	now := event.CreatedAt
	if now.IsZero() {
		now = service.now().UTC()
	}
	if taskRun.Status != task.TASK_RUN_CANCELLED {
		if taskRun.Status == task.TASK_RUN_QUEUED || taskRun.Status == task.TASK_RUN_RUNNING {
			if err := service.finishCurrentAttempt(ctx, taskRun, task.TASK_ATTEMPT_CANCELLED, now, event); err != nil {
				return err
			}
		}
		if _, err := service.tasks.UpdateTaskRunById(ctx, taskRun.ID, map[string]any{
			"status":   task.TASK_RUN_CANCELLED,
			"ended_at": now,
		}); err != nil {
			return fmt.Errorf("cancel task run: %w", err)
		}
	}
	if _, err := service.workflows.UpdateWorkflowRunByIdInternal(ctx, taskRun.WorkflowRunID, map[string]any{
		"status":       workflow.FAILED,
		"completed_at": now,
	}); err != nil {
		return fmt.Errorf("fail workflow after task cancellation: %w", err)
	}
	return nil
}

func (service *Service) resolveWorkflow(ctx context.Context, workflowRunID uuid.UUID) error {
	taskRuns, err := service.tasks.GetTaskRunsByWorkflowRunId(ctx, workflowRunID)
	if err != nil {
		return fmt.Errorf("list task runs for completion: %w", err)
	}
	allCompleted := len(taskRuns) > 0
	for _, taskRun := range taskRuns {
		if taskRun.Status == task.TASK_RUN_DEAD_LETTERED {
			_, err := service.workflows.UpdateWorkflowRunByIdInternal(ctx, workflowRunID, map[string]any{
				"status":       workflow.FAILED,
				"completed_at": service.now().UTC(),
			})
			return err
		}
		if taskRun.Status != task.TASK_RUN_COMPLETED {
			allCompleted = false
		}
	}
	if !allCompleted {
		return nil
	}
	if _, err := service.workflows.UpdateWorkflowRunByIdInternal(ctx, workflowRunID, map[string]any{
		"status":       workflow.COMPLETED,
		"completed_at": service.now().UTC(),
	}); err != nil {
		return fmt.Errorf("complete workflow run: %w", err)
	}
	return nil
}

func (service *Service) currentAttempt(
	ctx context.Context,
	taskRun task.TaskRun,
) (task.TaskAttempt, error) {
	attempt, err := service.tasks.GetTaskAttemptByTaskRunAndNumber(
		ctx,
		taskRun.ID,
		taskRun.RetryCount+1,
	)
	if err != nil {
		return task.TaskAttempt{}, fmt.Errorf("get current task attempt: %w", err)
	}
	return attempt, nil
}

func (service *Service) finishCurrentAttempt(
	ctx context.Context,
	taskRun task.TaskRun,
	status task.TaskAttemptStatus,
	endedAt time.Time,
	event task.TaskEvent,
) error {
	attempt, err := service.currentAttempt(ctx, taskRun)
	if err != nil {
		return err
	}
	updates := map[string]any{
		"status":   status,
		"ended_at": endedAt,
	}
	if event.WorkerID != "" {
		updates["worker_id"] = event.WorkerID
	}
	if event.Error != nil {
		updates["log"] = event.Error
	}
	if _, err := service.tasks.UpdateTaskAttemptById(ctx, attempt.ID, updates); err != nil {
		return fmt.Errorf("finish task attempt: %w", err)
	}
	return nil
}

func (service *Service) lockWorkflowRun(workflowRunID uuid.UUID) func() {
	service.runLocksMu.Lock()
	lock := service.runLocks[workflowRunID]
	if lock == nil {
		lock = &runLock{}
		service.runLocks[workflowRunID] = lock
	}
	lock.refs++
	service.runLocksMu.Unlock()

	lock.mutex.Lock()
	return func() {
		lock.mutex.Unlock()
		service.runLocksMu.Lock()
		lock.refs--
		if lock.refs == 0 {
			delete(service.runLocks, workflowRunID)
		}
		service.runLocksMu.Unlock()
	}
}

func (service *Service) loadWorkflowDefinitionGraph(
	ctx context.Context,
	workflowDefinitionID uuid.UUID,
) ([]task.TaskDefinition, []task.TaskDependency, error) {
	definitions, err := service.tasks.GetTaskDefinitionsByWorkflowDefinitionId(
		ctx,
		workflowDefinitionID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("list task definitions: %w", err)
	}
	dependencies, err := service.tasks.GetTaskDependenciesByWorkflowDefinitionId(
		ctx,
		workflowDefinitionID,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("list task dependencies: %w", err)
	}
	return definitions, dependencies, nil
}

func dependenciesCompleted(
	dependencyIDs []uuid.UUID,
	taskRuns map[uuid.UUID]task.TaskRun,
) bool {
	for _, dependencyID := range dependencyIDs {
		dependency, ok := taskRuns[dependencyID]
		if !ok || dependency.Status != task.TASK_RUN_COMPLETED {
			return false
		}
	}
	return true
}

func validateDependencyGraph(
	workflowDefinitionID uuid.UUID,
	definitions []task.TaskDefinition,
	dependencies []task.TaskDependency,
) (map[uuid.UUID]task.TaskDefinition, error) {
	definitionsByID := make(map[uuid.UUID]task.TaskDefinition, len(definitions))
	graph := make(map[uuid.UUID][]uuid.UUID, len(definitions))
	for _, definition := range definitions {
		if definition.ID == uuid.Nil {
			return nil, fmt.Errorf(
				"%w: workflow definition %s has a task definition without an ID",
				ErrInvalidWorkflowGraph,
				workflowDefinitionID,
			)
		}
		if definition.WorkflowDefinitionID != workflowDefinitionID {
			return nil, fmt.Errorf(
				"%w: task definition %s belongs to workflow definition %s, not %s",
				ErrInvalidWorkflowGraph,
				definition.ID,
				definition.WorkflowDefinitionID,
				workflowDefinitionID,
			)
		}
		if _, exists := definitionsByID[definition.ID]; exists {
			return nil, fmt.Errorf(
				"%w: duplicate task definition %s",
				ErrInvalidWorkflowGraph,
				definition.ID,
			)
		}
		definitionsByID[definition.ID] = definition
		graph[definition.ID] = nil
	}
	for _, dependency := range dependencies {
		if _, exists := definitionsByID[dependency.TaskID]; !exists {
			return nil, fmt.Errorf(
				"%w: dependent task definition %s is not in workflow definition %s",
				ErrInvalidWorkflowGraph,
				dependency.TaskID,
				workflowDefinitionID,
			)
		}
		if _, exists := definitionsByID[dependency.DependOnTaskID]; !exists {
			return nil, fmt.Errorf(
				"%w: prerequisite task definition %s is not in workflow definition %s",
				ErrInvalidWorkflowGraph,
				dependency.DependOnTaskID,
				workflowDefinitionID,
			)
		}
		if dependency.TaskID == dependency.DependOnTaskID {
			return nil, fmt.Errorf(
				"%w: task definition %s depends on itself",
				ErrInvalidWorkflowGraph,
				dependency.TaskID,
			)
		}
		graph[dependency.TaskID] = append(graph[dependency.TaskID], dependency.DependOnTaskID)
	}

	const (
		unvisited = iota
		visiting
		visited
	)
	state := make(map[uuid.UUID]int, len(graph))
	var visit func(uuid.UUID) error
	visit = func(definitionID uuid.UUID) error {
		switch state[definitionID] {
		case visiting:
			return fmt.Errorf(
				"%w: dependency cycle includes task definition %s",
				ErrInvalidWorkflowGraph,
				definitionID,
			)
		case visited:
			return nil
		}
		state[definitionID] = visiting
		for _, prerequisiteID := range graph[definitionID] {
			if err := visit(prerequisiteID); err != nil {
				return err
			}
		}
		state[definitionID] = visited
		return nil
	}
	for definitionID := range graph {
		if err := visit(definitionID); err != nil {
			return nil, err
		}
	}
	return definitionsByID, nil
}

func validateWorkflowRunMaterialization(
	workflowRunID uuid.UUID,
	taskRuns []task.TaskRun,
	definitionsByID map[uuid.UUID]task.TaskDefinition,
) (map[uuid.UUID]task.TaskRun, error) {
	byDefinition := make(map[uuid.UUID]task.TaskRun, len(taskRuns))
	for _, taskRun := range taskRuns {
		if taskRun.WorkflowRunID != workflowRunID {
			return nil, fmt.Errorf("task run %s belongs to another workflow run", taskRun.ID)
		}
		if taskRun.TaskDefinitionID == uuid.Nil {
			return nil, fmt.Errorf("%w: task run %s has no task definition", ErrInvalidWorkflowGraph, taskRun.ID)
		}
		if _, exists := definitionsByID[taskRun.TaskDefinitionID]; !exists {
			return nil, fmt.Errorf(
				"%w: task run %s references task definition %s outside this workflow definition",
				ErrInvalidWorkflowGraph,
				taskRun.ID,
				taskRun.TaskDefinitionID,
			)
		}
		if existing, exists := byDefinition[taskRun.TaskDefinitionID]; exists {
			return nil, fmt.Errorf(
				"%w: task runs %s and %s share definition %s",
				ErrInvalidWorkflowGraph,
				existing.ID,
				taskRun.ID,
				taskRun.TaskDefinitionID,
			)
		}
		byDefinition[taskRun.TaskDefinitionID] = taskRun
	}
	for definitionID := range definitionsByID {
		if _, exists := byDefinition[definitionID]; !exists {
			return nil, fmt.Errorf(
				"%w: workflow run %s has no task run for task definition %s",
				ErrInvalidWorkflowGraph,
				workflowRunID,
				definitionID,
			)
		}
	}
	return byDefinition, nil
}

func isTerminalWorkflow(status workflow.WorkflowRunStatus) bool {
	return status == workflow.COMPLETED || status == workflow.FAILED || status == workflow.CANCELLED
}

func eventEffectIsPersisted(
	eventType task.TaskEventType,
	taskStatus task.TaskRunStatus,
	workflowStatus workflow.WorkflowRunStatus,
) bool {
	switch eventType {
	case task.TaskCompleted:
		return taskStatus == task.TASK_RUN_COMPLETED && workflowStatus == workflow.COMPLETED
	case task.TaskFailed, task.TaskDeadLettered:
		return taskStatus == task.TASK_RUN_DEAD_LETTERED && workflowStatus == workflow.FAILED
	case task.TaskCancelled:
		return taskStatus == task.TASK_RUN_CANCELLED
	default:
		return false
	}
}

func wrapIfError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", operation, err)
}
