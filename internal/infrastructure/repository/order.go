package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/shopspring/decimal"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/order"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
	"github.com/srunas/market-ddd-cqrs-layout/internal/infrastructure/repository/sqlcgen"
)

type OrderRepository struct {
	q *sqlcgen.Queries
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{q: sqlcgen.New(stdlib.OpenDBFromPool(pool))}
}

func (r *OrderRepository) Save(ctx context.Context, o *order.Order) error {
	err := r.q.SaveOrder(ctx, sqlcgen.SaveOrderParams{
		ID:            uuid.UUID(o.ID),
		BuyerID:       uuid.UUID(o.BuyerID),
		Status:        sqlcgen.OrderStatus(o.Status),
		Total:         o.Total.String(),
		Currency:      string(o.Currency),
		PaymentMethod: sqlcgen.PaymentMethod(o.PaymentMethod),
		CreatedAt:     o.CreatedAt,
	})
	if err != nil {
		return err
	}

	for _, item := range o.Items {
		err = r.q.SaveOrderItem(ctx, sqlcgen.SaveOrderItemParams{
			OrderID:      uuid.UUID(o.ID),
			ProductID:    uuid.UUID(item.ProductID),
			Quantity:     int32(item.Quantity), //nolint:gosec // quantity не может превысить int32
			PriceAtOrder: item.PriceAtOrder.String(),
		})
		if err != nil {
			return fmt.Errorf("ошибка сохранения позиции заказа: %w", err)
		}
	}
	return nil
}

func (r *OrderRepository) FindByID(ctx context.Context, id types.OrderID) (*order.Order, error) {
	row, err := r.q.FindOrderByID(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}
	return r.toOrderDomain(ctx, row)
}

func (r *OrderRepository) FindByBuyerID(ctx context.Context, buyerID types.UserID) ([]*order.Order, error) {
	rows, err := r.q.FindOrdersByBuyerID(ctx, uuid.UUID(buyerID))
	if err != nil {
		return nil, err
	}

	orders := make([]*order.Order, 0, len(rows))
	for _, row := range rows {
		o, errOrder := r.toOrderDomain(ctx, row)
		if errOrder != nil {
			return nil, errOrder
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *OrderRepository) Update(ctx context.Context, o *order.Order) error {
	return r.q.UpdateOrder(ctx, sqlcgen.UpdateOrderParams{
		ID:     uuid.UUID(o.ID),
		Status: sqlcgen.OrderStatus(o.Status),
	})
}

func (r *OrderRepository) toOrderDomain(ctx context.Context, row sqlcgen.Order) (*order.Order, error) {
	itemRows, err := r.q.FindOrderItems(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	items := make([]order.Item, 0, len(itemRows))
	for _, item := range itemRows {
		price, errOrder := decimal.NewFromString(item.PriceAtOrder)
		if errOrder != nil {
			return nil, errOrder
		}
		items = append(items, order.Item{
			ProductID:    types.ProductID(item.ProductID),
			Quantity:     int64(item.Quantity),
			PriceAtOrder: price,
		})
	}

	total, err := decimal.NewFromString(row.Total)
	if err != nil {
		return nil, err
	}

	return &order.Order{
		ID:            types.OrderID(row.ID),
		BuyerID:       types.UserID(row.BuyerID),
		Status:        order.Status(row.Status),
		Total:         total,
		Currency:      order.Currency(row.Currency),
		PaymentMethod: order.PaymentMethod(row.PaymentMethod),
		CreatedAt:     row.CreatedAt,
		Items:         items,
	}, nil
}
