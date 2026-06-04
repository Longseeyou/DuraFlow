package user

import "context"

type UserService struct {
	Repository UserRepository
}

func (userService UserService) CreateUser(ctx context.Context, user User) CreateUserResponseDto {
	user, error := userService.Repository.CreateUser(ctx, user)
	if error != nil {
		// Handle error
	}
	return UserToCreateUserResponseDto(user)
}
func (userService UserService) SoftDeleteUser(ctx context.Context, user User) DeleteUserResponseDto {
	user, error := userService.Repository.SoftDeleteUser(ctx, user)
	if error != nil {
		// Handle error
	}
	return UserToDeleteUserResponseDto(user)
}

func (userService UserService) HardDeleteUser(ctx context.Context, user User) DeleteUserResponseDto {
	user, error := userService.Repository.HardDeleteUser(ctx, user)
	if error != nil {
		// Handle error
	}
	return UserToDeleteUserResponseDto(user)
}

func (userService UserService) RestoreUser(ctx context.Context, user User) UserResponseDto {
	user, error := userService.Repository.RestoreUser(ctx, user)
	if error != nil {
		// Handle error
	}
	return UserToUserResponseDto(user)
}

func (userService UserService) GetUser(ctx context.Context, user User) UserRequestDto {
	user, error := userService.Repository.GetUser(ctx, user)
	if error != nil {
		// Handle error
	}
	return UserToUserRequestDto(user)
}

func (userService UserService) UpdateUser(ctx context.Context, user User) UserResponseDto {
	user, error := userService.Repository.UpdateUser(ctx, user)
	if error != nil {
		// Handle error
	}
	return UserToUserResponseDto(user)
}
