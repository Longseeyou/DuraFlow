package executorimpl

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
)

type MockExecutor struct {
}

func (e MockExecutor) Execute(ctx context.Context, input string) (string, string, error) {
	r := rand.Float64()
	if r < 0.6 {
		return "", strconv.FormatFloat(r, 'f', -1, 32), nil
	} else if r < 0.9 {
		return "error", "error", fmt.Errorf("error")
	}
	os.Exit(1)
	return "", "", nil
}
