package user

func UserToCreateUserRequestDto(u User) CreateUserRequestDto {
	return CreateUserRequestDto{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.PasswordHash,
		Role:     u.Role,
	}
}

func UserToCreateUserResponseDto(u User) CreateUserResponseDto {
	return CreateUserResponseDto{
		Name:  u.Name,
		Email: u.Email,
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

func UserToUpdateUserResponseDto(u User) UpdateUserResponseDto {
	return UpdateUserResponseDto{
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

func UserToDeleteUserResponseDto(u User) DeleteUserResponseDto {
	return DeleteUserResponseDto{
		Name:  u.Name,
		Email: u.Email,
	}
}
