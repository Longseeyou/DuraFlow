package user

func UserToUserResponseDto(u User) UserResponseDto {
	return UserResponseDto{
		Name:  u.Name,
		Email: u.Email,
	}
}
