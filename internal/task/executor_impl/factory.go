package executorimpl

import (
	"errors"

	"github.com/Longseeyou/DuraFlow/internal/task"
)

func NewTaskExecutor(taskType task.TaskType) (task.TaskExecutor, error) {
	switch taskType {
	case task.MOCK_TASK:
		return MockExecutor{}, nil
	}
	return nil, errors.New("Invalid taskType")
}
