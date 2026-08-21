package postgresql

import (
	"context"
	"fmt"

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
	return &PostgresWorkflowRepository{database: database}
}

func (repo *PostgresWorkflowRepository) CreateWorkflow(
	ctx context.Context,
	w workflow.Workflow,
) (workflow.Workflow, error) {
	result := repo.database.WithContext(ctx).Create(&w)
	return w, result.Error
}

func (repo *PostgresWorkflowRepository) GetWorkflowsByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]workflow.Workflow, error) {
	var w []workflow.Workflow
	result := repo.database.WithContext(ctx).Where("user_id = ?", userID).Find(&w)
	return w, result.Error
}

func (repo *PostgresWorkflowRepository) GetWorkflowByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowID uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := repo.database.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, workflowID).
		First(&w)
	return w, result.Error
}

func (repo *PostgresWorkflowRepository) UpdateWorkflowByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowID uuid.UUID,
	newWorkflow map[string]any,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := repo.database.WithContext(ctx).
		Model(&w).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND id = ?", userID, workflowID).
		Updates(newWorkflow)
	return w, repository.CheckRowsAffected(result)
}

func (repo *PostgresWorkflowRepository) SoftDeleteWorkflow(
	ctx context.Context,
	userID uuid.UUID,
	workflowID uuid.UUID,
) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := repo.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("user_id = ? AND id = ?", userID, workflowID).
		Delete(&w)
	return w, repository.CheckRowsAffected(result)
}

// WorkflowDefinition

func (repo *PostgresWorkflowRepository) CreateWorkflowDefinition(
	ctx context.Context,
	wD workflow.WorkflowDefinition,
) (workflow.WorkflowDefinition, error) {
	result := repo.database.WithContext(ctx).Create(&wD)
	return wD, result.Error
}

func (repo *PostgresWorkflowRepository) GetWorkflowDefinitionsByWorkflow(
	ctx context.Context,
	userID uuid.UUID,
	workflowID uuid.UUID,
) ([]workflow.WorkflowDefinition, error) {
	wD := []workflow.WorkflowDefinition{}
	result := repo.database.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflows.id = ?", userID, workflowID).
		Find(&wD)
	return wD, result.Error
}

func (repo *PostgresWorkflowRepository) GetWorkflowDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := repo.database.WithContext(ctx).
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", userID, workflowDefinitionID).
		First(&wD)
	return wD, result.Error
}

func (repo *PostgresWorkflowRepository) filterWorkflowByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) *gorm.DB {
	return repo.database.WithContext(ctx).
		Where(`EXISTS (
      SELECT 1 
      FROM workflows
      WHERE workflows.id = workflow_definitions.workflow_id AND workflows.user_id = ?
    )`, userID).
		Where("workflow_definitions.id = ?", workflowDefinitionID)
}

func (repo *PostgresWorkflowRepository) UpdateWorkflowDefinitionByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
	newWorkflowDefinition map[string]any,
) (workflow.WorkflowDefinition, error) {
	var wD workflow.WorkflowDefinition

	filter := repo.filterWorkflowByUserAndID(ctx, userID, workflowDefinitionID)
	v, ok := newWorkflowDefinition["status"]
	if ok {
		newStatus, ok := v.(workflow.WorkflowDefinitionStatus)
		if !ok {
			return wD, fmt.Errorf("")
		}

		filter = filter.Where(
			"workflow_definitions.status IN ?",
			workflow.ValidPreviousWorkflowDefinitionStatus(newStatus),
		)
	}

	result := filter.
		Clauses(clause.Returning{}).
		Model(&wD).
		Updates(newWorkflowDefinition)
	return wD, repository.CheckRowsAffected(result)
}

func (repo *PostgresWorkflowRepository) SoftDeleteWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	var wD workflow.WorkflowDefinition
	result := repo.filterWorkflowByUserAndID(ctx, userID, workflowDefinitionID).
		Clauses(clause.Returning{}).
		Delete(&wD)
	return wD, repository.CheckRowsAffected(result)
}

// WorkflowRun

func (repo *PostgresWorkflowRepository) GetWorkflowRunsByWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) ([]workflow.WorkflowRun, error) {
	var wR []workflow.WorkflowRun
	result := repo.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_definitions.id = ?", userID, workflowDefinitionID).
		Find(&wR)
	return wR, result.Error
}

func (repo *PostgresWorkflowRepository) GetWorkflowRunByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowRunID uuid.UUID,
) (workflow.WorkflowRun, error) {
	var wR workflow.WorkflowRun
	result := repo.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Joins("JOIN workflows ON workflows.id = workflow_definitions.workflow_id").
		Where("workflows.user_id = ? AND workflow_runs.id = ?", userID, workflowRunID).First(&wR)
	return wR, result.Error
}

func (repo *PostgresWorkflowRepository) filterWorkflowRunByUserAndID(
	ctx context.Context, userID uuid.UUID,
	workflowRunID uuid.UUID,
) *gorm.DB {
	return repo.database.WithContext(ctx).
		Where(`EXISTS (
      SELECT 1
      FROM workflow_definitions
      JOIN workflows ON workflows.id = workflow_definitions.workflow_id
      WHERE workflow_definitions.id = workflow_runs.workflow_definition_id AND workflows.user_id = ?
    )`, userID).
		Where("workflow_runs.id = ?", workflowRunID)
}

func (repo *PostgresWorkflowRepository) UpdateWorkflowRunByUserAndID(
	ctx context.Context,
	userID uuid.UUID,
	workflowRunID uuid.UUID,
	newWorkflowRun map[string]any,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := repo.filterWorkflowRunByUserAndID(ctx, userID, workflowRunID).
		Clauses(clause.Returning{}).
		Model(&wR).
		Updates(newWorkflowRun)
	return wR, repository.CheckRowsAffected(result)
}

func (repo *PostgresWorkflowRepository) SoftDeleteWorkflowRun(
	ctx context.Context,
	userID uuid.UUID,
	workflowRunID uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := repo.filterWorkflowRunByUserAndID(ctx, userID, workflowRunID).
		Clauses(clause.Returning{}).
		Delete(&wR)
	return wR, repository.CheckRowsAffected(result)
}

func (repo *PostgresWorkflowRepository) CountActiveWorkflowRunByUserAndWorkflowDefinition(
	ctx context.Context,
	userID uuid.UUID,
	workflowDefinitionID uuid.UUID,
) (int64, error) {
	var count int64
	result := repo.database.WithContext(ctx).
		Joins("JOIN workflow_definitions ON workflow_definitions.id = workflow_runs.workflow_definition_id").
		Where("user_id = ? AND workflow_definition_id = ? AND status = ?", userID, workflowDefinitionID, workflow.WORKFLOW_RUN_RUNNING).
		Model(&workflow.WorkflowRun{}).
		Count(&count)
	return count, result.Error
}

// Internal

type PostgresWorkflowRepositoryInternal struct {
	database *gorm.DB
}

func NewPostgresWorkflowRepositoryInternal(database *gorm.DB) workflow.WorkflowRepositoryInternal {
	return &PostgresWorkflowRepositoryInternal{database: database}
}

// Workflow

func (repo *PostgresWorkflowRepositoryInternal) GetWorkflowByID(
	ctx context.Context,
	workflowID uuid.UUID,
) (workflow.Workflow, error) {
	var w workflow.Workflow
	result := repo.database.WithContext(ctx).Where("id = ?", workflowID).First(&w)
	return w, result.Error
}

func (repo *PostgresWorkflowRepositoryInternal) HardDeleteWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
) (workflow.Workflow, error) {
	var w workflow.Workflow
	result := repo.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("id = ?", workflowID).
		Unscoped().
		Delete(&w)
	return w, repository.CheckRowsAffected(result)
}

// WorkflowDefinition

func (repo *PostgresWorkflowRepositoryInternal) GetWorkflowDefinitionByID(
	ctx context.Context,
	workflowDefinitionID uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	var wD workflow.WorkflowDefinition
	result := repo.database.WithContext(ctx).Where("id = ?", workflowDefinitionID).First(&wD)
	return wD, result.Error
}

func (repo *PostgresWorkflowRepositoryInternal) HardDeleteWorkflowDefinition(
	ctx context.Context,
	workflowDefinitionID uuid.UUID,
) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := repo.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("id = ?", workflowDefinitionID).
		Unscoped().
		Delete(&wD)
	return wD, repository.CheckRowsAffected(result)
}

// WorkflowRun

func (repo *PostgresWorkflowRepositoryInternal) CreateWorkflowRun(
	ctx context.Context,
	wR workflow.WorkflowRun,
) (workflow.WorkflowRun, error) {
	result := repo.database.WithContext(ctx).Create(&wR)
	return wR, result.Error
}

func (repo *PostgresWorkflowRepositoryInternal) GetWorkflowRunByID(
	ctx context.Context,
	workflowRunID uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := repo.database.WithContext(ctx).Where("id = ?", workflowRunID).First(&wR)
	return wR, result.Error
}

func (repo *PostgresWorkflowRepositoryInternal) UpdateWorkflowRunByID(
	ctx context.Context,
	workflowRunID uuid.UUID,
	newWorkflowRun map[string]any,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := repo.database.WithContext(ctx).
		Where("id = ?", workflowRunID).
		Clauses(clause.Returning{}).
		Model(&wR).
		Updates(newWorkflowRun)
	return wR, repository.CheckRowsAffected(result)
}

func (repo *PostgresWorkflowRepositoryInternal) HardDeleteWorkflowRun(
	ctx context.Context,
	workflowRunID uuid.UUID,
) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := repo.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("id = ?", workflowRunID).
		Unscoped().
		Delete(&wR)
	return wR, repository.CheckRowsAffected(result)
}
