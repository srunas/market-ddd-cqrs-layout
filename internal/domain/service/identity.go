package service

import (
	"context"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type RegisterUserRequest struct {
	Username string
	Surname  string
	Email    string
	Password string //nolint:gosec // plain text пароль — хэшируется в auth.New через bcrypt
}

type RegisterUserResponse struct {
	UserID types.UserID
}

type LoginRequest struct {
	Username string
	Password string //nolint:gosec // plain text пароль — хэшируется в auth.New через bcrypt
}

type LoginResponse struct {
	Token string
}

type GetUserProfileRequest struct {
	UserID types.UserID
}

type GetUserProfileResponse struct {
	ID       types.UserID
	Username string
	Surname  string
	Email    string
	Role     string
}

type Identity interface {
	RegisterUser(ctx context.Context, r RegisterUserRequest) (RegisterUserResponse, error)
	LoginUser(ctx context.Context, req LoginRequest) (LoginResponse, error)
	GetUserProfile(ctx context.Context, req GetUserProfileRequest) (GetUserProfileResponse, error)
}
