package product

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyRUB Currency = "RUB"
)

type Product struct {
	ID          uuid.UUID
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

func NewProduct(name, description string, price decimal.Decimal, currency Currency, sellerID uuid.UUID) (*Product, error) {
	if name == "" || price.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("name or price must be greater than zero")
	}
	return &Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
		Currency:    currency,
		Stock:       0,
		SellerID:    sellerID,
		CategoryIDs: []uuid.UUID{},
		Attributes:  make(map[string]interface{}),
		UpdatedAt:   time.Now(),
	}, nil
}

func (p *Product) AddCategory(categoryID uuid.UUID) {
	p.CategoryIDs = append(p.CategoryIDs, categoryID)
	p.UpdatedAt = time.Now()
}

func (p *Product) UpdateStock(delta int) error {
	if p.Stock+delta < 0 {
		return fmt.Errorf("stock cannot be negative")
	}
	p.Stock += delta
	p.UpdatedAt = time.Now()
	return nil
}
