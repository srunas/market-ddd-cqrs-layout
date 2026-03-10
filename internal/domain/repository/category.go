package repository

import (
	"context"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/category"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type Category interface {
	Save(ctx context.Context, category *category.Category) error
	FindByID(ctx context.Context, id types.CategoryID) (*category.Category, error)
	FindAll(ctx context.Context) ([]*category.Category, error)
	FindByIDForUpdate(ctx context.Context, id types.CategoryID) (*category.Category, error)
}
