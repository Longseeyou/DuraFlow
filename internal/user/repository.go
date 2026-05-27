package user

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, user User) User
	SoftDeleteUser(ctx context.Context, user User) User
	HardDeleteUser(ctx context.Context, user User) User
	RestoreUser(ctx context.Context, user User) User
	GetUser(ctx context.Context, user User) User
	UpdateUser(ctx context.Context, user User) User
}
