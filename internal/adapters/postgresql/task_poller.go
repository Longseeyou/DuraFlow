package postgresql

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/orchestrator"
	"github.com/Longseeyou/DuraFlow/internal/task"
	"gorm.io/gorm"
)

type PostgresTaskPoller struct {
	database *gorm.DB
}

func NewPostgresTaskPoller(database *gorm.DB) orchestrator.TaskPoller {
	return PostgresTaskPoller{database: database}
}

func (postgresTaskPoller PostgresTaskPoller) PollTaskRun(
	ctx context.Context, numberOfTasks uint,
) ([]task.TaskRun, error) {
	var tRs []task.TaskRun
	result := postgresTaskPoller.database.WithContext(ctx).
		Where(
			"status = ANY(ARRAY[?, ?]) AND COALESCE(number_of_incomplete_predecessor_tasks, 0) = 0",
			string(task.TASK_RUN_PENDING),
			string(task.TASK_RUN_FAILED),
		).
		Limit(int(numberOfTasks)).
		Find(&tRs)
	return tRs, result.Error
}
