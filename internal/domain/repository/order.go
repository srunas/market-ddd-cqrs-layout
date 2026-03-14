package repository

import (
	"context"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/order"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type Order interface {
	Save(ctx context.Context, order *order.Order) error
	FindByID(ctx context.Context, id types.OrderID) (*order.Order, error)
	FindByBuyerID(ctx context.Context, id types.UserID) ([]*order.Order, error)
	Update(ctx context.Context, order *order.Order) error
}
