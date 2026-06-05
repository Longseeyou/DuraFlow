package user

import "context"

type UserService struct {
	Repository UserRepository
}

func (userService UserService) CreateUser(ctx context.Context, u User) UserResponseDto {
	u, error := userService.Repository.CreateUser(ctx, u)
	if error != nil {
		// Handle error
	}
	return UserToUserResponseDto(u)
}
func (userService UserService) SoftDeleteUser(ctx context.Context, u User) UserResponseDto {
	u, error := userService.Repository.SoftDeleteUser(ctx, u)
	if error != nil {
		// Handle error
	}
	return UserToUserResponseDto(u)
}

func (userService UserService) HardDeleteUser(ctx context.Context, u User) UserResponseDto {
	u, error := userService.Repository.HardDeleteUser(ctx, u)
	if error != nil {
		// Handle error
	}
	return UserToUserResponseDto(u)
}

func (userService UserService) RestoreUser(ctx context.Context, u User) UserResponseDto {
	u, error := userService.Repository.RestoreUser(ctx, u)
	if error != nil {
		// Handle error
	}
	return UserToUserResponseDto(u)
}

func (userService UserService) GetUser(ctx context.Context, u User) UserRequestDto {
	u, error := userService.Repository.GetUser(ctx, u)
	if error != nil {
		// Handle error
	}
	return UserToUserRequestDto(u)
}

func (userService UserService) UpdateUser(ctx context.Context, u User) UserResponseDto {
	u, error := userService.Repository.UpdateUser(ctx, u)
	if error != nil {
		// Handle error
	}
	return UserToUserResponseDto(u)
}
