package service

import (
	"context"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type CartItem struct {
	ProductID types.ProductID
	Quantity  int64
}

type AddItemRequest struct {
	BuyerID   types.UserID
	ProductID types.ProductID
	Quantity  int64
}

type AddItemResponse struct{}

type RemoveItemRequest struct {
	BuyerID   types.UserID
	ProductID types.ProductID
}

type RemoveItemResponse struct{}

type DecreaseQuantityRequest struct {
	BuyerID   types.UserID
	ProductID types.ProductID
	Quantity  int64
}

type DecreaseQuantityResponse struct{}

type GetCartRequest struct {
	BuyerID types.UserID
}

type GetCartResponse struct {
	ID    types.CartID
	Items []CartItem
}

type Cart interface {
	AddItem(ctx context.Context, req AddItemRequest) (AddItemResponse, error)
	RemoveItem(ctx context.Context, req RemoveItemRequest) (RemoveItemResponse, error)
	DecreaseQuantity(ctx context.Context, req DecreaseQuantityRequest) (DecreaseQuantityResponse, error)
	GetCart(ctx context.Context, req GetCartRequest) (GetCartResponse, error)
}
