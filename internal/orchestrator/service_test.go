package orchestrator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
)

func TestScheduleWorkflowQueuesOnlyRunnableTasks(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)

	rootDefinition := task.TaskDefinition{Name: "root", TaskType: task.MOCK_TASK}
	rootDefinition.ID = uuid.New()
	childDefinition := task.TaskDefinition{Name: "child", TaskType: task.MOCK_TASK}
	childDefinition.ID = uuid.New()
	fixture.store.definitions[rootDefinition.ID] = rootDefinition
	fixture.store.definitions[childDefinition.ID] = childDefinition

	root := fixture.addTaskRun(rootDefinition, task.TASK_RUN_PENDING, 0, 2)
	child := fixture.addTaskRun(childDefinition, task.TASK_RUN_PENDING, 0, 2)
	fixture.store.dependencies = []task.TaskDependency{{
		TaskID:         childDefinition.ID,
		DependOnTaskID: rootDefinition.ID,
	}}

	if err := fixture.service.ScheduleWorkflow(context.Background(), fixture.workflowRun.ID); err != nil {
		t.Fatalf("ScheduleWorkflow() error = %v", err)
	}

	if got := fixture.store.workflowRuns[fixture.workflowRun.ID].Status; got != workflow.RUNNING {
		t.Fatalf("workflow status = %v, want RUNNING", got)
	}
	if got := fixture.store.taskRuns[root.ID].Status; got != task.TASK_RUN_QUEUED {
		t.Fatalf("root status = %v, want QUEUED", got)
	}
	if got := fixture.store.taskRuns[child.ID].Status; got != task.TASK_RUN_PENDING {
		t.Fatalf("child status = %v, want PENDING", got)
	}
	if len(fixture.publisher.commands) != 1 {
		t.Fatalf("published commands = %d, want 1", len(fixture.publisher.commands))
	}
	command := fixture.publisher.commands[0]
	if command.TaskRunID != root.ID || command.Attempt != 1 || command.TaskKey != "root" {
		t.Fatalf("unexpected command: %+v", command)
	}
	if _, ok := fixture.store.attempts[attemptKey{root.ID, 1}]; !ok {
		t.Fatal("initial task attempt was not created")
	}
}

func TestCompletedTaskUnblocksDependencyAndCompletesWorkflow(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)

	rootDefinition := task.TaskDefinition{Name: "root", TaskType: task.MOCK_TASK}
	rootDefinition.ID = uuid.New()
	childDefinition := task.TaskDefinition{Name: "child", TaskType: task.MOCK_TASK}
	childDefinition.ID = uuid.New()
	fixture.store.definitions[rootDefinition.ID] = rootDefinition
	fixture.store.definitions[childDefinition.ID] = childDefinition
	root := fixture.addTaskRun(rootDefinition, task.TASK_RUN_PENDING, 0, 1)
	child := fixture.addTaskRun(childDefinition, task.TASK_RUN_PENDING, 0, 1)
	fixture.store.dependencies = []task.TaskDependency{{
		TaskID:         childDefinition.ID,
		DependOnTaskID: rootDefinition.ID,
	}}

	ctx := context.Background()
	if err := fixture.service.ScheduleWorkflow(ctx, fixture.workflowRun.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.startTask(ctx, root.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.completeTask(ctx, root.ID, `{"root":true}`); err != nil {
		t.Fatal(err)
	}

	if got := fixture.store.taskRuns[child.ID].Status; got != task.TASK_RUN_QUEUED {
		t.Fatalf("child status = %v, want QUEUED", got)
	}
	if len(fixture.publisher.commands) != 2 {
		t.Fatalf("published commands = %d, want 2", len(fixture.publisher.commands))
	}

	if err := fixture.startTask(ctx, child.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.completeTask(ctx, child.ID, `{"child":true}`); err != nil {
		t.Fatal(err)
	}
	if got := fixture.store.workflowRuns[fixture.workflowRun.ID].Status; got != workflow.COMPLETED {
		t.Fatalf("workflow status = %v, want COMPLETED", got)
	}
}

func TestFailedTaskRetriesThenDeadLettersWorkflow(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)

	definition := task.TaskDefinition{Name: "flaky", TaskType: task.MOCK_TASK}
	definition.ID = uuid.New()
	fixture.store.definitions[definition.ID] = definition
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 1)

	ctx := context.Background()
	if err := fixture.service.ScheduleWorkflow(ctx, fixture.workflowRun.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.startTask(ctx, taskRun.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.failTask(ctx, taskRun.ID, "first failure"); err != nil {
		t.Fatal(err)
	}

	retried := fixture.store.taskRuns[taskRun.ID]
	if retried.Status != task.TASK_RUN_QUEUED || retried.RetryCount != 1 {
		t.Fatalf("retried task = %+v, want QUEUED with retry count 1", retried)
	}
	if len(fixture.publisher.commands) != 2 || fixture.publisher.commands[1].Attempt != 2 {
		t.Fatalf("retry commands = %+v", fixture.publisher.commands)
	}

	if err := fixture.startTask(ctx, taskRun.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.failTask(ctx, taskRun.ID, "second failure"); err != nil {
		t.Fatal(err)
	}
	if got := fixture.store.taskRuns[taskRun.ID].Status; got != task.TASK_RUN_DEAD_LETTERED {
		t.Fatalf("task status = %v, want DEAD_LETTERED", got)
	}
	if got := fixture.store.workflowRuns[fixture.workflowRun.ID].Status; got != workflow.FAILED {
		t.Fatalf("workflow status = %v, want FAILED", got)
	}
}

func TestHandleTaskEventIsIdempotentByEventID(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	definition := task.TaskDefinition{Name: "once", TaskType: task.MOCK_TASK}
	definition.ID = uuid.New()
	fixture.store.definitions[definition.ID] = definition
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 1)

	ctx := context.Background()
	if err := fixture.service.ScheduleWorkflow(ctx, fixture.workflowRun.ID); err != nil {
		t.Fatal(err)
	}
	event := task.TaskEvent{
		WorkflowRunID: fixture.workflowRun.ID,
		TaskRunID:     taskRun.ID,
		EventType:     task.TaskStarted,
		Attempt:       1,
		WorkerID:      "worker-1",
	}
	event.ID = uuid.New()
	event.CreatedAt = now

	if err := fixture.service.HandleTaskEvent(ctx, event); err != nil {
		t.Fatal(err)
	}
	if err := fixture.service.HandleTaskEvent(ctx, event); err != nil {
		t.Fatal(err)
	}
	if len(fixture.store.events) != 1 {
		t.Fatalf("stored events = %d, want 1", len(fixture.store.events))
	}
}

func TestHandleTaskEventRejectsInvalidTransition(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	definition := task.TaskDefinition{Name: "pending", TaskType: task.MOCK_TASK}
	definition.ID = uuid.New()
	fixture.store.definitions[definition.ID] = definition
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 0)
	fixture.store.workflowRuns[fixture.workflowRun.ID] = workflow.WorkflowRun{
		BaseModel:            fixture.workflowRun.BaseModel,
		WorkflowDefinitionID: fixture.workflowRun.WorkflowDefinitionID,
		Status:               workflow.RUNNING,
	}

	event := task.TaskEvent{
		WorkflowRunID: fixture.workflowRun.ID,
		TaskRunID:     taskRun.ID,
		EventType:     task.TaskCompleted,
	}
	event.ID = uuid.New()

	err := fixture.service.HandleTaskEvent(context.Background(), event)
	if !errors.Is(err, ErrInvalidStateTransition) {
		t.Fatalf("HandleTaskEvent() error = %v, want ErrInvalidStateTransition", err)
	}
}

func TestScheduleWorkflowRestoresPendingTaskWhenPublishFails(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	definition := task.TaskDefinition{Name: "publish-me", TaskType: task.MOCK_TASK}
	definition.ID = uuid.New()
	fixture.store.definitions[definition.ID] = definition
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 1)
	fixture.publisher.err = errors.New("kafka unavailable")

	err := fixture.service.ScheduleWorkflow(context.Background(), fixture.workflowRun.ID)
	if err == nil {
		t.Fatal("ScheduleWorkflow() error = nil, want publisher error")
	}
	restored := fixture.store.taskRuns[taskRun.ID]
	if restored.Status != task.TASK_RUN_PENDING || restored.RetryCount != 0 {
		t.Fatalf("task after publish failure = %+v, want original pending state", restored)
	}
	if len(fixture.store.attempts) != 0 {
		t.Fatalf("attempts after publish failure = %d, want 0", len(fixture.store.attempts))
	}

	fixture.publisher.err = nil
	if err := fixture.service.ScheduleWorkflow(context.Background(), fixture.workflowRun.ID); err != nil {
		t.Fatalf("ScheduleWorkflow() retry error = %v", err)
	}
	if got := fixture.store.taskRuns[taskRun.ID].Status; got != task.TASK_RUN_QUEUED {
		t.Fatalf("task status after retry = %v, want QUEUED", got)
	}
}

func TestFailedEventResumesRetryWithoutDoubleIncrement(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	definition := task.TaskDefinition{Name: "flaky-publisher", TaskType: task.MOCK_TASK}
	definition.ID = uuid.New()
	fixture.store.definitions[definition.ID] = definition
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 2)

	ctx := context.Background()
	if err := fixture.service.ScheduleWorkflow(ctx, fixture.workflowRun.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.startTask(ctx, taskRun.ID); err != nil {
		t.Fatal(err)
	}

	fixture.publisher.err = errors.New("kafka unavailable")
	event := fixture.failedEvent(taskRun.ID, "worker failed")
	if err := fixture.service.HandleTaskEvent(ctx, event); err == nil {
		t.Fatal("HandleTaskEvent() error = nil, want publisher error")
	}
	restored := fixture.store.taskRuns[taskRun.ID]
	if restored.Status != task.TASK_RUN_FAILED || restored.RetryCount != 0 {
		t.Fatalf("task after retry publish failure = %+v, want FAILED with retry count 0", restored)
	}

	fixture.publisher.err = nil
	if err := fixture.service.HandleTaskEvent(ctx, event); err != nil {
		t.Fatalf("HandleTaskEvent() redelivery error = %v", err)
	}
	retried := fixture.store.taskRuns[taskRun.ID]
	if retried.Status != task.TASK_RUN_QUEUED || retried.RetryCount != 1 {
		t.Fatalf("task after redelivery = %+v, want QUEUED with retry count 1", retried)
	}
	if len(fixture.store.events) != 2 {
		// The start event plus the recovered failure event.
		t.Fatalf("stored events = %d, want 2", len(fixture.store.events))
	}
}

func TestCompletedEventCanRecordAfterStateWasAlreadyCommitted(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	definition := task.TaskDefinition{Name: "terminal", TaskType: task.MOCK_TASK}
	definition.ID = uuid.New()
	fixture.store.definitions[definition.ID] = definition
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 0)

	ctx := context.Background()
	if err := fixture.service.ScheduleWorkflow(ctx, fixture.workflowRun.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.startTask(ctx, taskRun.ID); err != nil {
		t.Fatal(err)
	}

	result := "done"
	event := task.TaskEvent{
		WorkflowRunID: fixture.workflowRun.ID,
		TaskRunID:     taskRun.ID,
		EventType:     task.TaskCompleted,
		Attempt:       1,
		Result:        &result,
	}
	event.ID = uuid.New()
	event.CreatedAt = now.Add(time.Minute)
	fixture.store.createEventErr = errors.New("temporary database error")

	if err := fixture.service.HandleTaskEvent(ctx, event); err == nil {
		t.Fatal("HandleTaskEvent() error = nil, want event persistence error")
	}
	if got := fixture.store.workflowRuns[fixture.workflowRun.ID].Status; got != workflow.COMPLETED {
		t.Fatalf("workflow status = %v, want already committed COMPLETED", got)
	}
	if err := fixture.service.HandleTaskEvent(ctx, event); err != nil {
		t.Fatalf("HandleTaskEvent() redelivery error = %v", err)
	}
	if len(fixture.store.events) != 2 {
		// The start event plus the recovered completion event.
		t.Fatalf("stored events = %d, want 2", len(fixture.store.events))
	}
}

func TestCancelWorkflowCancelsActiveAndPendingTasks(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	rootDefinition := task.TaskDefinition{Name: "active", TaskType: task.MOCK_TASK}
	rootDefinition.ID = uuid.New()
	childDefinition := task.TaskDefinition{Name: "pending", TaskType: task.MOCK_TASK}
	childDefinition.ID = uuid.New()
	fixture.store.definitions[rootDefinition.ID] = rootDefinition
	fixture.store.definitions[childDefinition.ID] = childDefinition
	root := fixture.addTaskRun(rootDefinition, task.TASK_RUN_PENDING, 0, 0)
	child := fixture.addTaskRun(childDefinition, task.TASK_RUN_PENDING, 0, 0)
	fixture.store.dependencies = []task.TaskDependency{{
		TaskID:         childDefinition.ID,
		DependOnTaskID: rootDefinition.ID,
	}}

	ctx := context.Background()
	if err := fixture.service.ScheduleWorkflow(ctx, fixture.workflowRun.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.startTask(ctx, root.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.service.CancelWorkflow(ctx, fixture.workflowRun.ID, "user request"); err != nil {
		t.Fatalf("CancelWorkflow() error = %v", err)
	}

	if got := fixture.store.workflowRuns[fixture.workflowRun.ID].Status; got != workflow.CANCELLED {
		t.Fatalf("workflow status = %v, want CANCELLED", got)
	}
	if got := fixture.store.taskRuns[root.ID].Status; got != task.TASK_RUN_CANCELLED {
		t.Fatalf("active task status = %v, want CANCELLED", got)
	}
	if got := fixture.store.taskRuns[child.ID].Status; got != task.TASK_RUN_CANCELLED {
		t.Fatalf("pending task status = %v, want CANCELLED", got)
	}
	attempt := fixture.store.attempts[attemptKey{root.ID, 1}]
	if attempt.Status != task.TASK_ATTEMPT_CANCELLED {
		t.Fatalf("active attempt status = %v, want CANCELLED", attempt.Status)
	}
}

func TestStaleAttemptEventIsRejected(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	definition := task.TaskDefinition{Name: "retry", TaskType: task.MOCK_TASK}
	definition.ID = uuid.New()
	fixture.store.definitions[definition.ID] = definition
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_QUEUED, 1, 2)
	fixture.store.workflowRuns[fixture.workflowRun.ID] = workflow.WorkflowRun{
		BaseModel:            fixture.workflowRun.BaseModel,
		WorkflowDefinitionID: fixture.workflowRun.WorkflowDefinitionID,
		Status:               workflow.RUNNING,
	}
	attempt := task.TaskAttempt{
		TaskRunID:     taskRun.ID,
		AttemptNumber: 2,
		Status:        task.TASK_ATTEMPT_QUEUED,
	}
	attempt.ID = uuid.New()
	fixture.store.attempts[attemptKey{taskRun.ID, 2}] = attempt

	event := task.TaskEvent{
		WorkflowRunID: fixture.workflowRun.ID,
		TaskRunID:     taskRun.ID,
		EventType:     task.TaskStarted,
		Attempt:       1,
	}
	event.ID = uuid.New()

	err := fixture.service.HandleTaskEvent(context.Background(), event)
	if !errors.Is(err, ErrMismatchedTaskAttempt) {
		t.Fatalf("HandleTaskEvent() error = %v, want ErrMismatchedTaskAttempt", err)
	}
}

func TestScheduleWorkflowRejectsDependencyCycleBeforeStarting(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	firstDefinition := task.TaskDefinition{Name: "first", TaskType: task.MOCK_TASK}
	firstDefinition.ID = uuid.New()
	secondDefinition := task.TaskDefinition{Name: "second", TaskType: task.MOCK_TASK}
	secondDefinition.ID = uuid.New()
	fixture.store.definitions[firstDefinition.ID] = firstDefinition
	fixture.store.definitions[secondDefinition.ID] = secondDefinition
	fixture.addTaskRun(firstDefinition, task.TASK_RUN_PENDING, 0, 0)
	fixture.addTaskRun(secondDefinition, task.TASK_RUN_PENDING, 0, 0)
	fixture.store.dependencies = []task.TaskDependency{
		{TaskID: firstDefinition.ID, DependOnTaskID: secondDefinition.ID},
		{TaskID: secondDefinition.ID, DependOnTaskID: firstDefinition.ID},
	}

	err := fixture.service.ScheduleWorkflow(context.Background(), fixture.workflowRun.ID)
	if !errors.Is(err, ErrInvalidWorkflowGraph) {
		t.Fatalf("ScheduleWorkflow() error = %v, want ErrInvalidWorkflowGraph", err)
	}
	if got := fixture.store.workflowRuns[fixture.workflowRun.ID].Status; got != workflow.PENDING {
		t.Fatalf("workflow status = %v, want PENDING", got)
	}
	if len(fixture.publisher.commands) != 0 {
		t.Fatalf("published commands = %d, want 0", len(fixture.publisher.commands))
	}
}

func TestScheduleWorkflowRejectsInvalidDefinitionBeforeListingTaskRuns(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	firstDefinition := fixture.addTaskDefinition("first")
	secondDefinition := fixture.addTaskDefinition("second")
	fixture.store.dependencies = []task.TaskDependency{
		{TaskID: firstDefinition.ID, DependOnTaskID: secondDefinition.ID},
		{TaskID: secondDefinition.ID, DependOnTaskID: firstDefinition.ID},
	}

	err := fixture.service.ScheduleWorkflow(context.Background(), fixture.workflowRun.ID)
	if !errors.Is(err, ErrInvalidWorkflowGraph) {
		t.Fatalf("ScheduleWorkflow() error = %v, want ErrInvalidWorkflowGraph", err)
	}
	if fixture.store.listTaskRunsCalls != 0 {
		t.Fatalf("task runs listed = %d, want 0", fixture.store.listTaskRunsCalls)
	}
	if got := fixture.store.workflowRuns[fixture.workflowRun.ID].Status; got != workflow.PENDING {
		t.Fatalf("workflow status = %v, want PENDING", got)
	}
	if len(fixture.publisher.commands) != 0 {
		t.Fatalf("published commands = %d, want 0", len(fixture.publisher.commands))
	}
}

func TestCompletedEventRecoversWhenTaskUpdateFailsAfterAttemptUpdate(t *testing.T) {
	now := time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	definition := task.TaskDefinition{Name: "recoverable", TaskType: task.MOCK_TASK}
	definition.ID = uuid.New()
	fixture.store.definitions[definition.ID] = definition
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 0)

	ctx := context.Background()
	if err := fixture.service.ScheduleWorkflow(ctx, fixture.workflowRun.ID); err != nil {
		t.Fatal(err)
	}
	if err := fixture.startTask(ctx, taskRun.ID); err != nil {
		t.Fatal(err)
	}

	result := "done"
	event := task.TaskEvent{
		WorkflowRunID: fixture.workflowRun.ID,
		TaskRunID:     taskRun.ID,
		EventType:     task.TaskCompleted,
		Attempt:       1,
		Result:        &result,
	}
	event.ID = uuid.New()
	event.CreatedAt = now.Add(time.Minute)
	fixture.store.failTaskUpdateStatus = map[task.TaskRunStatus]int{task.TASK_RUN_COMPLETED: 1}

	if err := fixture.service.HandleTaskEvent(ctx, event); err == nil {
		t.Fatal("HandleTaskEvent() error = nil, want injected task update error")
	}
	if got := fixture.store.taskRuns[taskRun.ID].Status; got != task.TASK_RUN_RUNNING {
		t.Fatalf("task status after failed update = %v, want RUNNING", got)
	}
	if got := fixture.store.attempts[attemptKey{taskRun.ID, 1}].Status; got != task.TASK_ATTEMPT_COMPLETED {
		t.Fatalf("attempt status after failed task update = %v, want COMPLETED", got)
	}

	if err := fixture.service.HandleTaskEvent(ctx, event); err != nil {
		t.Fatalf("HandleTaskEvent() redelivery error = %v", err)
	}
	if got := fixture.store.workflowRuns[fixture.workflowRun.ID].Status; got != workflow.COMPLETED {
		t.Fatalf("workflow status = %v, want COMPLETED", got)
	}
}

type fixture struct {
	workflowDefinition workflow.WorkflowDefinition
	workflowRun        workflow.WorkflowRun
	store              *memoryStore
	publisher          *recordingPublisher
	service            *Service
	now                time.Time
}

func newFixture(now time.Time) *fixture {
	workflowDefinition := workflow.WorkflowDefinition{Status: workflow.ACTIVATE}
	workflowDefinition.ID = uuid.New()
	workflowRun := workflow.WorkflowRun{
		WorkflowDefinitionID: workflowDefinition.ID,
		Status:               workflow.PENDING,
	}
	workflowRun.ID = uuid.New()
	store := &memoryStore{
		workflowDefinitions: make(map[uuid.UUID]workflow.WorkflowDefinition),
		workflowRuns:        make(map[uuid.UUID]workflow.WorkflowRun),
		definitions:         make(map[uuid.UUID]task.TaskDefinition),
		taskRuns:            make(map[uuid.UUID]task.TaskRun),
		attempts:            make(map[attemptKey]task.TaskAttempt),
		events:              make(map[uuid.UUID]task.TaskEvent),
	}
	store.workflowDefinitions[workflowDefinition.ID] = workflowDefinition
	store.workflowRuns[workflowRun.ID] = workflowRun
	publisher := &recordingPublisher{}
	service := NewService(store, store, publisher, WithClock(func() time.Time {
		return now
	}))
	return &fixture{
		workflowDefinition: workflowDefinition,
		workflowRun:        workflowRun,
		store:              store,
		publisher:          publisher,
		service:            service,
		now:                now,
	}
}

func (fixture *fixture) addTaskDefinition(name string) task.TaskDefinition {
	definition := task.TaskDefinition{
		WorkflowDefinitionID: fixture.workflowRun.WorkflowDefinitionID,
		Name:                 name,
		TaskType:             task.MOCK_TASK,
	}
	definition.ID = uuid.New()
	fixture.store.definitions[definition.ID] = definition
	return definition
}

func (fixture *fixture) addTaskRun(
	definition task.TaskDefinition,
	status task.TaskRunStatus,
	retryCount uint,
	maxRetries uint,
) task.TaskRun {
	if definition.WorkflowDefinitionID == uuid.Nil {
		definition.WorkflowDefinitionID = fixture.workflowRun.WorkflowDefinitionID
	}
	fixture.store.definitions[definition.ID] = definition
	taskRun := task.TaskRun{
		WorkflowRunID:    fixture.workflowRun.ID,
		TaskDefinitionID: definition.ID,
		TaskDefinition:   definition,
		Status:           status,
		RetryCount:       retryCount,
		MaxRetries:       maxRetries,
	}
	taskRun.ID = uuid.New()
	fixture.store.taskRuns[taskRun.ID] = taskRun
	return taskRun
}

func (fixture *fixture) startTask(ctx context.Context, taskRunID uuid.UUID) error {
	taskRun := fixture.store.taskRuns[taskRunID]
	event := task.TaskEvent{
		WorkflowRunID: fixture.workflowRun.ID,
		TaskRunID:     taskRunID,
		EventType:     task.TaskStarted,
		Attempt:       taskRun.RetryCount + 1,
		WorkerID:      "worker-1",
	}
	event.ID = uuid.New()
	event.CreatedAt = fixture.now
	return fixture.service.HandleTaskEvent(ctx, event)
}

func (fixture *fixture) completeTask(ctx context.Context, taskRunID uuid.UUID, result string) error {
	taskRun := fixture.store.taskRuns[taskRunID]
	event := task.TaskEvent{
		WorkflowRunID: fixture.workflowRun.ID,
		TaskRunID:     taskRunID,
		EventType:     task.TaskCompleted,
		Attempt:       taskRun.RetryCount + 1,
		Result:        &result,
		WorkerID:      "worker-1",
	}
	event.ID = uuid.New()
	event.CreatedAt = fixture.now.Add(time.Minute)
	return fixture.service.HandleTaskEvent(ctx, event)
}

func (fixture *fixture) failTask(ctx context.Context, taskRunID uuid.UUID, message string) error {
	return fixture.service.HandleTaskEvent(ctx, fixture.failedEvent(taskRunID, message))
}

func (fixture *fixture) failedEvent(taskRunID uuid.UUID, message string) task.TaskEvent {
	taskRun := fixture.store.taskRuns[taskRunID]
	event := task.TaskEvent{
		WorkflowRunID: fixture.workflowRun.ID,
		TaskRunID:     taskRunID,
		EventType:     task.TaskFailed,
		Attempt:       taskRun.RetryCount + 1,
		Error:         &message,
		WorkerID:      "worker-1",
	}
	event.ID = uuid.New()
	event.CreatedAt = fixture.now.Add(time.Minute)
	return event
}

type attemptKey struct {
	taskRunID uuid.UUID
	number    uint
}

type memoryStore struct {
	workflowDefinitions  workflowDefinitionMap
	workflowRuns         workflowRunMap
	definitions          map[uuid.UUID]task.TaskDefinition
	taskRuns             taskRunMap
	dependencies         []task.TaskDependency
	attempts             map[attemptKey]task.TaskAttempt
	events               map[uuid.UUID]task.TaskEvent
	createEventErr       error
	failTaskUpdateStatus map[task.TaskRunStatus]int
	listTaskRunsCalls    int
}

type workflowDefinitionMap map[uuid.UUID]workflow.WorkflowDefinition
type workflowRunMap map[uuid.UUID]workflow.WorkflowRun
type taskRunMap map[uuid.UUID]task.TaskRun

func (store *memoryStore) GetWorkflowDefinitionById(
	_ context.Context,
	id uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	definition, ok := store.workflowDefinitions[id]
	if !ok {
		return workflow.WorkflowDefinition{}, errors.New("workflow definition not found")
	}
	return definition, nil
}

func (store *memoryStore) GetWorkflowRunById(
	_ context.Context,
	id uuid.UUID,
) (workflow.WorkflowRun, error) {
	run, ok := store.workflowRuns[id]
	if !ok {
		return workflow.WorkflowRun{}, errors.New("workflow run not found")
	}
	return run, nil
}

func (store *memoryStore) UpdateWorkflowRunByIdInternal(
	_ context.Context,
	id uuid.UUID,
	updates map[string]any,
) (workflow.WorkflowRun, error) {
	run := store.workflowRuns[id]
	if status, ok := updates["status"].(workflow.WorkflowRunStatus); ok {
		run.Status = status
	}
	if startedAt, ok := updates["started_at"].(time.Time); ok {
		run.StartedAt = &startedAt
	}
	if endedAt, ok := updates["ended_at"].(time.Time); ok {
		run.EndedAt = &endedAt
	}
	if cancelledAt, ok := updates["cancelled_at"].(time.Time); ok {
		run.CancelledAt = &cancelledAt
	}
	store.workflowRuns[id] = run
	return run, nil
}

func (store *memoryStore) GetTaskDefinitionById(
	_ context.Context,
	id uuid.UUID,
) (task.TaskDefinition, error) {
	definition, ok := store.definitions[id]
	if !ok {
		return task.TaskDefinition{}, errors.New("task definition not found")
	}
	return definition, nil
}

func (store *memoryStore) GetTaskDefinitionsByWorkflowDefinitionId(
	_ context.Context,
	workflowDefinitionID uuid.UUID,
) ([]task.TaskDefinition, error) {
	var definitions []task.TaskDefinition
	for _, definition := range store.definitions {
		if definition.WorkflowDefinitionID == workflowDefinitionID {
			definitions = append(definitions, definition)
		}
	}
	return definitions, nil
}

func (store *memoryStore) GetTaskRunById(
	_ context.Context,
	id uuid.UUID,
) (task.TaskRun, error) {
	taskRun, ok := store.taskRuns[id]
	if !ok {
		return task.TaskRun{}, errors.New("task run not found")
	}
	return taskRun, nil
}

func (store *memoryStore) GetTaskRunsByWorkflowRunId(
	_ context.Context,
	workflowRunID uuid.UUID,
) ([]task.TaskRun, error) {
	store.listTaskRunsCalls++
	var taskRuns []task.TaskRun
	for _, taskRun := range store.taskRuns {
		if taskRun.WorkflowRunID == workflowRunID {
			taskRuns = append(taskRuns, taskRun)
		}
	}
	return taskRuns, nil
}

func (store *memoryStore) UpdateTaskRunById(
	_ context.Context,
	id uuid.UUID,
	updates map[string]any,
) (task.TaskRun, error) {
	taskRun := store.taskRuns[id]
	if status, ok := updates["status"].(task.TaskRunStatus); ok {
		if store.failTaskUpdateStatus[status] > 0 {
			store.failTaskUpdateStatus[status]--
			return task.TaskRun{}, errors.New("injected task update error")
		}
		taskRun.Status = status
	}
	if retryCount, ok := updates["retry_count"].(uint); ok {
		taskRun.RetryCount = retryCount
	}
	if value, exists := updates["scheduled_at"]; exists {
		taskRun.ScheduledAt = timePointer(value)
	}
	if value, exists := updates["started_at"]; exists {
		taskRun.StartedAt = timePointer(value)
	}
	if value, exists := updates["ended_at"]; exists {
		taskRun.EndedAt = timePointer(value)
	}
	if value, exists := updates["output"]; exists {
		taskRun.Output, _ = value.(*string)
	}
	store.taskRuns[id] = taskRun
	return taskRun, nil
}

func (store *memoryStore) CreateTaskAttempt(
	_ context.Context,
	attempt task.TaskAttempt,
) (task.TaskAttempt, error) {
	key := attemptKey{attempt.TaskRunID, attempt.AttemptNumber}
	if _, exists := store.attempts[key]; exists {
		return task.TaskAttempt{}, errors.New("duplicate task attempt")
	}
	attempt.ID = uuid.New()
	store.attempts[key] = attempt
	return attempt, nil
}

func (store *memoryStore) DeleteTaskAttemptById(
	_ context.Context,
	id uuid.UUID,
) error {
	for key, attempt := range store.attempts {
		if attempt.ID == id {
			delete(store.attempts, key)
			return nil
		}
	}
	return errors.New("task attempt not found")
}

func (store *memoryStore) GetTaskAttemptByTaskRunAndNumber(
	_ context.Context,
	taskRunID uuid.UUID,
	attemptNumber uint,
) (task.TaskAttempt, error) {
	attempt, ok := store.attempts[attemptKey{taskRunID, attemptNumber}]
	if !ok {
		return task.TaskAttempt{}, errors.New("task attempt not found")
	}
	return attempt, nil
}

func (store *memoryStore) UpdateTaskAttemptById(
	_ context.Context,
	id uuid.UUID,
	updates map[string]any,
) (task.TaskAttempt, error) {
	for key, attempt := range store.attempts {
		if attempt.ID != id {
			continue
		}
		if status, ok := updates["status"].(task.TaskAttemptStatus); ok {
			attempt.Status = status
		}
		if workerID, ok := updates["worker_id"].(string); ok {
			attempt.WorkerID = workerID
		}
		if startedAt, ok := updates["started_at"].(time.Time); ok {
			attempt.StartedAt = &startedAt
		}
		if endedAt, ok := updates["ended_at"].(time.Time); ok {
			attempt.EndedAt = &endedAt
		}
		store.attempts[key] = attempt
		return attempt, nil
	}
	return task.TaskAttempt{}, errors.New("task attempt not found")
}

func (store *memoryStore) TaskAttemptRunningIdempotency(
	_ context.Context,
	id uuid.UUID,
) (bool, error) {
	for _, attempt := range store.attempts {
		if attempt.ID == id {
			return attempt.Status == task.TASK_ATTEMPT_RUNNING, nil
		}
	}
	return false, errors.New("task attempt not found")
}

func (store *memoryStore) GetTaskDependenciesByWorkflowDefinitionId(
	_ context.Context,
	_ uuid.UUID,
) ([]task.TaskDependency, error) {
	return store.dependencies, nil
}

func (store *memoryStore) GetTaskDependencyById(
	_ context.Context,
	id uuid.UUID,
) (task.TaskDependency, error) {
	for _, dependency := range store.dependencies {
		if dependency.ID == id {
			return dependency, nil
		}
	}
	return task.TaskDependency{}, errors.New("task dependency not found")
}

func (store *memoryStore) TaskEventExists(
	_ context.Context,
	eventID uuid.UUID,
) (bool, error) {
	_, exists := store.events[eventID]
	return exists, nil
}

func (store *memoryStore) CreateTaskEvent(
	_ context.Context,
	event task.TaskEvent,
) (task.TaskEvent, error) {
	if store.createEventErr != nil {
		err := store.createEventErr
		store.createEventErr = nil
		return task.TaskEvent{}, err
	}
	if _, exists := store.events[event.ID]; exists {
		return task.TaskEvent{}, errors.New("duplicate task event")
	}
	store.events[event.ID] = event
	return event, nil
}

type recordingPublisher struct {
	commands []task.Command
	err      error
}

func (publisher *recordingPublisher) PublishTask(
	_ context.Context,
	command task.Command,
) error {
	if publisher.err != nil {
		return publisher.err
	}
	publisher.commands = append(publisher.commands, command)
	return nil
}

func timePointer(value any) *time.Time {
	if value == nil {
		return nil
	}
	timestamp, ok := value.(time.Time)
	if ok {
		return &timestamp
	}
	pointer, _ := value.(*time.Time)
	return pointer
}
