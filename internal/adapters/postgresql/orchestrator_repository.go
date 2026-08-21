package postgresql

import (
	"github.com/Longseeyou/DuraFlow/internal/orchestrator"
	"github.com/Longseeyou/DuraFlow/internal/task"
	"github.com/Longseeyou/DuraFlow/internal/workflow"
	"gorm.io/gorm"
)

type PostgresOrchestratorRepository struct {
	workflow.WorkflowRepository

	workflow.WorkflowRepositoryInternal
	task.TaskRepositoryInternal
}

func NewPostgresOrchestratorRepository(
	database *gorm.DB,
) orchestrator.OrchestratorRepository {
	return &PostgresOrchestratorRepository{
		WorkflowRepository: NewPostgresWorkflowRepository(database),

		WorkflowRepositoryInternal: NewPostgresWorkflowRepositoryInternal(database),
		TaskRepositoryInternal:     NewPostgresTaskRepositoryInternal(database),
	}
}
