package repository

import (
	"context"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/cart"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type Cart interface {
	Save(ctx context.Context, cart *cart.Cart) error
	Update(ctx context.Context, cart *cart.Cart) error
	FindByBuyerID(ctx context.Context, buyerID types.UserID) (*cart.Cart, error)
}
