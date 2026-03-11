package identity_service

import (
	"context"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/auth"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/user"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) RegisterUser(
	ctx context.Context,
	req service.RegisterUserRequest,
) (service.RegisterUserResponse, error) {
	ctx, cpan := prospan.Start(ctx)
	defer cpan.End()

	var userEntity *user.User

	err := s.txManager.Do(ctx, func(ctx context.Context) error {
		authEntity, err := auth.New(req.Username, req.Password)
		if err != nil {
			return err
		}

		userEntity, err = user.New(req.Username, req.Surname, req.Email, user.RoleBuyer)
		if err != nil {
			return err
		}

		if err = s.authRepo.Save(ctx, authEntity); err != nil {
			return err
		}
		return s.userRepo.Save(ctx, userEntity)
	})
	if err != nil {
		return service.RegisterUserResponse{}, err
	}
	return service.RegisterUserResponse{UserID: userEntity.ID}, nil
}
