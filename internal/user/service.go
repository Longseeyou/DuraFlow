package user

import "context"

type UserService struct {
	Repository UserRepository
}

func (userService UserService) CreateUser(ctx context.Context, u CreateUserRequestDto) UserResponseDto {
	userModel := CreateUserRequestDtoToUser(u)
	userModel, err := userService.Repository.CreateUser(ctx, userModel)
	if err != nil {
		// Handle error
	}
	return UserToUserResponseDto(userModel)
}
func (userService UserService) SoftDeleteUser(ctx context.Context, u DeleteUserRequestDto) UserResponseDto {
	userModel := DeleteUserRequestDtoToUser(u)
	userModel, err := userService.Repository.SoftDeleteUser(ctx, userModel)
	if err != nil {
		// Handle error
	}
	return UserToUserResponseDto(userModel)
}

func (userService UserService) HardDeleteUser(ctx context.Context, u DeleteUserRequestDto) UserResponseDto {
	userModel := DeleteUserRequestDtoToUser(u)
	userModel, err := userService.Repository.HardDeleteUser(ctx, userModel)
	if err != nil {
		// Handle error
	}
	return UserToUserResponseDto(userModel)
}

func (userService UserService) RestoreUser(ctx context.Context, u UserRequestDto) UserResponseDto {
	userModel := UserRequestDtoToUser(u)
	userModel, err := userService.Repository.RestoreUser(ctx, userModel)
	if err != nil {
		// Handle error
	}
	return UserToUserResponseDto(userModel)
}

func (userService UserService) GetUser(ctx context.Context, u UserRequestDto) UserResponseDto {
	userModel := UserRequestDtoToUser(u)
	userModel, err := userService.Repository.GetUser(ctx, userModel)
	if err != nil {
		// Handle error
	}
	return UserToUserResponseDto(userModel)
}

func (userService UserService) UpdateUser(ctx context.Context, u UpdateUserRequestDto) UserResponseDto {
	userModel := UpdateUserRequestDtoToUser(u)
	newUser := UpdateUserRequestDtoToMap(u)
	userModel, err := userService.Repository.UpdateUser(ctx, userModel, newUser)
	if err != nil {
		// Handle error
	}
	return UserToUserResponseDto(userModel)
}
