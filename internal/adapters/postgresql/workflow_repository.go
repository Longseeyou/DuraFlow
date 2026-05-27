package postgresql

import (
	"context"
	"fmt"

	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkflowRepositoryPostgres struct {
	database *gorm.DB
}

func NewWorkflowRepositoryPostgres(database *gorm.DB) workflow.WorkflowRepository {
	return WorkflowRepositoryPostgres{database: database}
}

func (workflowRepository WorkflowRepositoryPostgres) CreateWorkflow(ctx context.Context, w workflow.Workflow) workflow.Workflow {
	result := workflowRepository.database.WithContext(ctx).Create(&w)
	if result.Error != nil {
		fmt.Println("TODO CreateWorkflow(ctx context.Context, w workflow.Workflow) workflow.Workflow")
	}
	return w
}

func (workflowRepository WorkflowRepositoryPostgres) GetWorkflowByUser(ctx context.Context, uId uuid.UUID) []workflow.Workflow {
	var w []workflow.Workflow
	result := workflowRepository.database.WithContext(ctx).Where("user_id = ?", uId).Find(&w)
	if result.Error != nil {
		fmt.Println("TODO GetWorkflowByUser(ctx context.Context, uId uuid.UUID) []workflow.Workflow")
	}
	return w
}
func (workflowRepository WorkflowRepositoryPostgres) GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, w workflow.Workflow) workflow.Workflow {
	result := workflowRepository.database.WithContext(ctx).Where("user_id = ?", uId).First(&w)
	if result.Error != nil {
		fmt.Println("TODO GetWorkflowByUserAndId(ctx context.Context, uId uuid.UUID, w workflow.Workflow) workflow.Workflow")
	}
	return w
}

func (workflowRepository WorkflowRepositoryPostgres) UpdateWorkflowById(ctx context.Context, w workflow.Workflow) workflow.Workflow {
	result := workflowRepository.database.WithContext(ctx).Save(&w)
	if result.Error != nil {
		fmt.Println("TODO GetWorkflowByUser(ctx context.Context, uId uuid.UUID) []workflow.Workflow")
	}
	return w

}
