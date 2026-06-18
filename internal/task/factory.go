package task

import (
	"errors"

	executorimpl "github.com/Longseeyou/DuraFlow/internal/task/executor_impl"
)

func NewTaskExecutor(taskType TaskType) (TaskExecutor, error) {
	switch taskType {
	case MOCK_TASK:
		return executorimpl.MockExecutor{}, nil
	}
	return nil, errors.New("Invalid taskType")
}
