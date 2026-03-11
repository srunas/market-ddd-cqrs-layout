package repository

import (
	"context"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/auth"
)

type Auth interface {
	Save(ctx context.Context, a *auth.Auth) error
	FindByUsername(ctx context.Context, username string) (*auth.Auth, error)
	UpdateAuth(ctx context.Context, a *auth.Auth) error
}
