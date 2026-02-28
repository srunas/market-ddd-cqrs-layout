package product

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type Status string

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyRUB Currency = "RUB"
)

type Product struct {
	ID          types.ProductID
	Name        string
	Description string
	Price       decimal.Decimal
	CreatedAt   time.Time
	Currency    Currency
	Stock       int
	SellerID    uuid.UUID
	UpdatedAt   time.Time
	CategoryIDs []uuid.UUID
	Attributes  map[string]interface{}
	Status      Status
}

func New(name, description string, price decimal.Decimal, currency Currency, sellerID uuid.UUID) (*Product, error) {
	if name == "" || price.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("name or price must be greater than zero")
	}
	return &Product{
		ID:          types.NewProductID(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now().UTC(),
		Currency:    currency,
		Stock:       0,
		SellerID:    sellerID,
		CategoryIDs: []uuid.UUID{},
		Attributes:  make(map[string]interface{}),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func (p *Product) AddCategory(categoryID uuid.UUID) {
	p.CategoryIDs = append(p.CategoryIDs, categoryID)
	p.UpdatedAt = time.Now().UTC()
}

func (p *Product) UpdateStock(delta int) error {
	if p.Stock+delta < 0 {
		return errors.New("stock cannot be negative")
	}
	p.Stock += delta
	p.UpdatedAt = time.Now().UTC()
	return nil
}
