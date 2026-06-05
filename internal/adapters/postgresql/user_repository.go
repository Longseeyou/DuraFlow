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

func NewUserRepository(database *gorm.DB) user.UserRepository {
	return UserRepositoryPostgres{database: database}
}

func (userRepository UserRepositoryPostgres) CreateUser(ctx context.Context, u user.User) (user.User, error) {
	result := userRepository.database.WithContext(ctx).Create(&u)
	if result.Error != nil {
		fmt.Println("TODO CreateUser(ctx context.Context, user user.User) user.User")
		return u, result.Error
	}
	return u, nil
}

func (userRepository UserRepositoryPostgres) SoftDeleteUser(ctx context.Context, u user.User) (user.User, error) {
	result := userRepository.database.WithContext(ctx).Delete(&u)
	if result.Error != nil {
		fmt.Println("TODO SoftDeleteUser(ctx context.Context, user user.User) user.User")
		return u, result.Error
	}

	return u, nil
}

func (userRepository UserRepositoryPostgres) HardDeleteUser(ctx context.Context, u user.User) (user.User, error) {
	result := userRepository.database.WithContext(ctx).Unscoped().Delete(&u)
	if result.Error != nil {
		fmt.Println("TODO HardDeleteUser(ctx context.Context, user user.User) user.User")
		return u, result.Error
	}
	return u, nil
}

func (userRepository UserRepositoryPostgres) RestoreUser(ctx context.Context, u user.User) (user.User, error) {
	result := userRepository.database.WithContext(ctx).Unscoped().Model(&u).Where("id = ?", u.ID).Update("deleted_at", nil)
	if result.Error != nil {
		fmt.Println("TODO RestoreUser(ctx context.Context, user user.User) user.User")
		return u, result.Error
	}
	return u, nil
}

func (userRepository UserRepositoryPostgres) GetUser(ctx context.Context, u user.User) (user.User, error) {
	result := userRepository.database.WithContext(ctx).Where("id = ?", u.ID).First(&u)
	if result.Error != nil {
		fmt.Println("TODO GetUser(ctx context.Context, uId uuid.UUID) user.User")
		return u, result.Error
	}
	return u, nil
}

func (userRepository UserRepositoryPostgres) UpdateUser(ctx context.Context, u user.User, newUser map[string]any) (user.User, error) {
	result := userRepository.database.WithContext(ctx).Model(&u).Where("id = ?", u.ID).Updates(newUser)
	if result.Error != nil {
		fmt.Println("TODO UpdateUser(ctx context.Context, user user.User) user.User")
		return u, result.Error
	}
	result = userRepository.database.WithContext(ctx).Where("id = ?", u.ID).First(&u)
	if result.Error != nil {
		fmt.Println("TODO UpdateUser(ctx context.Context, user user.User) user.User")
		return u, result.Error
	}
	return u, nil
}
