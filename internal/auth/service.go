package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Longseeyou/DuraFlow/internal/user"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	defaultTokenTTL = 24 * time.Hour
	tokenType       = "Bearer"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type AuthService struct {
	Repository user.UserRepository
	TokenAuth  *jwtauth.JWTAuth
	TokenTTL   time.Duration
}

func NewAuthService(
	repository user.UserRepository,
	secret string,
	tokenTTL time.Duration,
) *AuthService {
	if tokenTTL == 0 {
		tokenTTL = defaultTokenTTL
	}

	return &AuthService{
		Repository: repository,
		TokenAuth: jwtauth.New(
			"HS256",
			[]byte(secret),
			nil,
			jwt.WithAcceptableSkew(30*time.Second),
		),
		TokenTTL: tokenTTL,
	}
}

func (authService AuthService) Login(
	ctx context.Context,
	request LoginRequestDto,
) (TokenResponseDto, error) {
	u, err := authService.Repository.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return TokenResponseDto{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.PasswordHash),
		[]byte(request.Password),
	); err != nil {
		return TokenResponseDto{}, ErrInvalidCredentials
	}

	ttl := authService.TokenTTL
	if ttl == 0 {
		ttl = defaultTokenTTL
	}

	claims := map[string]any{
		"sub":     u.ID.String(),
		"user_id": u.ID.String(),
		"email":   u.Email,
		"role":    u.Role,
		"exp":     jwtauth.ExpireIn(ttl),
	}
	jwtauth.SetIssuedNow(claims)

	_, tokenString, err := authService.TokenAuth.Encode(claims)
	if err != nil {
		return TokenResponseDto{}, err
	}

	return TokenResponseDto{
		TokenType:   tokenType,
		AccessToken: tokenString,
		ExpiresIn:   int64(ttl.Seconds()),
	}, nil
}
