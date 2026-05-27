package user

import "context"

type UserService struct {
	Repository UserRepository
}

func (userService UserService) CreateUser(ctx context.Context, user User) User {
	return userService.Repository.CreateUser(ctx, user)
}
func (userService UserService) SoftDeleteUser(ctx context.Context, user User) User {
	return userService.Repository.SoftDeleteUser(ctx, user)
}

func (userService UserService) HardDeleteUser(ctx context.Context, user User) User {
	return userService.Repository.HardDeleteUser(ctx, user)
}

func (userService UserService) RestoreUser(ctx context.Context, user User) User {
	return userService.Repository.RestoreUser(ctx, user)
}

func (userService UserService) GetUser(ctx context.Context, user User) User {
	return userService.Repository.GetUser(ctx, user)
}

func (userService UserService) UpdateUser(ctx context.Context, user User) User {
	return userService.Repository.UpdateUser(ctx, user)
}
