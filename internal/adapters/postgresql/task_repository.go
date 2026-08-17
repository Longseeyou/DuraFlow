package postgresql

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/shared/repository"
	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresTaskRepository struct {
	database *gorm.DB
}

func NewPostgresTaskRepository(database *gorm.DB) task.TaskRepository {
	return &PostgresTaskRepository{database: database}
}

func (repo *PostgresTaskRepository) CreateTaskDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
	tD task.TaskDefinition,
) (task.TaskDefinition, error) {
	var created task.TaskDefinition

	result := repo.database.WithContext(ctx).Raw(`
		INSERT INTO task_definitions (
			id,
			workflow_definition_id,
			name,
			description,
			task_type,
			timeout
		)
		SELECT gen_random_uuid(), ?, ?, ?, ?, ?
		FROM workflow_definitions
		JOIN workflows ON workflows.id = workflow_definitions.workflow_id
		WHERE workflow_definitions.id = ?
		  AND workflows.user_id = ?
		  AND workflow_definitions.status = ?
		RETURNING id,
			workflow_definition_id,
			name,
			description,
			task_type,
			timeout
	`,
		workflowDefinitionID,
		tD.Name,
		tD.Description,
		tD.TaskType,
		tD.Timeout,
		workflowDefinitionID,
		userID,
		workflow.WORKFLOW_DEFINITION_EDITING,
	).Scan(&created)

	if result.Error != nil {
		return task.TaskDefinition{}, result.Error
	}

	if result.RowsAffected == 0 {
		return task.TaskDefinition{}, task.ErrWorkflowDefinitionNotEditing
	}

	return created, nil
}

func (repo *PostgresTaskRepository) GetTaskDefinitionsByWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) ([]task.TaskDefinition, error) {
	var tDs []task.TaskDefinition
	result := repo.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", userID, workflowDefinitionID).
		Find(&tDs)
	return tDs, result.Error
}

func (repo *PostgresTaskRepository) GetTaskDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskDefinitionID uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := repo.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_definitions.id = ?", userID, taskDefinitionID).
		First(&tD)
	return tD, result.Error
}

func (repo *PostgresTaskRepository) filterTaskDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskDefinitionID uuid.UUID,
) *gorm.DB {
	return repo.database.WithContext(ctx).Where(`EXISTS (
      SELECT workflow_definitions.id
      FROM workflow_definitions
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE workflow_definitions.id = task_definitions.workflow_definition_id AND workflows.user_id = ?
    )`, userID).
		Where("task_definitions.id = ?", taskDefinitionID)
}

func (repo *PostgresTaskRepository) UpdateTaskDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskDefinitionID uuid.UUID,
	newTaskDefinition map[string]any,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := repo.filterTaskDefinitionByUserAndID(ctx, userID, taskDefinitionID).
		Where("workflow_definitions.status = ?", workflow.WORKFLOW_DEFINITION_EDITING).
		Clauses(clause.Returning{}).
		Model(&tD).
		Updates(newTaskDefinition)
	return tD, repository.CheckRowsAffected(result)
}

func (repo *PostgresTaskRepository) SoftDeleteTaskDefinition(
	ctx context.Context,
	userID uuid.UUID,
	taskDefinitionID uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := repo.filterTaskDefinitionByUserAndID(ctx, userID, taskDefinitionID).
		Where("workflow_definitions.status = ?", workflow.WORKFLOW_DEFINITION_EDITING).
		Clauses(clause.Returning{}).
		Delete(&tD)
	return tD, repository.CheckRowsAffected(result)
}

// TaskDependency

func (repo *PostgresTaskRepository) CreateTaskDependency(
	ctx context.Context,
	tDp task.TaskDependency,
) (task.TaskDependency, error) {
	result := repo.database.WithContext(ctx).Create(&tDp)
	return tDp, result.Error
}

func (repo *PostgresTaskRepository) GetTaskDependenciesByWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) ([]task.TaskDependency, error) {
	var tDps []task.TaskDependency
	result := repo.database.WithContext(ctx).
		Joins("JOIN task_definitions td_task ON td_task.id = task_dependencies.task_id").
		Joins("JOIN task_definitions td_dep ON td_dep.id = task_dependencies.depend_on_task_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = td_task.workflow_definition_id AND workflow_definitions.id = td_dep.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", userID, workflowDefinitionID).
		Find(&tDps)
	return tDps, result.Error
}

func (repo *PostgresTaskRepository) GetTaskDependencyByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskDependencyID uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := repo.database.WithContext(ctx).
		Joins("JOIN task_definitions td_task ON td_task.id = task_dependencies.task_id").
		Joins("JOIN task_definitions td_dep ON td_dep.id = task_dependencies.depend_on_task_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = td_task.workflow_definition_id AND workflow_definitions.id = td_dep.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_dependencies.id = ?", userID, taskDependencyID).
		First(&tDp)
	return tDp, result.Error
}

func (repo *PostgresTaskRepository) filterTaskDependencyByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskDependencyID uuid.UUID,
) *gorm.DB {
	return repo.database.WithContext(ctx).Where(`EXISTS (
      SELECT 1
      FROM task_definitions
      JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE task_definitions.id = task_dependencies.task_id AND workflows.user_id = ?
    )`, userID).
		Where(`EXISTS (
      SELECT 1
      FROM task_definitions td_dep
      JOIN task_definitions td_task ON td_task.id = task_dependencies.task_id
      WHERE td_dep.id = task_dependencies.depend_on_task_id AND td_dep.workflow_definition_id = td_task.workflow_definition_id 
    )`).
		Where("task_dependencies.id = ?", taskDependencyID)
}

func (repo *PostgresTaskRepository) UpdateTaskDependencyByUserAndID(
	ctx context.Context, userID uuid.UUID,
	taskDependencyID uuid.UUID,
	newTaskDependency map[string]any,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := repo.filterTaskDependencyByUserAndID(ctx, userID, taskDependencyID).
		Where("workflow_definitions.status = ?", workflow.WORKFLOW_DEFINITION_EDITING).
		Clauses(clause.Returning{}).
		Model(&tDp).
		Updates(newTaskDependency)
	return tDp, repository.CheckRowsAffected(result)
}

func (repo *PostgresTaskRepository) SoftDeleteTaskDependency(
	ctx context.Context,
	userID uuid.UUID,
	taskDependencyID uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := repo.filterTaskDependencyByUserAndID(ctx, userID, taskDependencyID).
		Where("workflow_definitions.status = ?", workflow.WORKFLOW_DEFINITION_EDITING).
		Clauses(clause.Returning{}).
		Delete(&tDp)
	return tDp, repository.CheckRowsAffected(result)
}

// TaskRun

func (repo *PostgresTaskRepository) GetTaskRunsByWorkflowRun(
	ctx context.Context,
	userID uuid.UUID,
	workflowRunID uuid.UUID,
) ([]task.TaskRun, error) {
	var tRs []task.TaskRun
	result := repo.database.WithContext(ctx).
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.id = workflow_runs.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_runs.id = ?", userID, workflowRunID).
		Find(&tRs)
	return tRs, result.Error
}

func (repo *PostgresTaskRepository) GetTaskRunByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskRunID uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := repo.database.WithContext(ctx).
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.id = workflow_runs.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", userID, taskRunID).First(&tR)
	return tR, result.Error
}

func (repo *PostgresTaskRepository) filterTaskRunByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskRunID uuid.UUID,
) *gorm.DB {
	return repo.database.WithContext(ctx).Where(`EXISTS (
      SELECT 1
      FROM workflow_runs
      JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE workflow_runs.id = task_runs.workflow_run_id AND workflows.user_id = ?
    )`, userID).
		Where(`EXISTS (
      SELECT 1
      FROM workflow_runs
      JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id
      WHERE workflow_runs.id = task_runs.workflow_run_id AND workflow_runs.workflow_definition_id = task_definitions.workflow_definition_id 
    )`).
		Where("task_runs.id = ?", taskRunID)
}

func (repo *PostgresTaskRepository) UpdateTaskRunByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskRunID uuid.UUID,
	newTaskRun map[string]any,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := repo.filterTaskRunByUserAndID(ctx, userID, taskRunID).
		Clauses(clause.Returning{}).
		Model(&tR).
		Updates(newTaskRun)
	return tR, repository.CheckRowsAffected(result)
}

func (repo *PostgresTaskRepository) SoftDeleteTaskRun(
	ctx context.Context,
	userID uuid.UUID,
	taskRunID uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := repo.filterTaskRunByUserAndID(ctx, userID, taskRunID).
		Clauses(clause.Returning{}).
		Delete(&tR)
	return tR, repository.CheckRowsAffected(result)
}

// TaskAttempt

func (repo *PostgresTaskRepository) GetTaskAttemptsByTaskRun(
	ctx context.Context,
	userID uuid.UUID,
	taskRunID uuid.UUID,
) ([]task.TaskAttempt, error) {
	var tAs []task.TaskAttempt
	result := repo.database.WithContext(ctx).
		Joins("JOIN task_runs ON task_runs.id = task_attempts.task_run_id").
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id AND workflow_definitions.id = task_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", userID, taskRunID).
		Find(&tAs)
	return tAs, result.Error
}

func (repo *PostgresTaskRepository) GetTaskAttemptByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskAttemptID uuid.UUID,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := repo.database.WithContext(ctx).
		Joins("JOIN task_runs ON task_runs.id = task_attempts.task_run_id").
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id AND workflow_definitions.id = task_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_attempts.id = ?", userID, taskAttemptID).
		First(&tA)
	return tA, result.Error
}

func (repo *PostgresTaskRepository) filterTaskAttemptByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	taskAttemptID uuid.UUID,
) *gorm.DB {
	return repo.database.WithContext(ctx).Where(`EXISTS (
      SELECT 1
      FROM task_runs
      JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id
      JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE task_runs.id = task_attempts.task_run_id AND workflows.user_id = ?
    )`, userID).
		Where(`EXISTS (
      SELECT 1
      FROM task_runs
      JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id 
      JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id
      WHERE task_runs.id = task_attempts.task_run_id AND workflow_runs.workflow_definition_id = task_definitions.workflow_definition_id 
    )`).
		Where("task_attempts.id = ?", taskAttemptID)
}

func (repo *PostgresTaskRepository) SoftDeleteTaskAttempt(
	ctx context.Context,
	userID uuid.UUID,
	taskAttemptID uuid.UUID,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := repo.filterTaskAttemptByUserAndID(ctx, userID, taskAttemptID).
		Clauses(clause.Returning{}).
		Delete(&tA)
	return tA, repository.CheckRowsAffected(result)
}

// Internal

type PostgresTaskRepositoryInternal struct {
	database *gorm.DB
}

func NewPostgresTaskRepositoryInternal(database *gorm.DB) task.TaskRepositoryInternal {
	return &PostgresTaskRepositoryInternal{database: database}
}

func (repo *PostgresTaskRepositoryInternal) GetTaskDefinitionByID(
	ctx context.Context,
	taskDefinitionID uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := repo.database.WithContext(ctx).Where("id = ?", taskDefinitionID).First(&tD)
	return tD, result.Error
}

func (repo *PostgresTaskRepositoryInternal) GetTaskDefinitionsByWorkflowDefinition(
	ctx context.Context,
	workflowDefinitionID uuid.UUID,
) ([]task.TaskDefinition, error) {
	var tDs []task.TaskDefinition
	result := repo.database.WithContext(ctx).
		Where("workflow_definition_id = ?", workflowDefinitionID).
		Find(&tDs)
	return tDs, result.Error
}

func (repo *PostgresTaskRepositoryInternal) HardDeleteTaskDefinition(
	ctx context.Context,
	taskDefinitionID uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := repo.database.WithContext(ctx).
		Where("id = ?", taskDefinitionID).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tD)
	return tD, repository.CheckRowsAffected(result)
}

func (repo *PostgresTaskRepositoryInternal) GetTaskDependencyByID(
	ctx context.Context,
	taskDependencyID uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := repo.database.WithContext(ctx).Where("id = ?", taskDependencyID).First(&tDp)
	return tDp, result.Error
}

func (repo *PostgresTaskRepositoryInternal) GetTaskDependenciesByWorkflowDefinition(
	ctx context.Context,
	workflowDefinitionID uuid.UUID,
) ([]task.TaskDependency, error) {
	var tDps []task.TaskDependency
	result := repo.database.WithContext(ctx).
		Joins("JOIN task_definitions ON task_definitions.id = task_dependencies.task_id").
		Where("task_definitions.workflow_definition_id = ?", workflowDefinitionID).
		Find(&tDps)
	return tDps, result.Error
}

func (repo *PostgresTaskRepositoryInternal) HardDeleteTaskDependency(
	ctx context.Context,
	taskDependencyID uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := repo.database.WithContext(ctx).
		Where("id = ?", taskDependencyID).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tDp)
	return tDp, repository.CheckRowsAffected(result)
}

func (repo *PostgresTaskRepositoryInternal) CreateTaskRun(
	ctx context.Context,
	tR task.TaskRun,
) (task.TaskRun, error) {
	result := repo.database.WithContext(ctx).Create(&tR)
	return tR, result.Error
}

func (repo *PostgresTaskRepositoryInternal) GetTaskRunByID(
	ctx context.Context,
	taskRunID uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := repo.database.WithContext(ctx).
		Preload("TaskDefinition").
		Where("id = ?", taskRunID).
		First(&tR)
	return tR, result.Error
}

func (repo *PostgresTaskRepositoryInternal) GetTaskRunsByWorkflowRunID(
	ctx context.Context,
	workflowRunID uuid.UUID,
) ([]task.TaskRun, error) {
	var tRs []task.TaskRun
	result := repo.database.WithContext(ctx).
		Preload("TaskDefinition").
		Where("workflow_run_id = ?", workflowRunID).
		Find(&tRs)
	return tRs, result.Error
}

func (repo *PostgresTaskRepositoryInternal) UpdateTaskRunByID(
	ctx context.Context,
	taskRunID uuid.UUID,
	newTaskRun map[string]any,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := repo.database.WithContext(ctx).
		Where("id = ?", taskRunID).
		Clauses(clause.Returning{}).
		Model(&tR).
		Updates(newTaskRun)
	return tR, repository.CheckRowsAffected(result)
}

func (repo *PostgresTaskRepositoryInternal) HardDeleteTaskRun(
	ctx context.Context,
	taskRunID uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := repo.database.WithContext(ctx).
		Where("id = ?", taskRunID).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tR)
	return tR, repository.CheckRowsAffected(result)
}

func (repo *PostgresTaskRepositoryInternal) TaskRunIdempotency(
	ctx context.Context,
	taskRunID uuid.UUID,
	newStatus task.TaskRunStatus,
) (bool, error) {
	var tR task.TaskRun
	result := repo.database.WithContext(ctx).
		Where("id = ? AND status IN ?", taskRunID, task.ValidPreviousTaskRunStatus(newStatus)).
		Clauses(clause.Returning{}).
		Model(&tR).
		Updates(map[string]any{"status": newStatus})
	return result.RowsAffected == 1, result.Error
}

func (repo *PostgresTaskRepositoryInternal) MarkTaskRunRunning(
	ctx context.Context,
	taskRunID uuid.UUID,
	attemptNumber uint,
) (bool, error) {
	var tR task.TaskRun
	result := repo.database.WithContext(ctx).
		Where(
			"id = ? AND status = ? AND retry_count = ?",
			taskRunID,
			task.TASK_RUN_QUEUED,
			attemptNumber-1,
		).
		Clauses(clause.Returning{}).
		Model(&tR).
		Updates(map[string]any{"status": task.TASK_RUN_RUNNING})
	return result.RowsAffected == 1, result.Error
}

func (repo *PostgresTaskRepositoryInternal) GetPredecessorTaskRuns(
	ctx context.Context,
	taskRunID uuid.UUID,
) ([]task.TaskRun, error) {
	var tRs []task.TaskRun
	result := repo.database.WithContext(ctx).Raw(`
		SELECT dep_tr.*
    FROM task_runs tr
    JOIN task_dependencies d ON d.task_id = tr.task_definition_id
    JOIN task_runs dep_tr ON dep_tr.task_definition_id = d.depend_on_task_id AND dep_tr.workflow_run_id = tr.workflow_run_id
    WHERE tr.id = ?`, taskRunID).Scan(&tRs)
	return tRs, result.Error
}

func (repo *PostgresTaskRepositoryInternal) CreateTaskAttempt(
	ctx context.Context,
	taskAttempt task.TaskAttempt,
) (task.TaskAttempt, error) {
	result := repo.database.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "task_run_id"},
				{Name: "attempt_number"},
			},
			DoNothing: true,
		}).
		Create(&taskAttempt)
	if result.RowsAffected == 0 {
		return taskAttempt, nil
	}
	return taskAttempt, result.Error
}

func (repo *PostgresTaskRepositoryInternal) HardDeleteTaskAttempt(
	ctx context.Context,
	taskAttemptID uuid.UUID,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := repo.database.WithContext(ctx).
		Where("id = ?", taskAttemptID).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tA)
	return tA, repository.CheckRowsAffected(result)
}

func (repo *PostgresTaskRepositoryInternal) GetTaskAttemptByTaskRunAndNumber(
	ctx context.Context,
	taskRunID uuid.UUID,
	attemptNumber uint,
) (task.TaskAttempt, error) {
	var attempt task.TaskAttempt
	result := repo.database.WithContext(ctx).
		Where("task_run_id = ? AND attempt_number = ?", taskRunID, attemptNumber).
		First(&attempt)
	return attempt, result.Error
}

func (repo *PostgresTaskRepositoryInternal) UpdateTaskAttemptByID(
	ctx context.Context,
	taskAttemptID uuid.UUID,
	newTaskAttempt map[string]any,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := repo.database.WithContext(ctx).
		Where("id = ?", taskAttemptID).
		Clauses(clause.Returning{}).
		Model(&tA).
		Updates(newTaskAttempt)
	return tA, repository.CheckRowsAffected(result)
}

func (repo *PostgresTaskRepositoryInternal) GetTimedOutTaskRuns(
	ctx context.Context,
) ([]task.TaskRun, error) {
	var tRs []task.TaskRun
	result := repo.database.WithContext(ctx).
		Where("status = ? AND timeout_at < CURRENT_TIMESTAMP", task.TASK_RUN_RUNNING).
		Find(&tRs)
	return tRs, result.Error
}
