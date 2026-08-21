package user

type CreateUserRequestDto struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6"`
	RetypedPassword string `json:"retyped_password" validate:"required,eqfield=Password"`
	Role            Role   `json:"role"`
}

type UpdateUserRequestDto struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
	Role     *Role   `json:"role,omitempty"`
}

type UserResponseDto struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
