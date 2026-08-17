package postgresql

import (
	"context"

	"github.com/Longseeyou/DuraFlow/internal/shared/repository"
	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostgresUserRepository struct {
	database *gorm.DB
}

func NewPostgresUserRepository(database *gorm.DB) user.UserRepository {
	return &PostgresUserRepository{database: database}
}

func (repo *PostgresUserRepository) CreateUser(ctx context.Context, u user.User) (user.User, error) {
	result := repo.database.WithContext(ctx).Create(&u)
	return u, result.Error
}

func (repo *PostgresUserRepository) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (user.User, error) {
	var u user.User
	result := repo.database.WithContext(ctx).Where("id = ?", userID).First(&u)
	return u, result.Error
}

func (repo *PostgresUserRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (user.User, error) {
	var u user.User
	result := repo.database.WithContext(ctx).Where("email = ?", email).First(&u)
	return u, result.Error
}

func (repo *PostgresUserRepository) UpdateUserByID(
	ctx context.Context,
	userID uuid.UUID,
	newUser map[string]any,
) (user.User, error) {
	var u user.User
	result := repo.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("id = ?", userID).
		Model(&u).
		Updates(newUser)
	return u, repository.CheckRowsAffected(result)
}

func (repo *PostgresUserRepository) SoftDeleteUser(
	ctx context.Context,
	userID uuid.UUID,
) (user.User, error) {
	var u user.User
	result := repo.database.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where("id = ?", userID).
		Delete(&u)
	return u, repository.CheckRowsAffected(result)
}

func (repo *PostgresUserRepository) HardDeleteUser(
	ctx context.Context,
	userID uuid.UUID,
) (user.User, error) {
	var u user.User
	result := repo.database.WithContext(ctx).
		Unscoped().
		Clauses(clause.Returning{}).
		Where("id = ?", userID).
		Delete(&u)
	return u, repository.CheckRowsAffected(result)
}

func (repo *PostgresUserRepository) RestoreUser(
	ctx context.Context,
	userID uuid.UUID,
) (user.User, error) {
	var u user.User
	result := repo.database.WithContext(ctx).
		Unscoped().
		Clauses(clause.Returning{}).
		Model(&u).
		Where("id = ?", userID).
		Update("deleted_at", nil)
	return u, repository.CheckRowsAffected(result)
}
