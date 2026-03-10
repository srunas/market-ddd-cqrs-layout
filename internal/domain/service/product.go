package service

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/product"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type ProductItem struct {
	ID          types.ProductID
	Name        string
	Price       decimal.Decimal
	Status      product.Status
	CategoryIDs []types.CategoryID
}

type CreateProductRequest struct {
	Name        string
	Description string
	Price       decimal.Decimal
	Currency    product.Currency
	SellerID    types.UserID
	CategoryIDs []types.CategoryID
}

type CreateProductResponse struct {
	ID types.ProductID
}

type PublishProductRequest struct {
	ID types.ProductID
}

type PublishProductResponse struct {
}

type ListProductsRequest struct {
	CategoryID *types.CategoryID
	Status     *product.Status
	MinPrice   *decimal.Decimal
	MaxPrice   *decimal.Decimal
	Page       int
	PageSize   int
}

type ListProductsResponse struct {
	Items []*ProductItem
}

type Product interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (CreateProductResponse, error)
	PublishProduct(ctx context.Context, req PublishProductRequest) (PublishProductResponse, error)
	ListProducts(ctx context.Context, req ListProductsRequest) (ListProductsResponse, error)
}
