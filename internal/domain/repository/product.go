package repository

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/product"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type ProductFilter struct {
	CategoryID *types.CategoryID
	Status     *product.Status
	MinPrice   *decimal.Decimal
	MaxPrice   *decimal.Decimal
	Limit      int
	Offset     int
}

type Product interface {
	Save(ctx context.Context, product *product.Product) error
	Update(ctx context.Context, product *product.Product) error
	FindByID(ctx context.Context, id types.ProductID) (*product.Product, error)
	GetProductList(ctx context.Context, filter ProductFilter) ([]*product.Product, error)
}
