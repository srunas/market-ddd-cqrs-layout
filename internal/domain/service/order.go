package service

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/order"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type PlaceOrderRequest struct {
	BuyerID       types.UserID
	PaymentMethod order.PaymentMethod
	Currency      order.Currency
}

type PlaceOrderResponse struct {
	Order *order.Order
}

type CancelOrderRequest struct {
	OrderID types.OrderID
	BuyerID types.UserID
}

type CancelOrderResponse struct {
}

type GetOrderRequest struct {
	OrderID types.OrderID
}

type GetOrderResponse struct {
	Order *order.Order
}

type ListOrdersRequest struct {
	BuyerID types.UserID
}

type OrderSummary struct {
	ID     types.OrderID
	Status order.Status
	Total  decimal.Decimal
}

type ListOrdersResponse struct {
	OrderSummary []*OrderSummary
}

type Order interface {
	PlaceOrder(ctx context.Context, req PlaceOrderRequest) (PlaceOrderResponse, error)
	CancelOrder(ctx context.Context, req CancelOrderRequest) (CancelOrderResponse, error)
	GetOrder(ctx context.Context, req GetOrderRequest) (GetOrderResponse, error)
	ListOrders(ctx context.Context, req ListOrdersRequest) (ListOrdersResponse, error)
}
