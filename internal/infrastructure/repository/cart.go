package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/cart"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
	"github.com/srunas/market-ddd-cqrs-layout/internal/infrastructure/repository/sqlcgen"
)

type CartRepository struct {
	q *sqlcgen.Queries
}

func NewCartRepository(pool *pgxpool.Pool) *CartRepository {
	return &CartRepository{q: sqlcgen.New(stdlib.OpenDBFromPool(pool))}
}

func (r *CartRepository) Save(ctx context.Context, c *cart.Cart) error {
	err := r.q.SaveCart(ctx, sqlcgen.SaveCartParams{
		ID:        uuid.UUID(c.ID),
		BuyerID:   uuid.UUID(c.BuyerID),
		CreatedAt: c.CreatedAt,
	})
	if err != nil {
		return err
	}
	for _, item := range c.Items {
		err = r.q.SaveCartItem(ctx, sqlcgen.SaveCartItemParams{
			CartID:    uuid.UUID(c.ID),
			ProductID: uuid.UUID(item.ProductID),
			Quantity:  int32(item.Quantity), //nolint:gosec // quantity не может превысить int32
		})
		if err != nil {
			return fmt.Errorf("ошибка сохранения позиции корзины: %w", err)
		}
	}
	return nil
}

func (r *CartRepository) Update(ctx context.Context, c *cart.Cart) error {
	// Сохраняем каждую позицию — ON CONFLICT DO UPDATE обновит quantity
	for _, item := range c.Items {
		err := r.q.SaveCartItem(ctx, sqlcgen.SaveCartItemParams{
			CartID:    uuid.UUID(c.ID),
			ProductID: uuid.UUID(item.ProductID),
			Quantity:  int32(item.Quantity), //nolint:gosec // quantity не может превысить int32
		})
		if err != nil {
			return fmt.Errorf("ошибка обновления позиции корзины: %w", err)
		}
	}
	return nil
}

func (r *CartRepository) FindByBuyerID(ctx context.Context, buyerID types.UserID) (*cart.Cart, error) {
	row, err := r.q.FindCartByBuyerID(ctx, uuid.UUID(buyerID))
	if err != nil {
		return nil, err
	}

	itemRows, err := r.q.FindCartItems(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	items := make([]cart.CartItem, 0, len(itemRows))
	for _, item := range itemRows {
		items = append(items, cart.CartItem{
			ProductID: types.ProductID(item.ProductID),
			Quantity:  int64(item.Quantity),
		})
	}

	return &cart.Cart{
		ID:        types.CartID(row.ID),
		BuyerID:   types.UserID(row.BuyerID),
		CreatedAt: row.CreatedAt,
		Items:     items,
	}, nil
}
