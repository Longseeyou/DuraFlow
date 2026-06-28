package task

import (
	"time"

	"github.com/google/uuid"
)

// Command is the task-dispatch message sent to workers.
type Command struct {
	WorkflowRunID uuid.UUID `json:"workflow_run_id"`
	TaskRunID     uuid.UUID `json:"task_run_id"`
	TaskKey       string    `json:"task_key"`
	TaskType      TaskType  `json:"task_type"`
	Payload       *string   `json:"payload,omitempty"`
	Attempt       uint      `json:"attempt"`
	CreatedAt     time.Time `json:"created_at"`
}
