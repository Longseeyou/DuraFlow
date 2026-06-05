package user

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (User, error)
	SoftDeleteUser(ctx context.Context, user User) (User, error)
	HardDeleteUser(ctx context.Context, user User) (User, error)
	RestoreUser(ctx context.Context, user User) (User, error)
	GetUser(ctx context.Context, user User) (User, error)
	UpdateUser(ctx context.Context, user User, newUser map[string]any) (User, error)
}
