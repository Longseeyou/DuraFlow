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

func (wRP WorkflowRepositoryPostgres) CreateWorkflow(ctx context.Context, w workflow.Workflow) (workflow.Workflow, error) {
	result := wRP.database.WithContext(ctx).Create(&w)
	return w, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowByUser(ctx context.Context, uId uuid.UUID) ([]workflow.Workflow, error) {
	var w []workflow.Workflow
	result := wRP.database.WithContext(ctx).Where("user_id = ?", uId).Find(&w)
	return w, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRP.database.WithContext(ctx).Where("user_id = ? AND id = ?", uId, wId).First(&w)
	return w, result.Error
}

func (wRP WorkflowRepositoryPostgres) UpdateWorkflowById(ctx context.Context, uId uuid.UUID, wId uuid.UUID, newW map[string]any) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRP.database.WithContext(ctx).Model(&w).Clauses(clause.Returning{}).Where("user_id = ? AND id = ?", uId, wId).Updates(newW)
	return w, result.Error
}

func (wRP WorkflowRepositoryPostgres) SoftDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRP.database.WithContext(ctx).Where("user_id = ? AND id = ?", uId, wId).Delete(&w)
	return w, result.Error
}

func (wRP WorkflowRepositoryPostgres) HardDeleteWorkflow(ctx context.Context, uId uuid.UUID, wId uuid.UUID) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRP.database.WithContext(ctx).Where("user_id = ? AND id = ?", uId, wId).Unscoped().Delete(&w)
	return w, result.Error
}

// WorkflowDefinition

func (wRP WorkflowRepositoryPostgres) CreateWorkflowDefinition(ctx context.Context, wD workflow.WorkflowDefinition) (workflow.WorkflowDefinition, error) {
	result := wRP.database.WithContext(ctx).Create(&wD)
	return wD, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowDefinitionByWorkflow(ctx context.Context, wId uuid.UUID) ([]workflow.WorkflowDefinition, error) {
	wD := []workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).Where("workflow_id = ?", wId).Find(&wD)
	return wD, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowDefinitionByWorkflowAndId(ctx context.Context, wId uuid.UUID, wDId uuid.UUID) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).Where("workflow_id = ? AND id = ?", wId, wDId).First(&wD)
	return wD, result.Error
}

func (wRP WorkflowRepositoryPostgres) UpdateWorkflowDefinitionById(ctx context.Context, wId uuid.UUID, wDId uuid.UUID, newWD map[string]any) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).Model(&wD).Clauses(clause.Returning{}).Where("workflow_id = ? AND id = ?", wId, wDId).Updates(newWD)
	return wD, result.Error
}

func (wRP WorkflowRepositoryPostgres) SoftDeleteWorkflowDefinition(ctx context.Context, wId uuid.UUID, wDId uuid.UUID) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).Where("workflow_id = ? AND id = ?", wId, wDId).Delete(&wD)
	return wD, result.Error
}

func (wRP WorkflowRepositoryPostgres) HardDeleteWorkflowDefinition(ctx context.Context, wId uuid.UUID, wDId uuid.UUID) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRP.database.WithContext(ctx).Where("workflow_id = ? AND id = ?", wId, wDId).Unscoped().Delete(&wD)
	return wD, result.Error
}

// WorkflowRun

func (wRP WorkflowRepositoryPostgres) CreateWorkflowRun(ctx context.Context, wR workflow.WorkflowRun) (workflow.WorkflowRun, error) {
	result := wRP.database.WithContext(ctx).Create(&wR)
	return wR, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowRunByWorkflowDefinition(ctx context.Context, wDId uuid.UUID) ([]workflow.WorkflowRun, error) {
	wR := []workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).Where("workflow_definition_id = ?", wDId).Find(&wR)
	return wR, result.Error
}

func (wRP WorkflowRepositoryPostgres) GetWorkflowRunByWorkflowDefinitionAndId(ctx context.Context, wDId uuid.UUID, wRId uuid.UUID) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).Where("workflow_definition_id = ? AND id = ?", wDId, wRId).First(&wR)
	return wR, result.Error
}

func (wRP WorkflowRepositoryPostgres) UpdateWorkflowRunById(ctx context.Context, wDId uuid.UUID, wRId uuid.UUID, newWR map[string]any) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).Model(&wR).Clauses(clause.Returning{}).Where("workflow_definition_id = ? AND id = ?", wDId, wRId).Updates(newWR)
	return wR, result.Error
}

func (wRP WorkflowRepositoryPostgres) SoftDeleteWorkflowRun(ctx context.Context, wDId uuid.UUID, wRId uuid.UUID) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).Where("workflow_definition_id = ? AND id = ?", wDId, wRId).Delete(&wR)
	return wR, result.Error
}

func (wRP WorkflowRepositoryPostgres) HardDeleteWorkflowRun(ctx context.Context, wDId uuid.UUID, wRId uuid.UUID) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRP.database.WithContext(ctx).Where("workflow_definition_id = ? AND id = ?", wDId, wRId).Unscoped().Delete(&wR)
	return wR, result.Error
}

// Internal

type WorkflowRepositoryInternalPostgres struct {
	database *gorm.DB
}

func NewWorkflowRepositoryInternalPostgres(database *gorm.DB) workflow.WorkflowRepositoryInternal {
	return WorkflowRepositoryInternalPostgres{database: database}
}

func (wRIP WorkflowRepositoryInternalPostgres) GetWorkflowById(ctx context.Context, wId uuid.UUID) (workflow.Workflow, error) {
	w := workflow.Workflow{}
	result := wRIP.database.WithContext(ctx).Where("id = ?", wId).First(&w)
	return w, result.Error
}

func (wRIP WorkflowRepositoryInternalPostgres) GetWorkflowDefinitionById(ctx context.Context, wDId uuid.UUID) (workflow.WorkflowDefinition, error) {
	wD := workflow.WorkflowDefinition{}
	result := wRIP.database.WithContext(ctx).Where("id = ?", wDId).First(&wD)
	return wD, result.Error
}

func (wRIP WorkflowRepositoryInternalPostgres) GetWorkflowRunById(ctx context.Context, wRId uuid.UUID) (workflow.WorkflowRun, error) {
	wR := workflow.WorkflowRun{}
	result := wRIP.database.WithContext(ctx).Where("id = ?", wRId).First(&wR)
	return wR, result.Error
}
