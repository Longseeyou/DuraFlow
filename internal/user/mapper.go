package user

func UserToCreateUserRequestDto(u User) CreateUserRequestDto {
	return CreateUserRequestDto{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.PasswordHash,
		Role:     u.Role,
	}
}

func UserToUpdateUserRequestDto(u User) UpdateUserRequestDto {
	return UpdateUserRequestDto{
		Name:     &u.Name,
		Email:    &u.Email,
		Password: &u.PasswordHash,
		Role:     &u.Role,
	}
}

func UserToUpdateUserResponseDto(u User) UserResponseDto {
	return UserResponseDto{
		Name:  u.Name,
		Email: u.Email,
	}
}

func UserToUserRequestDto(u User) UserRequestDto {
	return UserRequestDto{
		ID:    &u.ID,
		Name:  &u.Name,
		Email: &u.Email,
	}
}

func UserToUserResponseDto(u User) UserResponseDto {
	return UserResponseDto{
		Name:  u.Name,
		Email: u.Email,
	}
}

func UserToDeleteUserRequestDto(u User) DeleteUserRequestDto {
	return DeleteUserRequestDto{
		ID: &u.ID,
	}
}

func UserToDeleteUserResponseDto(u User) UserResponseDto {
	return UserResponseDto{
		Name:  u.Name,
		Email: u.Email,
	}
}
