package executorimpl

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/task"
)

type MockExecutor struct {
}

func (e MockExecutor) Execute(
	ctx context.Context,
	input string,
) (string, string, task.TaskAttemptStatus, error) {
	r := rand.Float64()
	slog.Info("MockExecutor Execute", "r", r)
	time.Sleep(time.Duration(r * 10 * time.Hour.Seconds()))

	if r < 0.7 {
		return "", strconv.FormatFloat(r, 'f', -1, 32), task.TASK_ATTEMPT_COMPLETED, nil
	} else if r < 0.8 {
		return "task error", "task error", task.TASK_ATTEMPT_FAILED, nil
	} else if r < 0.9 {
		return "executor error", "executor error", task.TASK_ATTEMPT_FAILED, fmt.Errorf(
			"MockExecutor Execute executor error",
		)
	} else {
		os.Exit(1)
		return "", "", task.TASK_ATTEMPT_FAILED, fmt.Errorf("")
	}
}
