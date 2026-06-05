package user

func CreateUserRequestDtoToUser(uDto CreateUserRequestDto) User {
	return User{
		Name:         uDto.Name,
		Email:        uDto.Email,
		PasswordHash: uDto.Password,
		Role:         uDto.Role,
		IsActive:     true,
	}
}

func UserRequestDtoToUser(uDto UserRequestDto) User {
	u := User{}
	if uDto.ID != nil {
		u.ID = *uDto.ID
	}
	if uDto.Name != nil {
		u.Name = *uDto.Name
	}
	if uDto.Email != nil {
		u.Email = *uDto.Email
	}
	return u
}

func DeleteUserRequestDtoToUser(uDto DeleteUserRequestDto) User {
	u := User{}
	if uDto.ID != nil {
		u.ID = *uDto.ID
	}
	return u
}

func UpdateUserRequestDtoToUser(uDto UpdateUserRequestDto) User {
	u := User{}
	if uDto.ID != nil {
		u.ID = *uDto.ID
	}
	return u
}

func UpdateUserRequestDtoToMap(uDto UpdateUserRequestDto) map[string]any {
	newUser := map[string]any{}
	if uDto.Name != nil {
		newUser["name"] = *uDto.Name
	}
	if uDto.Email != nil {
		newUser["email"] = *uDto.Email
	}
	if uDto.Password != nil {
		newUser["password_hash"] = *uDto.Password
	}
	if uDto.Role != nil {
		newUser["role"] = *uDto.Role
	}
	return newUser
}

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
