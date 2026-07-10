package orchestrator

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/task"
)

func TestTaskPollerClaimsAndDispatchesTask(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	definition := fixture.addTaskDefinition("poll-me")
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 1)

	if err := fixture.service.taskPoller.QueueTask(
		context.Background(),
		taskRun,
		taskRun,
	); err != nil {
		t.Fatalf("QueueTask() error = %v", err)
	}
	if got := fixture.store.taskRuns[taskRun.ID].Status; got != task.TASK_RUN_QUEUED {
		t.Fatalf("task status = %v, want QUEUED", got)
	}
	attempt := fixture.store.attempts[attemptKey{taskRun.ID, 1}]
	if attempt.Status != task.TASK_ATTEMPT_QUEUED {
		t.Fatalf("attempt status = %v, want QUEUED", attempt.Status)
	}
	if len(fixture.publisher.commands) != 1 {
		t.Fatalf("published commands = %d, want 1", len(fixture.publisher.commands))
	}
	command := fixture.publisher.commands[0]
	if command.TaskRunID != taskRun.ID || command.TaskKey != definition.Name || command.Attempt != 1 {
		t.Fatalf("unexpected command: %+v", command)
	}
}

func TestTaskPollerRestoresClaimWhenPublishFails(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	fixture := newFixture(now)
	definition := fixture.addTaskDefinition("poll-me")
	taskRun := fixture.addTaskRun(definition, task.TASK_RUN_PENDING, 0, 1)
	fixture.publisher.err = errors.New("broker unavailable")

	err := fixture.service.taskPoller.QueueTask(context.Background(), taskRun, taskRun)
	if err == nil {
		t.Fatal("QueueTask() error = nil, want publish error")
	}
	if got := fixture.store.taskRuns[taskRun.ID].Status; got != task.TASK_RUN_PENDING {
		t.Fatalf("task status = %v, want PENDING", got)
	}
	if len(fixture.store.attempts) != 0 {
		t.Fatalf("unpublished attempts = %d, want 0", len(fixture.store.attempts))
	}
}
