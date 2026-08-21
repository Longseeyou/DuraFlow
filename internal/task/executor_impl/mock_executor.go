package executorimpl

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strconv"
	"time"
)

type MockExecutor struct{}

func (e MockExecutor) Execute(
	ctx context.Context,
	input string,
) (string, string, error) {
	r := rand.Float64()
	slog.Info("MockExecutor Execute", "r", r)
	time.Sleep(time.Duration(r * float64(time.Second)))

	if r < 0.8 {
		return "", strconv.FormatFloat(r, 'f', -1, 32), nil
	} else if r < 0.9 {
		return "task error", "task error", errors.New("task error")
	} else {
		return "executor error", "executor error", fmt.Errorf(
			"executor error",
		)
	}
}
