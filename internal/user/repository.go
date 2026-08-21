package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	UserStore
}

type UserStore interface {
	CreateUser(ctx context.Context, u User) (User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	UpdateUserByID(ctx context.Context, userID uuid.UUID, newUser map[string]any) (User, error)
	SoftDeleteUser(ctx context.Context, userID uuid.UUID) (User, error)
	HardDeleteUser(ctx context.Context, userID uuid.UUID) (User, error)
	RestoreUser(ctx context.Context, userID uuid.UUID) (User, error)
}
