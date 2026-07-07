package postgresql

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/shared/repository"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresWorkflowRepository struct {
	database *gorm.DB
}

func NewPostgresWorkflowRepository(database *gorm.DB) workflow.WorkflowRepository {
	return PostgresWorkflowRepository{database: database}
}

func (pWR PostgresWorkflowRepository) CreateWorkflow(
	ctx context.Context,
	w workflow.Workflow,
) (workflow.Workflow, error) {
	result := pWR.database.WithContext(ctx).Create(&w)
	return w, result.Error
}

func (pWR PostgresWorkflowRepository) GetWorkflowByUser(
	ctx context.Context,
	uId uuid.UUID,
) ([]workflow.Workflow, error) {
	var w []workflow.Workflow
	result := pWR.database.WithContext(ctx).Where("user_id = ?", uId).Find(&w)
	return w, result.Error
}

func (pWR PostgresWorkflowRepository) GetWorkflowByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := pWR.database.WithContext(ctx).Where("user_id = ? AND id = ?", uId, wId).First(&w)
	return w, result.Error
}

func (pWR PostgresWorkflowRepository) UpdateWorkflowById(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
	newW map[string]any,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := pWR.database.WithContext(ctx).
		Model(&w).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND id = ?", uId, wId).
		Updates(newW)
	return w, repository.CheckRowsAffected(result)
}

func (pWR PostgresWorkflowRepository) SoftDeleteWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := pWR.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND id = ?", uId, wId).
		Delete(&w)
	return w, repository.CheckRowsAffected(result)
}

func (pWR PostgresWorkflowRepository) HardDeleteWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := pWR.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND id = ?", uId, wId).
		Unscoped().
		Delete(&w)
	return w, repository.CheckRowsAffected(result)
}

// WorkflowDefinition

func (pWR PostgresWorkflowRepository) CreateWorkflowDefinition(
	ctx context.Context,
	wD workflow.WorkflowDefinition,
) (workflow.WorkflowDefinition, error) {
	result := pWR.database.WithContext(ctx).Create(&wD)
	return wD, result.Error
}

func (pWR PostgresWorkflowRepository) GetWorkflowDefinitionByWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) ([]workflow.WorkflowDefinition, error) {
	wD := []workflow.WorkflowDefinition{}
	result := pWR.database.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflows.id = ?", uId, wId).
		Find(&wD)
	return wD, result.Error
}

func (pWR PostgresWorkflowRepository) GetWorkflowDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := pWR.database.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).First(&wD)
	return wD, result.Error
}

func (pWR PostgresWorkflowRepository) filterWorkflowByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) *gorm.DB {
	return pWR.database.WithContext(ctx).
		Where(`EXISTS (
      SELECT 1 
      FROM workflows
      WHERE workflows.id = workflow_definitions.workflow_id AND workflows.user_id = ?
    )`, uId).
		Where("workflow_definitions.id = ?", wDId)
}

func (pWR PostgresWorkflowRepository) UpdateWorkflowDefinitionById(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
	newWD map[string]any,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := pWR.filterWorkflowByUserAndId(ctx, uId, wDId).
		Clauses(clause.Returning{}).
		Model(&wD).
		Updates(newWD)
	return wD, repository.CheckRowsAffected(result)
}

func (pWR PostgresWorkflowRepository) SoftDeleteWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := pWR.filterWorkflowByUserAndId(ctx, uId, wDId).
		Clauses(clause.Returning{}).
		Delete(&wD)
	return wD, repository.CheckRowsAffected(result)
}

func (pWR PostgresWorkflowRepository) HardDeleteWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := pWR.filterWorkflowByUserAndId(ctx, uId, wDId).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&wD)
	return wD, repository.CheckRowsAffected(result)
}

// WorkflowRun

func (pWR PostgresWorkflowRepository) CreateWorkflowRun(
	ctx context.Context,
	wR workflow.WorkflowRun,
) (workflow.WorkflowRun, error) {
	result := pWR.database.WithContext(ctx).Create(&wR)
	return wR, result.Error
}

func (pWR PostgresWorkflowRepository) GetWorkflowRunByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]workflow.WorkflowRun, error) {
	wR := []workflow.WorkflowRun{}
	result := pWR.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		Find(&wR)
	return wR, result.Error
}

func (pWR PostgresWorkflowRepository) GetWorkflowRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := pWR.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_runs.id = ?", uId, wRId).First(&wR)
	return wR, result.Error
}

func (pWR PostgresWorkflowRepository) filterWorkflowRunByUserAndId(
	ctx context.Context, uId uuid.UUID,
	wRId uuid.UUID,
) *gorm.DB {
	return pWR.database.WithContext(ctx).
		Where(`EXISTS (
      SELECT 1
      FROM workflow_definitions
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE workflow_definitions.id = workflow_runs.workflow_definition_id AND workflows.user_id = ?
    )`, uId).
		Where("workflow_runs.id = ?", wRId)
}

func (pWR PostgresWorkflowRepository) UpdateWorkflowRunById(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
	newWR map[string]any,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := pWR.filterWorkflowRunByUserAndId(ctx, uId, wRId).
		Clauses(clause.Returning{}).
		Model(&wR).
		Updates(newWR)
	return wR, repository.CheckRowsAffected(result)
}

func (pWR PostgresWorkflowRepository) SoftDeleteWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := pWR.filterWorkflowRunByUserAndId(ctx, uId, wRId).
		Clauses(clause.Returning{}).
		Delete(&wR)
	return wR, repository.CheckRowsAffected(result)
}

func (pWR PostgresWorkflowRepository) HardDeleteWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := pWR.filterWorkflowRunByUserAndId(ctx, uId, wRId).
		Clauses(clause.Returning{}).
		Unscoped().
		Delete(&wR)
	return wR, repository.CheckRowsAffected(result)
}

// Internal

type PostgresWorkflowRepositoryInternal struct {
	database *gorm.DB
}

func NewPostgresWorkflowRepositoryInternal(database *gorm.DB) workflow.WorkflowRepositoryInternal {
	return PostgresWorkflowRepositoryInternal{database: database}
}

func (pWRI PostgresWorkflowRepositoryInternal) GetWorkflowById(
	ctx context.Context,
	wId uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := pWRI.database.WithContext(ctx).Where("id = ?", wId).First(&w)
	return w, result.Error
}

func (pWRI PostgresWorkflowRepositoryInternal) GetWorkflowDefinitionById(
	ctx context.Context,
	wDId uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := pWRI.database.WithContext(ctx).Where("id = ?", wDId).First(&wD)
	return wD, result.Error
}

func (pWRI PostgresWorkflowRepositoryInternal) GetWorkflowRunById(
	ctx context.Context,
	wRId uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := pWRI.database.WithContext(ctx).Where("id = ?", wRId).First(&wR)
	return wR, result.Error
}

func (pWRI PostgresWorkflowRepositoryInternal) UpdateWorkflowRunByIdInternal(
	ctx context.Context,
	wRId uuid.UUID,
	newWR map[string]any,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := pWRI.database.WithContext(ctx).
		Where("id = ?", wRId).
		Clauses(clause.Returning{}).
		Model(&wR).
		Updates(newWR)
	return wR, repository.CheckRowsAffected(result)
}

func (pWRI PostgresWorkflowRepositoryInternal) UpdateWorkflowRunByIdAndStatusInternal(
	ctx context.Context,
	wRId uuid.UUID,
	statuses []workflow.WorkflowRunStatus,
	newWR map[string]any,
) (workflow.WorkflowRun, bool, error) {
	wR := workflow.WorkflowRun{}
	result := pWRI.database.WithContext(ctx).
		Where("id = ? AND status IN ?", wRId, statuses).
		Clauses(clause.Returning{}).
		Model(&wR).
		Updates(newWR)
	if result.Error != nil {
		return workflow.WorkflowRun{}, false, result.Error
	}
	if result.RowsAffected == 0 {
		return workflow.WorkflowRun{}, false, nil
	}
	return wR, true, nil
}
