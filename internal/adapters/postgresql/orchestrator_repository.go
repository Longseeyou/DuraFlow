package postgresql

import (
	"gorm.io/gorm"
)

type PostgresOrchestratorRepository struct {
	PostgresWorkflowRepository

	PostgresWorkflowRepositoryInternal
	PostgresTaskRepositoryInternal
}

func NewPostgresOrchestratorRepository(
	database *gorm.DB,
) *PostgresOrchestratorRepository {
	return &PostgresOrchestratorRepository{
		PostgresWorkflowRepository: PostgresWorkflowRepository{
			database: database,
		},
		PostgresTaskRepositoryInternal: PostgresTaskRepositoryInternal{
			database: database,
		},
		PostgresWorkflowRepositoryInternal: PostgresWorkflowRepositoryInternal{
			database: database,
		},
	}
}
