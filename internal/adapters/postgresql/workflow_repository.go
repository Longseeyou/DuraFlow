package postgresql

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/shared/repository"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkflowRepositoryPostgres struct {
	database *gorm.DB
}

func NewWorkflowRepositoryPostgres(database *gorm.DB) workflow.WorkflowRepository {
	return WorkflowRepositoryPostgres{database: database}
}

func (wRP WorkflowRepositoryPostgres) CreateWorkflow(
	ctx context.Context,
	w workflow.Workflow,
) (workflow.Workflow, error) {
	result := wRP.database.WithContext(ctx).Create(&w)
	return w, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowByUser(
	ctx context.Context,
	uId uuid.UUID,
) ([]workflow.Workflow, error) {
	var w []workflow.Workflow
	result := wRP.database.WithContext(ctx).Where("user_id = ?", uId).Find(&w)
	return w, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRP.database.WithContext(ctx).Where("user_id = ? AND id = ?", uId, wId).First(&w)
	return w, result.Error
}

func (wRP WorkflowRepositoryPostgres) UpdateWorkflowById(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
	newW map[string]any,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRP.database.WithContext(ctx).
		Model(&w).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND id = ?", uId, wId).
		Updates(newW)
	return w, repository.CheckRowsAffected(result)
}

func (wRP WorkflowRepositoryPostgres) SoftDeleteWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND id = ?", uId, wId).
		Delete(&w)
	return w, repository.CheckRowsAffected(result)
}

func (wRP WorkflowRepositoryPostgres) HardDeleteWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND id = ?", uId, wId).
		Unscoped().
		Delete(&w)
	return w, repository.CheckRowsAffected(result)
}

// WorkflowDefinition

func (wRP WorkflowRepositoryPostgres) CreateWorkflowDefinition(
	ctx context.Context,
	wD workflow.WorkflowDefinition,
) (workflow.WorkflowDefinition, error) {
	result := wRP.database.WithContext(ctx).Create(&wD)
	return wD, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowDefinitionByWorkflow(
	ctx context.Context,
	uId uuid.UUID,
	wId uuid.UUID,
) ([]workflow.WorkflowDefinition, error) {
	wD := []workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.workflow_id = ?", uId, wId).
		Find(&wD)
	return wD, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowDefinitionByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		First(&wD)
	return wD, result.Error
}

func (wRP WorkflowRepositoryPostgres) UpdateWorkflowDefinitionById(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
	newWD map[string]any,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).
		Model(&wD).
		Clauses(clause.Returning{}).
		Table("workflows").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		Updates(newWD)
	return wD, repository.CheckRowsAffected(result)
}

func (wRP WorkflowRepositoryPostgres) SoftDeleteWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		Delete(&wD)
	return wD, repository.CheckRowsAffected(result)
}

func (wRP WorkflowRepositoryPostgres) HardDeleteWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		Unscoped().
		Delete(&wD)
	return wD, repository.CheckRowsAffected(result)
}

// WorkflowRun

func (wRP WorkflowRepositoryPostgres) CreateWorkflowRun(
	ctx context.Context,
	wR workflow.WorkflowRun,
) (workflow.WorkflowRun, error) {
	result := wRP.database.WithContext(ctx).Create(&wR)
	return wR, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowRunByWorkflowDefinition(
	ctx context.Context,
	uId uuid.UUID,
	wDId uuid.UUID,
) ([]workflow.WorkflowRun, error) {
	wR := []workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", uId, wDId).
		Find(&wR)
	return wR, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowRunByUserAndId(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Where("workflows.user_id = ? AND workflow_runs.id = ?", uId, wRId).
		First(&wR)
	return wR, result.Error
}

func (wRP WorkflowRepositoryPostgres) UpdateWorkflowRunById(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
	newWR map[string]any,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).
		Model(&wR).
		Clauses(clause.Returning{}).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Where("workflows.user_id = ? AND workflow_runs.id = ?", uId, wRId).
		Updates(newWR)
	return wR, repository.CheckRowsAffected(result)
}

func (wRP WorkflowRepositoryPostgres) SoftDeleteWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Where("workflows.user_id = ? AND workflow_runs.id = ?", uId, wRId).
		Delete(&wR)
	return wR, repository.CheckRowsAffected(result)
}

func (wRP WorkflowRepositoryPostgres) HardDeleteWorkflowRun(
	ctx context.Context,
	uId uuid.UUID,
	wRId uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Where("workflows.user_id = ? AND workflow_runs.id = ?", uId, wRId).
		Unscoped().
		Delete(&wR)
	return wR, repository.CheckRowsAffected(result)
}

// Internal

type WorkflowRepositoryInternalPostgres struct {
	database *gorm.DB
}

func NewWorkflowRepositoryInternalPostgres(database *gorm.DB) workflow.WorkflowRepositoryInternal {
	return WorkflowRepositoryInternalPostgres{database: database}
}

func (wRIP WorkflowRepositoryInternalPostgres) GetWorkflowById(
	ctx context.Context,
	wId uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRIP.database.WithContext(ctx).Where("id = ?", wId).First(&w)
	return w, result.Error
}

func (wRIP WorkflowRepositoryInternalPostgres) GetWorkflowDefinitionById(
	ctx context.Context,
	wDId uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRIP.database.WithContext(ctx).Where("id = ?", wDId).First(&wD)
	return wD, result.Error
}

func (wRIP WorkflowRepositoryInternalPostgres) GetWorkflowRunById(
	ctx context.Context,
	wRId uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRIP.database.WithContext(ctx).Where("id = ?", wRId).First(&wR)
	return wR, result.Error
}
