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

func (tRP TaskRepositoryPostgres) CreateTaskDefinition(ctx context.Context, tD task.TaskDefinition) (task.TaskDefinition, error) {
	result := tRP.database.WithContext(ctx).Create(&tD)
	return tD, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskDefinitionByWorkflowDefinition(ctx context.Context, wDId uuid.UUID) ([]task.TaskDefinition, error) {
	var tD []task.TaskDefinition
	result := tRP.database.WithContext(ctx).Where("workflow_definition_id = ?", wDId).Find(&tD)
	return tD, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskDefinitionByWorkflowDefinitionAndId(ctx context.Context, wDId uuid.UUID, tDId uuid.UUID) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRP.database.WithContext(ctx).Where("workflow_definition_id = ? AND id = ?", wDId, tDId).First(&tD)
	return tD, result.Error
}

func (tRP TaskRepositoryPostgres) UpdateTaskDefinitionById(ctx context.Context, wDId uuid.UUID, tDId uuid.UUID, newTD map[string]any) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRP.database.WithContext(ctx).Model(&tD).Clauses(clause.Returning{}).Where("workflow_definition_id = ? AND id = ?", wDId, tDId).Updates(newTD)
	return tD, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) SoftDeleteTaskDefinition(ctx context.Context, wDId uuid.UUID, tDId uuid.UUID) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRP.database.WithContext(ctx).Where("workflow_definition_id = ? AND id = ?", wDId, tDId).Delete(&tD)
	return tD, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) HardDeleteTaskDefinition(ctx context.Context, wDId uuid.UUID, tDId uuid.UUID) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRP.database.WithContext(ctx).Where("workflow_definition_id = ? AND id = ?", wDId, tDId).Unscoped().Delete(&tD)
	return tD, repository.CheckRowsAffected(result)
}

// TaskDependency

func (tRP TaskRepositoryPostgres) CreateTaskDependency(ctx context.Context, tDp task.TaskDependency) (task.TaskDependency, error) {
	result := tRP.database.WithContext(ctx).Create(&tDp)
	return tDp, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskDependencyByWorkflowDefinition(ctx context.Context, wDId uuid.UUID) ([]task.TaskDependency, error) {
	var tDp []task.TaskDependency
	result := tRP.database.WithContext(ctx).Joins("JOIN task_definitions t ON t.id = task_dependencys.task_id").Joins("JOIN task_definitions tdp ON tdp.id = task_dependencys.depend_on_task_id").Where("t.workflow_definition_id = ? AND tdp.workflow_definition_id = ?", wDId, wDId).Find(&tDp)
	return tDp, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskDependencyByWorkflowDefinitionAndId(ctx context.Context, wDId uuid.UUID, tDpId uuid.UUID) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRP.database.WithContext(ctx).Joins("JOIN task_definitions t ON t.id = task_dependencys.task_id").Joins("JOIN task_definitions tdp ON tdp.id = task_dependencys.depend_on_task_id").Where("t.workflow_definition_id = ? AND tdp.workflow_definition_id = ? AND task_dependencys.id = ?", wDId, wDId, tDpId).First(&tDp)
	return tDp, result.Error
}

func (tRP TaskRepositoryPostgres) UpdateTaskDependencyById(ctx context.Context, tDpId uuid.UUID, newTDp map[string]any) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRP.database.WithContext(ctx).Model(&tDp).Clauses(clause.Returning{}).Where("id = ?", tDpId).Updates(newTDp)
	return tDp, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) SoftDeleteTaskDependency(ctx context.Context, tDpId uuid.UUID) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRP.database.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", tDpId).Delete(&tDp)
	return tDp, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) HardDeleteTaskDependency(ctx context.Context, tDpId uuid.UUID) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRP.database.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", tDpId).Unscoped().Delete(&tDp)
	return tDp, repository.CheckRowsAffected(result)
}

// TaskRun

func (tRP TaskRepositoryPostgres) CreateTaskRun(ctx context.Context, tR task.TaskRun) (task.TaskRun, error) {
	result := tRP.database.WithContext(ctx).Create(&tR)
	return tR, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskRunByWorkflowRun(ctx context.Context, wRId uuid.UUID) ([]task.TaskRun, error) {
	var tR []task.TaskRun
	result := tRP.database.WithContext(ctx).Where("workflow_run_id = ?", wRId).Find(&tR)
	return tR, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskRunByWorkflowRunAndId(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID) (task.TaskRun, error) {
	var tR task.TaskRun
	result := tRP.database.WithContext(ctx).Where("workflow_run_id = ? AND id = ?", wRId, tRId).First(&tR)
	return tR, result.Error
}

func (tRP TaskRepositoryPostgres) UpdateTaskRunById(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID, newTR map[string]any) (task.TaskRun, error) {
	var tR task.TaskRun
	result := tRP.database.WithContext(ctx).Model(&tR).Clauses(clause.Returning{}).Where("workflow_run_id = ? AND id = ?", wRId, tRId).Updates(newTR)
	return tR, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) SoftDeleteTaskRun(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID) (task.TaskRun, error) {
	var tR task.TaskRun
	result := tRP.database.WithContext(ctx).Clauses(clause.Returning{}).Where("workflow_run_id = ? AND id = ?", wRId, tRId).Delete(&tR)
	return tR, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) HardDeleteTaskRun(ctx context.Context, wRId uuid.UUID, tRId uuid.UUID) (task.TaskRun, error) {
	var tR task.TaskRun
	result := tRP.database.WithContext(ctx).Clauses(clause.Returning{}).Where("workflow_run_id = ? AND id = ?", wRId, tRId).Unscoped().Delete(&tR)
	return tR, repository.CheckRowsAffected(result)
}

// TaskAttempt

func (tRP TaskRepositoryPostgres) CreateTaskAttempt(ctx context.Context, tA task.TaskAttempt) (task.TaskAttempt, error) {
	result := tRP.database.WithContext(ctx).Create(&tA)
	return tA, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskAttemptByTaskRun(ctx context.Context, tRId uuid.UUID) ([]task.TaskAttempt, error) {
	var tA []task.TaskAttempt
	result := tRP.database.WithContext(ctx).Where("task_run_id = ?", tRId).Find(&tA)
	return tA, result.Error
}

func (tRP TaskRepositoryPostgres) GetTaskAttemptByTaskRunAndId(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := tRP.database.WithContext(ctx).Where("task_run_id = ? AND id = ?", tRId, tAId).First(&tA)
	return tA, result.Error
}

func (tRP TaskRepositoryPostgres) UpdateTaskAttemptById(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID, newTA map[string]any) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := tRP.database.WithContext(ctx).Model(&tA).Clauses(clause.Returning{}).Where("task_run_id = ? AND id = ?", tRId, tAId).Updates(newTA)
	return tA, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) SoftDeleteTaskAttempt(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := tRP.database.WithContext(ctx).Clauses(clause.Returning{}).Where("task_run_id = ? AND id = ?", tRId, tAId).Delete(&tA)
	return tA, repository.CheckRowsAffected(result)
}

func (tRP TaskRepositoryPostgres) HardDeleteTaskAttempt(ctx context.Context, tRId uuid.UUID, tAId uuid.UUID) (task.TaskAttempt, error) {
	var tA task.TaskAttempt
	result := tRP.database.WithContext(ctx).Clauses(clause.Returning{}).Where("task_run_id = ? AND id = ?", tRId, tAId).Unscoped().Delete(&tA)
	return tA, repository.CheckRowsAffected(result)
}

// Internal

type TaskRepositoryInternalPostgres struct {
	database *gorm.DB
}

func NewTaskRepositoryInternalPostgres(database *gorm.DB) task.TaskRepositoryInternal {
	return TaskRepositoryInternalPostgres{database: database}
}

func (tRIP TaskRepositoryInternalPostgres) GetTaskDefinitionById(ctx context.Context, tDId uuid.UUID) (task.TaskDefinition, error) {
	var tD task.TaskDefinition
	result := tRIP.database.WithContext(ctx).Where("id = ?", tDId).First(&tD)
	return tD, result.Error
}

func (tRIP TaskRepositoryInternalPostgres) GetTaskDependencyById(ctx context.Context, tDpId uuid.UUID) (task.TaskDependency, error) {
	var tDp task.TaskDependency
	result := tRIP.database.WithContext(ctx).Where("id = ?", tDpId).First(&tDp)
	return tDp, result.Error
}
