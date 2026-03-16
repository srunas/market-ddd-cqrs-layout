package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/category"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
	"github.com/srunas/market-ddd-cqrs-layout/internal/infrastructure/repository/sqlcgen"
)

type CategoryRepository struct {
	q *sqlcgen.Queries
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{q: sqlcgen.New(stdlib.OpenDBFromPool(pool))}
}

func (r *CategoryRepository) Save(ctx context.Context, c *category.Category) error {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	var parentID uuid.NullUUID
	if c.ParentID != nil {
		parentID = uuid.NullUUID{UUID: uuid.UUID(*c.ParentID), Valid: true}
	}
	return r.q.SaveCategory(ctx, sqlcgen.SaveCategoryParams{
		ID:        uuid.UUID(c.ID),
		Name:      c.Name,
		ParentID:  parentID,
		Level:     int32(c.Level), //nolint:gosec // level дерева категорий не может превысить int32
		CreatedAt: c.CreatedAt,
	})
}

func (r *CategoryRepository) FindByID(ctx context.Context, id types.CategoryID) (*category.Category, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	row, err := r.q.FindCategoryByID(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}
	return toCategoryDomain(row), nil
}

func (r *CategoryRepository) FindByIDForUpdate(ctx context.Context, id types.CategoryID) (*category.Category, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	row, err := r.q.FindCategoryByIDForUpdate(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}
	return toCategoryDomain(row), nil
}

func (r *CategoryRepository) FindAll(ctx context.Context) ([]*category.Category, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	rows, err := r.q.FindAllCategories(ctx)
	if err != nil {
		return nil, err
	}
	categories := make([]*category.Category, 0, len(rows))
	for _, row := range rows {
		categories = append(categories, toCategoryDomain(row))
	}
	return categories, nil
}

func toCategoryDomain(row sqlcgen.Category) *category.Category {
	var parentID *types.CategoryID
	if row.ParentID.Valid {
		id := types.CategoryID(row.ParentID.UUID)
		parentID = &id
	}
	return &category.Category{
		ID:        types.CategoryID(row.ID),
		Name:      row.Name,
		ParentID:  parentID,
		Level:     int(row.Level),
		CreatedAt: row.CreatedAt,
		Status:    category.StatusActive,
	}
}
