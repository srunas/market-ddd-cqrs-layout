package order

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/openpgp/errors"
)

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "created"
	OrderStatusProcessed OrderStatus = "processed"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"
	OrderStatusFailed    OrderStatus = "failed"
)

type PaymentMethod string

const (
	PaymentMethodCard PaymentMethod = "card"
)

type Currency string

type Order struct {
	ID            uuid.UUID
	BuyerID       uuid.UUID
	Status        OrderStatus
	Total         decimal.Decimal
	Currency      Currency
	PaymentMethod PaymentMethod
	CreatedAt     time.Time
	Items         []OrderItem
}

type OrderItem struct {
	ProductID    uuid.UUID
	Quantity     decimal.Decimal
	PriceAtOrder decimal.Decimal
}

func NewOrder(buyerID uuid.UUID, currency Currency, method PaymentMethod) *Order {
	return &Order{
		ID:            uuid.New(),
		BuyerID:       buyerID,
		Status:        OrderStatusCreated,
		Currency:      currency,
		PaymentMethod: method,
		CreatedAt:     time.Now(),
		Items:         []OrderItem{},
	}
}

func (o *Order) AddItem(productID uuid.UUID, quantity decimal.Decimal, price decimal.Decimal) error {
	if quantity.IsZero() {
		return fmt.Errorf("quantity must be positive")
	}
	o.Items = append(o.Items, OrderItem{
		ProductID:    productID,
		Quantity:     quantity,
		PriceAtOrder: price,
	})
	o.Total = o.Total.Add(price.Mul(decimal.NewFromInt(int64(quantity))))
	return nil
}

func (o *Order) Process() error {
	if o.Status != OrderStatusCreated {
		return fmt.Errorf("can only process order with status created")
	}
	o.Status = OrderStatusProcessed
	return nil
}

func (o *Order) Complete(success bool) error {
	if o.Status != OrderStatusProcessed {
		return fmt.Errorf("can only complete processing orders")
	}
	if success {
		o.Status = OrderStatusCompleted
	} else {
		o.Status = OrderStatusFailed
	}
	return nil
}
