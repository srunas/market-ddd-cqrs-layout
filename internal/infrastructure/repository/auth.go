package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/auth"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
	"github.com/srunas/market-ddd-cqrs-layout/internal/infrastructure/repository/sqlcgen"
)

type AuthRepository struct {
	q *sqlcgen.Queries
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{q: sqlcgen.New(stdlib.OpenDBFromPool(pool))}
}

func (r *AuthRepository) Save(ctx context.Context, a *auth.Auth) error {
	return r.q.Save(ctx, sqlcgen.SaveParams{
		ID:        uuid.UUID(a.ID),
		Username:  a.Username(),
		Password:  a.Password(),
		CreatedAt: a.CreatedAt,
		AuthAt:    sql.NullTime{Valid: false},
	})
}

func (r *AuthRepository) UpdateAuth(ctx context.Context, a *auth.Auth) error {
	return r.q.UpdateAuth(ctx, sqlcgen.UpdateAuthParams{
		Username: a.Username(),
		AuthAt:   sql.NullTime{Valid: true, Time: a.AuthAt},
	})
}

func (r *AuthRepository) FindByUsername(ctx context.Context, username string) (*auth.Auth, error) {
	row, err := r.q.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return toAuthDomain(row), nil
}

func toAuthDomain(row sqlcgen.Auth) *auth.Auth {
	return auth.NewFromDB(
		types.AuthID(row.ID),
		row.Username,
		row.Password,
		row.AuthAt.Time,
		row.CreatedAt,
	)
}
