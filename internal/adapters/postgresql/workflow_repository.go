package postgresql

import (
	"context"

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

func (workflowRepository WorkflowRepositoryPostgres) CreateWorkflow(ctx context.Context, w workflow.Workflow) (workflow.Workflow, error) {
	result := workflowRepository.database.WithContext(ctx).Create(&w)
	return w, result.Error
}

func (workflowRepository WorkflowRepositoryPostgres) GetWorkflowByUser(ctx context.Context, uId uuid.UUID) ([]workflow.Workflow, error) {
	var w []workflow.Workflow
	result := workflowRepository.database.WithContext(ctx).Where("user_id = ?", uId).Find(&w)
	return w, result.Error
}
func (workflowRepository WorkflowRepositoryPostgres) GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := workflowRepository.database.WithContext(ctx).Where("user_id = ? AND id = ?", uId, wId).First(&w)
	return w, result.Error
}

func (workflowRepository WorkflowRepositoryPostgres) UpdateWorkflowById(ctx context.Context, uId uuid.UUID, wId uuid.UUID, newW map[string]any) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := workflowRepository.database.WithContext(ctx).Model(&w).Clauses(clause.Returning{}).Where("user_id = ? AND id = ?", uId, wId).Updates(newW)
	return w, result.Error
}

func (workflowRepository WorkflowRepositoryPostgres) SoftDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := workflowRepository.database.WithContext(ctx).Where("user_id = ? AND id = ?", uId, wId).Delete(&w)
	return w, result.Error
}

func (workflowRepository WorkflowRepositoryPostgres) HardDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := workflowRepository.database.WithContext(ctx).Where("user_id = ? AND id = ?", uId, wId).Unscoped().Delete(&w)
	return w, result.Error
}

// WorkflowDefinition

func (workflowRepository WorkflowRepositoryPostgres) CreateWorkflowDefinition(ctx context.Context, wD workflow.WorkflowDefinition) (workflow.WorkflowDefinition, error) {
	result := workflowRepository.database.WithContext(ctx).Create(&wD)
	return wD, result.Error
}

func (workflowRepository WorkflowRepositoryPostgres) GetWorkflowDefinitionByUserAndWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) ([]workflow.WorkflowDefinition, error) {
	wD := []workflow.WorkflowDefinition{}
	result := workflowRepository.database.WithContext(ctx).Joins("JOIN workflows ON workflows.user_id = ? AND workflows.id = ? AND workflows.id = workflow_definitions.workflow_id", uId, wId).Find(&wD)
	return wD, result.Error
}

func (workflowRepository WorkflowRepositoryPostgres) GetWorkflowDefinitionByUserAndWorkflowAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID, wDId uuid.UUID) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := workflowRepository.database.WithContext(ctx).Joins("JOIN workflows ON workflows.user_id = ? AND workflows.id = ? AND workflows.id = workflow_definitions.workflow_id", uId, wId).Where("id = ?", wDId).Find(&wD)
	return wD, result.Error
}
