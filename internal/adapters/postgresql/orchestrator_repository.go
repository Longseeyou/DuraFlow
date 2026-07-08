package postgresql

import (
	"github.com/Longseeyou/DuraFlow/internal/orchestrator"
	"gorm.io/gorm"
)

type OrchestratorRepositories struct {
	Workflows orchestrator.WorkflowRepository
	Tasks     orchestrator.TaskRepository
}

func NewPostgresOrchestratorRepositories(database *gorm.DB) OrchestratorRepositories {
	return OrchestratorRepositories{
		Workflows: NewPostgresOrchestratorWorkflowRepository(database),
		Tasks:     NewPostgresOrchestratorTaskRepository(database),
	}
}

func NewPostgresOrchestratorWorkflowRepository(
	database *gorm.DB,
) orchestrator.WorkflowRepository {
	return NewPostgresWorkflowRepositoryInternal(database)
}

func NewPostgresOrchestratorTaskRepository(database *gorm.DB) orchestrator.TaskRepository {
	return NewPostgresTaskRepositoryInternal(database)
}
