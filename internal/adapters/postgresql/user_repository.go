package postgresql

import (
	"context"
	"fmt"

	"github.com/Longseeyou/DuraFlow/internal/user"
	"gorm.io/gorm"
)

type UserRepositoryPostgres struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepositoryPostgres {
	return UserRepositoryPostgres{database: database}
}

func (userRepository UserRepositoryPostgres) CreateUser(ctx context.Context, user user.User) user.User {
	result := userRepository.database.WithContext(ctx).Create(&user)
	if result.Error != nil {
		fmt.Println("TODO CreateUser(ctx context.Context, user user.User) user.User")
	}
	return user
}

func (userRepository UserRepositoryPostgres) SoftDeleteUser(ctx context.Context, user user.User) user.User {
	result := userRepository.database.WithContext(ctx).Delete(&user)
	if result.Error != nil {
		fmt.Println("TODO SoftDeleteUser(ctx context.Context, user user.User) user.User")
		return user
	}

	return user
}

func (userRepository UserRepositoryPostgres) HardDeleteUser(ctx context.Context, user user.User) user.User {
	result := userRepository.database.WithContext(ctx).Unscoped().Delete(&user)
	if result.Error != nil {
		fmt.Println("TODO HardDeleteUser(ctx context.Context, user user.User) user.User")
	}
	return user
}

func (userRepository UserRepositoryPostgres) RestoreUser(ctx context.Context, user user.User) user.User {
	result := userRepository.database.WithContext(ctx).Unscoped().Model(&user).Where("id = ?", user.ID).Update("deleted_at", nil)
	if result.Error != nil {
		fmt.Println("TODO RestoreUser(ctx context.Context, user user.User) user.User")
	}
	return user
}

func (userRepository UserRepositoryPostgres) GetUser(ctx context.Context, user user.User) user.User {
	result := userRepository.database.WithContext(ctx).Where("id = ?", user.ID).First(&user)
	if result.Error != nil {
		fmt.Println("TODO GetUser(ctx context.Context, uId uuid.UUID) user.User")
	}
	return user
}

func (userRepository UserRepositoryPostgres) UpdateUser(ctx context.Context, user user.User) user.User {
	result := userRepository.database.WithContext(ctx).Save(&user)
	if result.Error != nil {
		fmt.Println("TODO UpdateUser(ctx context.Context, user user.User) user.User")
	}
	return user
}
