package task

import (
	"github.com/Longseeyou/DuraFlow/internal/shared/model"
	"github.com/google/uuid"
)

type TaskEventType string

const (
	TaskStarted      TaskEventType = "task.started"
	TaskCompleted    TaskEventType = "task.completed"
	TaskFailed       TaskEventType = "task.failed"
	TaskCancelled    TaskEventType = "task.cancelled"
	TaskDeadLettered TaskEventType = "task.dead_lettered"
)

// TaskEvent is both the worker event contract and its persisted event record.
// Related domain objects are referenced by ID and loaded through repositories.
type TaskEvent struct {
	model.BaseModel
	WorkflowRunID uuid.UUID     `json:"workflow_run_id"`
	TaskRunID     uuid.UUID     `json:"task_run_id"`
	EventType     TaskEventType `json:"event_type"`
	Attempt       uint          `json:"attempt"`
	Result        *string       `json:"result,omitempty"`
	Error         *string       `json:"error,omitempty"`
	WorkerID      string        `json:"worker_id,omitempty"`
}
