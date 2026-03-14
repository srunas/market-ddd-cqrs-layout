package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/user"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
	"github.com/srunas/market-ddd-cqrs-layout/internal/infrastructure/repository/sqlcgen"
)

type UserRepository struct {
	q *sqlcgen.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{q: sqlcgen.New(stdlib.OpenDBFromPool(pool))}
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	return r.q.SaveUser(ctx, sqlcgen.SaveUserParams{
		ID:        uuid.UUID(u.ID),
		Username:  u.Username,
		Surname:   sql.NullString{String: u.Surname, Valid: u.Surname != ""},
		Role:      sqlcgen.UserRole(u.Role),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		Enabled:   u.Enabled,
	})
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	return r.q.UpdateUser(ctx, sqlcgen.UpdateUserParams{
		ID:       uuid.UUID(u.ID),
		Username: u.Username,
		Surname:  sql.NullString{String: u.Surname, Valid: u.Surname != ""},
		Email:    u.Email,
		Enabled:  u.Enabled,
	})
}

func (r *UserRepository) FindByID(ctx context.Context, id types.UserID) (*user.User, error) {
	row, err := r.q.FindUserByID(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}
	return toUserDomain(row), nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	row, err := r.q.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return toUserDomain(row), nil
}

func toUserDomain(row sqlcgen.User) *user.User {
	return &user.User{
		ID:        types.UserID(row.ID),
		Username:  row.Username,
		Surname:   row.Surname.String,
		Role:      user.Role(row.Role),
		Email:     row.Email,
		Enabled:   row.Enabled,
		CreatedAt: row.CreatedAt,
		AuthAt:    row.AuthAt.Time,
	}
}
