package executorimpl

import (
	"fmt"

	"github.com/Longseeyou/DuraFlow/internal/task"
)

func NewTaskExecutor(taskType task.TaskType) (task.TaskExecutor, error) {
	switch taskType {
	case task.MOCK_TASK:
		return MockExecutor{}, nil
	default:
		return nil, fmt.Errorf("NewTaskExecutor: Invalid taskType")
	}
}
