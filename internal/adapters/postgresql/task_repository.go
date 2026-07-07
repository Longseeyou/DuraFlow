package postgresql

import (
	"context"
	"errors"

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

func (pTR *PostgresTaskRepository) CreateTaskDefinition(
	ctx context.Context,
	tD task.TaskDefinition,
) (task.TaskDefinition, error) {
	result := pTR.database.WithContext(ctx).Create(&tD)
	return tD, result.Error
}

func (pTR *PostgresTaskRepository) GetTaskDefinitionByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]task.TaskDefinition, error) {
	var tDs []task.TaskDefinition
	result := pTR.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		Find(&tDs)
	return tDs, result.Error
}

func (pTR *PostgresTaskRepository) GetTaskDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := pTR.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_definitions.id = ?", uId, tDId).First(&tD)
	return tD, result.Error
}

func (pTR *PostgresTaskRepository) filterTaskDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) *gorm.DB {
	return pTR.database.WithContext(ctx).Where(`EXISTS (
      SELECT workflow_definitions.id
      FROM workflow_definitions
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE workflow_definitions.id = task_definitions.workflow_definition_id AND workflows.user_id = ?
    )`, uId).
		Where("task_definitions.id = ?", tDId)
}

func (pTR *PostgresTaskRepository) UpdateTaskDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
	newTD map[string]any,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := pTR.filterTaskDefinitionByUserAndId(ctx, uId, tDId).
		Where("workflow_definitions.status = ?", workflow.WORKFLOW_DEFINITION_EDITING).
		Clauses(clause.Returning{}).
		Model(&tD).
		Updates(newTD)
	return tD, repository.CheckRowsAffected(result)
}

func (pTR *PostgresTaskRepository) SoftDeleteTaskDefinition(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := pTR.filterTaskDefinitionByUserAndId(ctx, uId, tDId).
		Where("workflow_definitions.status = ?", workflow.WORKFLOW_DEFINITION_EDITING).
		Clauses(clause.Returning{}).
		Delete(&tD)
	return tD, repository.CheckRowsAffected(result)
}

// TaskDependency

func (pTR *PostgresTaskRepository) CreateTaskDependency(
	ctx context.Context,
	tDp task.TaskDependency,
) (task.TaskDependency, error) {
	result := pTR.database.WithContext(ctx).Create(&tDp)
	return tDp, result.Error
}

func (pTR *PostgresTaskRepository) GetTaskDependencyByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]task.TaskDependency, error) {
	var tDps []task.TaskDependency
	result := pTR.database.WithContext(ctx).
		Joins("JOIN task_definitions td_task ON td_task.id = task_dependencies.task_id").
		Joins("JOIN task_definitions td_dep ON td_dep.id = task_dependencies.depend_on_task_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = td_task.workflow_definition_id AND workflow_definitions.id = td_dep.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		Find(&tDps)
	return tDps, result.Error
}

func (pTR *PostgresTaskRepository) GetTaskDependencyByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := pTR.database.WithContext(ctx).
		Joins("JOIN task_definitions td_task ON td_task.id = task_dependencies.task_id").
		Joins("JOIN task_definitions td_dep ON td_dep.id = task_dependencies.depend_on_task_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = td_task.workflow_definition_id AND workflow_definitions.id = td_dep.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_dependencies.id = ?", uId, tDpId).First(&tDp)
	return tDp, result.Error
}

func (pTR *PostgresTaskRepository) filterTaskDependencyByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) *gorm.DB {
	return pTR.database.WithContext(ctx).Where(`EXISTS (
      SELECT 1
      FROM task_definitions
      JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE task_definitions.id = task_dependencies.task_id AND workflows.user_id = ?
    )`, uId).
		Where(`EXISTS (
      SELECT 1
      FROM task_definitions td_dep
      JOIN task_definitions td_task ON td_task.id = task_dependencies.task_id
      WHERE td_dep.id = task_dependencies.depend_on_task_id AND td_dep.workflow_definition_id = td_task.workflow_definition_id 
    )`).
		Where("task_dependencies.id = ?", tDpId)
}

func (pTR *PostgresTaskRepository) UpdateTaskDependencyByUserAndId(
	ctx context.Context, uId uuid.UUID,
	tDpId uuid.UUID,
	newTDp map[string]any,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := pTR.filterTaskDependencyByUserAndId(ctx, uId, tDpId).
		Where("workflow_definitions.status = ?", workflow.WORKFLOW_DEFINITION_EDITING).
		Clauses(clause.Returning{}).
		Model(&tDp).
		Updates(newTDp)
	return tDp, repository.CheckRowsAffected(result)
}

func (pTR *PostgresTaskRepository) SoftDeleteTaskDependency(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := pTR.filterTaskDependencyByUserAndId(ctx, uId, tDpId).
		Where("workflow_definitions.status = ?", workflow.WORKFLOW_DEFINITION_EDITING).
		Clauses(clause.Returning{}).
		Delete(&tDp)
	return tDp, repository.CheckRowsAffected(result)
}

// TaskRun

func (pTR *PostgresTaskRepository) CreateTaskRun(
	ctx context.Context,
	tR task.TaskRun,
) (task.TaskRun, error) {
	result := pTR.database.WithContext(ctx).Create(&tR)
	return tR, result.Error
}

func (pTR *PostgresTaskRepository) GetTaskRunByWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) ([]task.TaskRun, error) {
	var tRs []task.TaskRun
	result := pTR.database.WithContext(ctx).
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.id = workflow_runs.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_runs.id = ?", uId, wRId).
		Find(&tRs)
	return tRs, result.Error
}

func (pTR *PostgresTaskRepository) GetTaskRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := pTR.database.WithContext(ctx).
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.id = workflow_runs.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tRId).First(&tR)
	return tR, result.Error
}

func (pTR *PostgresTaskRepository) filterTaskRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) *gorm.DB {
	return pTR.database.WithContext(ctx).Where(`EXISTS (
      SELECT 1
      FROM workflow_runs
      JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE workflow_runs.id = task_runs.workflow_run_id AND workflows.user_id = ?
    )`, uId).
		Where(`EXISTS (
      SELECT 1
      FROM workflow_runs
      JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id
      WHERE workflow_runs.id = task_runs.workflow_run_id AND workflow_runs.workflow_definition_id = task_definitions.workflow_definition_id 
    )`).
		Where("task_runs.id = ?", tRId)
}

func (pTR *PostgresTaskRepository) UpdateTaskRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
	newTR map[string]any,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := pTR.filterTaskRunByUserAndId(ctx, uId, tRId).
		Clauses(clause.Returning{}).
		Model(&tR).
		Updates(newTR)
	return tR, repository.CheckRowsAffected(result)
}

func (pTR *PostgresTaskRepository) SoftDeleteTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := pTR.filterTaskRunByUserAndId(ctx, uId, tRId).Clauses(clause.Returning{}).Delete(&tR)
	return tR, repository.CheckRowsAffected(result)
}

// TaskAttempt

func (pTR *PostgresTaskRepository) GetTaskAttemptByTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) ([]task.TaskAttempt, error) {
	var tAs []task.TaskAttempt
	result := pTR.database.WithContext(ctx).
		Joins("JOIN task_runs ON task_runs.id = task_attempts.task_run_id").
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id AND workflow_definitions.id = task_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tRId).
		Find(&tAs)
	return tAs, result.Error
}

func (pTR *PostgresTaskRepository) GetTaskAttemptByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := pTR.database.WithContext(ctx).
		Joins("JOIN task_runs ON task_runs.id = task_attempts.task_run_id").
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id AND workflow_definitions.id = task_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_attempts.id = ?", uId, tAId).
		First(&tA)
	return tA, result.Error
}

func (pTR *PostgresTaskRepository) filterTaskAttemptByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) *gorm.DB {
	return pTR.database.WithContext(ctx).Where(`EXISTS (
      SELECT 1
      FROM task_runs
      JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id
      JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE task_runs.id = task_attempts.task_run_id AND workflows.user_id = ?
    )`, uId).
		Where(`EXISTS (
      SELECT 1
      FROM task_runs
      JOIN workflow_runs ON workflow_runs.id = task_runs.workflow_run_id 
      JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id
      WHERE task_runs.id = task_attempts.task_run_id AND workflow_runs.workflow_definition_id = task_definitions.workflow_definition_id 
    )`).
		Where("task_attempts.id = ?", tAId)
}

func (pTR *PostgresTaskRepository) SoftDeleteTaskAttempt(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := pTR.filterTaskAttemptByUserAndId(ctx, uId, tAId).
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

func (pTRI *PostgresTaskRepositoryInternal) GetTaskDefinitionById(
	ctx context.Context,
	tDId uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := pTRI.database.WithContext(ctx).Where("id = ?", tDId).First(&tD)
	return tD, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) GetTaskDefinitionsByWorkflowDefinition(
	ctx context.Context,
	wDId uuid.UUID,
) ([]task.TaskDefinition, error) {
	var tDs []task.TaskDefinition
	result := pTRI.database.WithContext(ctx).
		Where("workflow_definition_id = ?", wDId).
		Find(&tDs)
	return tDs, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) HardDeleteTaskDefinition(
	ctx context.Context,
	tDId uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := pTRI.database.WithContext(ctx).
		Where("id = ?", tDId).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tD)
	return tD, repository.CheckRowsAffected(result)
}

func (pTRI *PostgresTaskRepositoryInternal) GetTaskDependencyById(
	ctx context.Context,
	tDpId uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := pTRI.database.WithContext(ctx).Where("id = ?", tDpId).First(&tDp)
	return tDp, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) GetTaskDependenciesByWorkflowDefinitionId(
	ctx context.Context,
	wDId uuid.UUID,
) ([]task.TaskDependency, error) {
	var tDps []task.TaskDependency
	result := pTRI.database.WithContext(ctx).
		Joins("JOIN task_definitions ON task_definitions.id = task_dependencies.task_id").
		Where("task_definitions.workflow_definition_id = ?", wDId).
		Find(&tDps)
	return tDps, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) HardDeleteTaskDependency(
	ctx context.Context,
	tDpId uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := pTRI.database.WithContext(ctx).
		Where("id = ?", tDpId).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tDp)
	return tDp, repository.CheckRowsAffected(result)
}

func (pTRI *PostgresTaskRepositoryInternal) GetTaskRunById(
	ctx context.Context,
	tRId uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := pTRI.database.WithContext(ctx).
		Preload("TaskDefinition").
		Where("id = ?", tRId).
		First(&tR)
	return tR, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) GetTaskRunsByWorkflowRunId(
	ctx context.Context,
	wRId uuid.UUID,
) ([]task.TaskRun, error) {
	var tRs []task.TaskRun
	result := pTRI.database.WithContext(ctx).
		Preload("TaskDefinition").
		Where("workflow_run_id = ?", wRId).
		Find(&tRs)
	return tRs, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) UpdateTaskRunById(
	ctx context.Context,
	tRId uuid.UUID,
	newTR map[string]any,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := pTRI.database.WithContext(ctx).
		Where("id = ?", tRId).
		Clauses(clause.Returning{}).
		Model(&tR).
		Updates(newTR)
	return tR, repository.CheckRowsAffected(result)
}

func (pTRI *PostgresTaskRepositoryInternal) HardDeleteTaskRun(
	ctx context.Context,
	tRId uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := pTRI.database.WithContext(ctx).
		Where("id = ?", tRId).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tR)
	return tR, repository.CheckRowsAffected(result)
}

func (pTRI *PostgresTaskRepositoryInternal) TaskRunIdempotency(
	ctx context.Context,
	tRId uuid.UUID,
	newStatus task.TaskRunStatus,
) (bool, error) {
	var tR task.TaskRun
	result := pTRI.database.WithContext(ctx).
		Where("id = ? AND status IN ?", tRId, task.ValidPreviousTaskRunStatus(newStatus)).
		Clauses(clause.Returning{}).
		Model(&tR).
		Updates(map[string]any{"status": newStatus})
	return result.RowsAffected == 1, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) CreateTaskAttempt(
	ctx context.Context,
	tA task.TaskAttempt,
) (task.TaskAttempt, error) {
	result := pTRI.database.WithContext(ctx).Create(&tA)
	return tA, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) HardDeleteTaskAttempt(
	ctx context.Context,
	tAId uuid.UUID,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := pTRI.database.WithContext(ctx).
		Where("id = ?", tAId).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tA)
	return tA, repository.CheckRowsAffected(result)
}

func (pTRI *PostgresTaskRepositoryInternal) GetTaskAttemptByTaskRunAndNumber(
	ctx context.Context,
	tRId uuid.UUID,
	attemptNumber uint,
) (task.TaskAttempt, error) {
	var attempt task.TaskAttempt
	result := pTRI.database.WithContext(ctx).
		Where("task_run_id = ? AND attempt_number = ?", tRId, attemptNumber).
		First(&attempt)
	return attempt, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) UpdateTaskAttemptById(
	ctx context.Context,
	tAId uuid.UUID,
	newTA map[string]any,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := pTRI.database.WithContext(ctx).
		Where("id = ?", tAId).
		Clauses(clause.Returning{}).
		Model(&tA).
		Updates(newTA)
	return tA, repository.CheckRowsAffected(result)
}

func (pTRI *PostgresTaskRepositoryInternal) TaskAttemptIdempotency(
	ctx context.Context,
	tAId uuid.UUID,
	newStatus task.TaskAttemptStatus,
) (bool, error) {
	var tA task.TaskAttempt
	result := pTRI.database.WithContext(ctx).
		Where("id = ? AND status IN ?", tAId, task.ValidPreviousTaskAttemptStatus(newStatus)).
		Clauses(clause.Returning{}).
		Model(&tA).
		Updates(map[string]any{"status": newStatus})
	return result.RowsAffected == 1, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) TaskEventExists(
	ctx context.Context,
	eventId uuid.UUID,
) (bool, error) {
	var event task.TaskEvent
	result := pTRI.database.WithContext(ctx).Select("id").Where("id = ?", eventId).First(&event)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return result.Error == nil, result.Error
}

func (pTRI *PostgresTaskRepositoryInternal) CreateTaskEvent(
	ctx context.Context,
	event task.TaskEvent,
) (task.TaskEvent, error) {
	result := pTRI.database.WithContext(ctx).Create(&event)
	return event, result.Error
}
