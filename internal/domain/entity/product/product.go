package product

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type Status string

type Currency string

const (
	StatusDraft     Status = "DRAFT"
	StatusPublished Status = "PUBLISHED"
	StatusArchived  Status = "ARCHIVED"
)

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
	SellerID    types.UserID
	UpdatedAt   time.Time
	CategoryIDs []types.CategoryID
	Attributes  map[string]interface{}
	Status      Status
}

func New(name, description string, price decimal.Decimal, currency Currency, sellerID types.UserID) (*Product, error) {
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
		CategoryIDs: []types.CategoryID{},
		Attributes:  make(map[string]interface{}),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func (p *Product) Publish() error {
	if p.Status != StatusDraft {
		return errors.New("cannot publish product with status not Draft")
	}

	if p.Stock <= 0 {
		return errors.New("cannot publish product with zero stock")
	}

	if len(p.CategoryIDs) == 0 {
		return errors.New("cannot publish product with zero category IDs")
	}

	p.Status = StatusPublished
	p.UpdatedAt = time.Now().UTC()

	return nil
}

func (p *Product) AddCategory(categoryID types.CategoryID) error {
	p.CategoryIDs = append(p.CategoryIDs, categoryID)
	p.UpdatedAt = time.Now().UTC()
	return nil
}

func (p *Product) UpdateStock(delta int) error {
	if p.Stock+delta < 0 {
		return errors.New("stock cannot be negative")
	}
	p.Stock += delta
	p.UpdatedAt = time.Now().UTC()
	return nil
}
