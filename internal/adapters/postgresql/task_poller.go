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
	ctx context.Context,
) ([]task.TaskRun, error) {
	var tRs []task.TaskRun
	result := postgresTaskPoller.database.WithContext(ctx).
		Raw(`
        SELECT tr.*
        FROM task_runs tr
        JOIN task_definitions td ON td.id = tr.task_definition_id
        LEFT JOIN task_dependencies d ON d.task_id = td.id
        LEFT JOIN task_definitions tdp ON tdp.id = d.depend_on_task_id
        LEFT JOIN task_runs dep_tr ON dep_tr.task_definition_id = tdp.id
            AND dep_tr.workflow_run_id = tr.workflow_run_id
        WHERE tr.status = ANY(ARRAY[?, ?])
        GROUP BY tr.id
        HAVING COUNT(d.depend_on_task_id) = 0
            OR (
                COUNT(dep_tr.id) = COUNT(d.depend_on_task_id)
                AND COUNT(DISTINCT dep_tr.status) = 1
                AND MIN(dep_tr.status) = ?
            )
    `,
			string(task.TASK_RUN_PENDING),
			string(task.TASK_RUN_FAILED),
			string(task.TASK_RUN_COMPLETED),
		).
		Scan(&tRs)
	return tRs, result.Error
}
