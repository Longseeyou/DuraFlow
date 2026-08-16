package user

import (
	"context"
	"errors"

	"github.com/Longseeyou/DuraFlow/internal/shared/custom_error"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordsDoNotMatch = errors.New("password and retyped password do not match")

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	dto CreateUserRequestDto,
) (UserResponseDto, error) {
	if dto.Password != dto.RetypedPassword {
		return UserResponseDto{}, ErrPasswordsDoNotMatch
	}

	u := User{
		Name:     dto.Name,
		Email:    dto.Email,
		Role:     USER,
		IsActive: true,
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponseDto{}, err
	}
	u.PasswordHash = string(passwordHash)

	u, err = s.repository.CreateUser(ctx, u)
	if err != nil {
		return UserResponseDto{}, err
	}

	return UserToUserResponseDto(u), nil
}

func (s *UserService) GetUserByID(
	ctx context.Context,
	userID uuid.UUID,
) (UserResponseDto, error) {
	u, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return UserResponseDto{}, err
	}

	return UserToUserResponseDto(u), nil
}

func (s *UserService) UpdateUserByID(
	ctx context.Context,
	userID uuid.UUID,
	dto UpdateUserRequestDto,
) (UserResponseDto, error) {
	newUser := map[string]any{}
	if dto.Name != nil {
		newUser["name"] = *dto.Name
	}
	if dto.Email != nil {
		newUser["email"] = *dto.Email
	}
	if dto.Password != nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*dto.Password), bcrypt.DefaultCost)
		if err != nil {
			return UserResponseDto{}, err
		}
		newUser["password_hash"] = string(passwordHash)
	}
	if dto.Role != nil {
		newUser["role"] = *dto.Role
	}

	if len(newUser) == 0 {
		return UserResponseDto{}, custom_error.ErrUpdateInvalidRequest
	}

	u, err := s.repository.UpdateUserByID(ctx, userID, newUser)
	if err != nil {
		return UserResponseDto{}, err
	}

	return UserToUserResponseDto(u), nil
}

func (s *UserService) SoftDeleteUser(
	ctx context.Context,
	userID uuid.UUID,
) (UserResponseDto, error) {
	u, err := s.repository.SoftDeleteUser(ctx, userID)
	if err != nil {
		return UserResponseDto{}, err
	}

	return UserToUserResponseDto(u), nil
}

func (s *UserService) HardDeleteUser(
	ctx context.Context,
	userID uuid.UUID,
) (UserResponseDto, error) {
	u, err := s.repository.HardDeleteUser(ctx, userID)
	if err != nil {
		return UserResponseDto{}, err
	}

	return UserToUserResponseDto(u), nil
}

func (s *UserService) RestoreUser(
	ctx context.Context,
	userID uuid.UUID,
) (UserResponseDto, error) {
	u, err := s.repository.RestoreUser(ctx, userID)
	if err != nil {
		return UserResponseDto{}, err
	}

	return UserToUserResponseDto(u), nil
}
