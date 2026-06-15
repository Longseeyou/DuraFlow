package postgresql

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/shared/repository"
	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TaskRepositoryPostgres struct {
	database *gorm.DB
}

func NewTaskRepositoryPostgres(database *gorm.DB) task.TaskRepository {
	return TaskRepositoryPostgres{database: database}
}

func (tRP TaskRepositoryPostgres) CreateTaskDefinition(
	ctx context.Context,
	tD task.TaskDefinition,
) (task.TaskDefinition, error) {
	result := tRP.database.WithContext(ctx).Create(&tD)
	return tD, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskDefinitionByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]task.TaskDefinition, error) {
	var tD []task.TaskDefinition
	result := tRP.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		Find(&tD)
	return tD, result.Error
}

func (tRP TaskRepositoryPostgres) filterTaskDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) *gorm.DB {
	return tRP.database.WithContext(ctx).Where(`EXISTS (
      SELECT workflow_definitions.id
      FROM workflow_definitions
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE workflow_definitions.id = task_definitions.workflow_definition_id AND workflows.user_id = ?
    )`, uId).
		Where("task_definitions.id = ?", tDId)
}

func (tRP TaskRepositoryPostgres) GetTaskDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRP.filterTaskDefinitionByUserAndId(ctx, uId, tDId).First(&tD)
	return tD, result.Error
}

func (tRP TaskRepositoryPostgres) UpdateTaskDefinitionById(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
	newTD map[string]any,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRP.filterTaskDefinitionByUserAndId(ctx, uId, tDId).
		Clauses(clause.Returning{}).
		Model(&tD).
		Updates(newTD)
	return tD, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) SoftDeleteTaskDefinition(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRP.filterTaskDefinitionByUserAndId(ctx, uId, tDId).
		Clauses(clause.Returning{}).
		Delete(&tD)
	return tD, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) HardDeleteTaskDefinition(
	ctx context.Context,
	uId uuid.UUID,
	tDId uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRP.filterTaskDefinitionByUserAndId(ctx, uId, tDId).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tD)
	return tD, repository.CheckRowsAffected(result)
}

// TaskDependency

func (tRP TaskRepositoryPostgres) CreateTaskDependency(
	ctx context.Context,
	tDp task.TaskDependency,
) (task.TaskDependency, error) {
	result := tRP.database.WithContext(ctx).Create(&tDp)
	return tDp, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskDependencyByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]task.TaskDependency, error) {
	var tDp []task.TaskDependency
	result := tRP.database.WithContext(ctx).
		Joins("JOIN task_definitions td_task ON td_task.id = task_dependencys.task_id").
		Joins("JOIN task_definitions td_dep ON td_dep.id = task_dependencys.depend_on_task_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = td_task.workflow_definition_id AND workflow_definitions.id = td_dep.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		Find(&tDp)
	return tDp, result.Error
}

func (tRP TaskRepositoryPostgres) filterTaskDependencyByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) *gorm.DB {
	return tRP.database.WithContext(ctx).Where(`EXISTS (
      SELECT 1
      FROM task_definitions
      JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE task_definitions.id = task_dependencys.task_id AND workflows.user_id = ?
    )`, uId).
		Where(`EXISTS (
      SELECT 1
      FROM task_definitions td_dep
      JOIN task_definitions td_task ON td_task.id = task_dependencies.task_id
      WHERE td_dep.id = task_dependencies.depend_on_task_id AND td_dep.workflow_definition_id = td_task.workflow_definition_id 
    )`).
		Where("task_dependencys.id = ?", tDpId)
}

func (tRP TaskRepositoryPostgres) GetTaskDependencyByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRP.filterTaskDependencyByUserAndId(ctx, uId, tDpId).First(&tDp)
	return tDp, result.Error
}

func (tRP TaskRepositoryPostgres) UpdateTaskDependencyById(
	ctx context.Context, uId uuid.UUID,
	tDpId uuid.UUID,
	newTDp map[string]any,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRP.filterTaskDependencyByUserAndId(ctx, uId, tDpId).
		Clauses(clause.Returning{}).
		Model(&tDp).
		Updates(newTDp)
	return tDp, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) SoftDeleteTaskDependency(
	ctx context.Context,
	uId uuid.UUID,
	tDpId uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRP.filterTaskDependencyByUserAndId(ctx, uId, tDpId).
		Clauses(clause.Returning{}).
		Delete(&tDp)
	return tDp, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) HardDeleteTaskDependency(
	ctx context.Context, uId uuid.UUID,
	tDpId uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRP.filterTaskDependencyByUserAndId(ctx, uId, tDpId).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&tDp)
	return tDp, repository.CheckRowsAffected(result)
}

// TaskRun

func (tRP TaskRepositoryPostgres) CreateTaskRun(
	ctx context.Context,
	tR task.TaskRun,
) (task.TaskRun, error) {
	result := tRP.database.WithContext(ctx).Create(&tR)
	return tR, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskRunByWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) ([]task.TaskRun, error) {
	var tR []task.TaskRun
	result := tRP.database.WithContext(ctx).
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.id = workflow_runs.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_runs.id = ?", uId, wRId).
		Find(&tR)
	return tR, result.Error
}

func (tRP TaskRepositoryPostgres) filterTaskRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) *gorm.DB {
	return tRP.database.WithContext(ctx).Where(`EXISTS (
      SELECT 1
      FROM workflow_definitions
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE workflow_definitions.id = task_definitions.workflow_definition_id AND workflows.user_id = ?
    )`, uId).
		Where("task_definitions.id = ?", tDId)
}

func (tRP TaskRepositoryPostgres) GetTaskRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := tRP.database.WithContext(ctx).
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tRId).
		First(&tR)
	return tR, result.Error
}

func (tRP TaskRepositoryPostgres) UpdateTaskRunById(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
	newTR map[string]any,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := tRP.database.WithContext(ctx).
		Model(&tR).
		Clauses(clause.Returning{}).
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tRId).
		Updates(newTR)
	return tR, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) SoftDeleteTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := tRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tRId).
		Delete(&tR)
	return tR, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) HardDeleteTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) (task.TaskRun, error) {
	var tR task.TaskRun
	result := tRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tRId).
		Unscoped().
		Delete(&tR)
	return tR, repository.CheckRowsAffected(result)
}

// TaskAttempt

func (tRP TaskRepositoryPostgres) CreateTaskAttempt(
	ctx context.Context,
	tA task.TaskAttempt,
) (task.TaskAttempt, error) {
	result := tRP.database.WithContext(ctx).Create(&tA)
	return tA, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskAttemptByTaskRun(
	ctx context.Context,
	uId uuid.UUID,
	tRId uuid.UUID,
) ([]task.TaskAttempt, error) {
	var tA []task.TaskAttempt
	result := tRP.database.WithContext(ctx).
		Joins("JOIN task_attempts ON task_attempts.task_run_id = task_runs.id").
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tRId).
		Find(&tA)
	return tA, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskAttemptByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := tRP.database.WithContext(ctx).
		Joins("JOIN task_attempts ON task_attempts.task_run_id = task_runs.id").
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tAId).
		First(&tA)
	return tA, result.Error
}

func (tRP TaskRepositoryPostgres) UpdateTaskAttemptById(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
	newTA map[string]any,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := tRP.database.WithContext(ctx).
		Model(&tA).
		Clauses(clause.Returning{}).
		Joins("JOIN task_attempts ON task_attempts.task_run_id = task_runs.id").
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tAId).
		Updates(newTA)
	return tA, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) SoftDeleteTaskAttempt(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := tRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Joins("JOIN task_attempts ON task_attempts.task_run_id = task_runs.id").
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tAId).
		Delete(&tA)
	return tA, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) HardDeleteTaskAttempt(
	ctx context.Context,
	uId uuid.UUID,
	tAId uuid.UUID,
) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := tRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Joins("JOIN task_attempts ON task_attempts.task_run_id = task_runs.id").
		Joins("JOIN task_definitions ON task_definitions.id = task_runs.task_definition_id").
		Joins("JOIN workflow_runs ON workflow_run.id = task_runs.workflow_run_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = task_definitions.workflow_definition_id AND workflow_definitions.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND task_runs.id = ?", uId, tAId).
		Unscoped().
		Delete(&tA)
	return tA, repository.CheckRowsAffected(result)
}

// Internal

type TaskRepositoryInternalPostgres struct {
	database *gorm.DB
}

func NewTaskRepositoryInternalPostgres(database *gorm.DB) task.TaskRepositoryInternal {
	return TaskRepositoryInternalPostgres{database: database}
}

func (tRIP TaskRepositoryInternalPostgres) GetTaskDefinitionById(
	ctx context.Context,
	tDId uuid.UUID,
) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRIP.database.WithContext(ctx).Where("id = ?", tDId).First(&tD)
	return tD, result.Error
}

func (tRIP TaskRepositoryInternalPostgres) GetTaskDependencyById(
	ctx context.Context,
	tDpId uuid.UUID,
) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRIP.database.WithContext(ctx).Where("id = ?", tDpId).First(&tDp)
	return tDp, result.Error
}
