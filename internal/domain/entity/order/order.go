package order

import (
	"errors"
	_ "fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/cart"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

var (
	ErrInvalidQuantity           = errors.New("quantity must be positive")
	ErrInvalidStatus             = errors.New("invalid order status for this operation")
	ErrNotInCreatedState         = errors.New("can only process order with status 'created'")
	ErrNotInProcessedState       = errors.New("can only complete processing orders")
	ErrEmptyCart                 = errors.New("cart is empty")
	ErrNilCart                   = errors.New("cart is nil")
	ErrItemAlreadyInOrder        = errors.New("item already exists")
	ErrCannotCancelFinishedOrder = errors.New("cannot cancel finished order")
)

type Status string

const (
	StatusCreated   Status = "created"
	StatusProcessed Status = "processed"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
)

type PaymentMethod string

const (
	PaymentCard   PaymentMethod = "card"
	PaymentCash   PaymentMethod = "cash"
	PaymentWallet PaymentMethod = "wallet"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
)

type Order struct {
	ID            types.OrderID
	BuyerID       types.UserID
	Status        Status
	Total         decimal.Decimal
	Currency      Currency
	PaymentMethod PaymentMethod
	CreatedAt     time.Time
	Items         []Item
}

type Item struct {
	ProductID    types.ProductID
	Quantity     int64
	PriceAtOrder decimal.Decimal
}

func New(buyerID types.UserID, currency Currency, method PaymentMethod) *Order {
	return &Order{
		ID:            types.NewOrderID(),
		BuyerID:       buyerID,
		Status:        StatusCreated,
		Currency:      currency,
		PaymentMethod: method,
		CreatedAt:     time.Now().UTC(),
		Items:         []Item{},
		Total:         decimal.Zero,
	}
}

func NewFromCart(cart *cart.Cart) (*Order, error) {
	if cart == nil {
		return nil, ErrNilCart
	}
	if len(cart.Items) == 0 {
		return nil, ErrEmptyCart
	}

	order := &Order{
		ID:            types.NewOrderID(),
		BuyerID:       cart.BuyerID,
		Status:        StatusCreated,
		Currency:      CurrencyUSD,
		PaymentMethod: PaymentCard,
		CreatedAt:     time.Now().UTC(),
		Items:         make([]Item, len(cart.Items)),
		Total:         decimal.Zero,
	}

	for _, ci := range cart.Items {
		unitPrice := decimal.NewFromInt(100) //nolint:mnd // placeholder price until real product service integration

		if ci.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}

		item := Item{
			ProductID:    ci.ProductID,
			Quantity:     ci.Quantity,
			PriceAtOrder: unitPrice,
		}

		order.Items = append(order.Items, item)
		order.Total = order.Total.Add(unitPrice.Mul(decimal.NewFromInt(ci.Quantity)))
	}
	return order, nil
}

func (o *Order) AddItem(productID types.ProductID, quantity int64, unitPrice decimal.Decimal) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	for _, item := range o.Items {
		if item.ProductID == productID {
			return errors.New("item already exists in order")
		}
	}

	o.Items = append(o.Items, Item{
		ProductID:    productID,
		Quantity:     quantity,
		PriceAtOrder: unitPrice,
	})

	itemTotal := unitPrice.Mul(decimal.NewFromInt(quantity))
	o.Total = o.Total.Add(itemTotal)

	return nil
}

func (o *Order) Process() error {
	if o.Status != StatusCreated {
		return ErrNotInCreatedState
	}
	o.Status = StatusProcessed
	return nil
}

func (o *Order) Complete(success bool) error {
	if o.Status != StatusProcessed {
		return ErrNotInProcessedState
	}
	if success {
		o.Status = StatusCompleted
	} else {
		o.Status = StatusFailed
	}
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == StatusCompleted || o.Status == StatusFailed {
		return errors.New("cannot cancel completed or failed order")
	}
	o.Status = StatusCancelled
	return nil
}
