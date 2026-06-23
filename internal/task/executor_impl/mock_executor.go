package executorimpl

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"

	"github.com/Longseeyou/DuraFlow/internal/task"
)

type MockExecutor struct {
}

func (e MockExecutor) Execute(
	ctx context.Context,
	input string,
) (string, string, task.TaskRunStatus, error) {
	r := rand.Float64()
	fmt.Println("r:", r)
	if r < 0.7 {
		return "", strconv.FormatFloat(r, 'f', -1, 32), task.COMPLETED, nil
	} else if r < 0.8 {
		return "task error", "task error", task.FAILED, nil
	} else if r < 0.9 {
		return "executor error", "executor error", task.FAILED, fmt.Errorf("executor error")
	}
	os.Exit(1)
	return "", "", task.FAILED, nil
}
