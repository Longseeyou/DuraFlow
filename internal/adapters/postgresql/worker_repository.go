package postgresql

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/worker"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresqlWorkerRepository struct {
	database *gorm.DB

	task.TaskRepositoryInternal
}

func NewPostgresqlWorkerRepository(database *gorm.DB) worker.WorkerRepository {
	return &PostgresqlWorkerRepository{
		database: database,

		TaskRepositoryInternal: NewPostgresTaskRepositoryInternal(database),
	}
}

func (repo *PostgresqlWorkerRepository) UpdateTaskRunByIDAndCreateTaskAttempt(
	ctx context.Context,
	taskRunID uuid.UUID,
	newTaskRun map[string]any,
	taskAttempt task.TaskAttempt,
) (task.TaskRun, task.TaskAttempt, error) {
	var tR task.TaskRun
	var tA task.TaskAttempt

	err := repo.database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &PostgresqlWorkerRepository{
			database: tx,
		}

		var err error

		tR, err = txRepo.UpdateTaskRunByID(
			ctx,
			taskRunID,
			newTaskRun,
		)
		if err != nil {
			return err
		}

		tA, err = txRepo.CreateTaskAttempt(ctx, taskAttempt)
		if err != nil {
			return err
		}

		return nil
	})

	return tR, tA, err
}
