package identity_service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

// tokenTTL — время жизни токена.
const tokenTTL = 24 * time.Hour

// Claims — данные которые кладём внутрь JWT токена.
// Включаем только то что нужно для авторизации — не храним пароли и личные данные.
type Claims struct {
	jwt.RegisteredClaims

	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// LoginUser проверяет credentials пользователя и возвращает JWT токен.
// Обновляет время последней авторизации при успешном входе.
func (s *Implementation) LoginUser(
	ctx context.Context,
	req service.LoginRequest,
) (service.LoginResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	authEntity, err := s.authRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return service.LoginResponse{}, errors.New("неверный username или пароль")
	}

	if !authEntity.ValidatePassword(req.Password) {
		return service.LoginResponse{}, errors.New("неверный username или пароль")
	}

	userEntity, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return service.LoginResponse{}, err
	}

	authEntity.UpdateAuthTime()
	if err = s.authRepo.UpdateAuth(ctx, authEntity); err != nil {
		return service.LoginResponse{}, err
	}

	token, err := s.generateToken(userEntity.ID, string(userEntity.Role))
	if err != nil {
		return service.LoginResponse{}, err
	}

	return service.LoginResponse{Token: token}, nil
}

func (s *Implementation) generateToken(userID types.UserID, role string) (string, error) {
	claims := Claims{
		UserID: uuid.UUID(userID).String(),
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "marketplace",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
