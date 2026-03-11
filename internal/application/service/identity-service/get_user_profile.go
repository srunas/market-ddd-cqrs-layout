package identity_service

import (
	"context"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

// GetUserProfile возвращает профиль пользователя и его ID.
func (s *Implementation) GetUserProfile(
	ctx context.Context,
	req service.GetUserProfileRequest,
) (service.GetUserProfileResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	user, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return service.GetUserProfileResponse{}, err
	}

	return service.GetUserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Surname:  user.Surname,
		Email:    user.Email,
		Role:     string(user.Role),
	}, nil
}
