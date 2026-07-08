package orchestrator

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/shared/message"
	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/google/uuid"
)

func TestSchedulerProcessesWorkflowRunMessage(t *testing.T) {
	fixture := newFixture(time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC))
	definition := fixture.addTaskDefinition("scheduled")
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 0)

	payload, err := json.Marshal(struct {
		WorkflowRunID uuid.UUID `json:"workflow_run_id"`
	}{
		WorkflowRunID: fixture.workflowRun.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	scheduler := NewScheduler(fixture.service, nil)
	err = scheduler.ProcessMessage(context.Background(), message.Message{
		Topic: DefaultWorkflowRunTopic,
		Value: payload,
	})
	if err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if got := fixture.store.taskRuns[taskRun.ID].Status; got != task.TASK_RUN_QUEUED {
		t.Fatalf("task status = %v, want QUEUED", got)
	}
}

func TestSchedulerProcessesTaskEventMessage(t *testing.T) {
	fixture := newFixture(time.Date(2026, 6, 22, 12, 0, 0, 0, time.UTC))
	definition := fixture.addTaskDefinition("started")
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 0)

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
	event.CreatedAt = fixture.now
	payload, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}

	scheduler := NewScheduler(fixture.service, nil)
	err = scheduler.ProcessMessage(ctx, message.Message{
		Topic: DefaultTaskEventTopic,
		Value: payload,
	})
	if err != nil {
		t.Fatalf("ProcessMessage() error = %v", err)
	}
	if got := fixture.store.taskRuns[taskRun.ID].Status; got != task.TASK_RUN_RUNNING {
		t.Fatalf("task status = %v, want RUNNING", got)
	}
}
