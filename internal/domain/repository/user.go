package repository

import (
	"context"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/user"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type User interface {
	Save(ctx context.Context, user *user.User) error
	FindByID(ctx context.Context, id types.UserID) (*user.User, error)
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
}
