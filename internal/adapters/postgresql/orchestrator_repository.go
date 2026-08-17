package postgresql

import (
	"gorm.io/gorm"
)

type PostgresOrchestratorRepository struct {
	PostgresTaskRepositoryInternal
	PostgresWorkflowRepositoryInternal
}

func NewPostgresOrchestratorRepository(database *gorm.DB) *PostgresOrchestratorRepository {
	return &PostgresOrchestratorRepository{
		PostgresTaskRepositoryInternal: PostgresTaskRepositoryInternal{
			database: database,
		},
		PostgresWorkflowRepositoryInternal: PostgresWorkflowRepositoryInternal{
			database: database,
		},
	}
}
