package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/shopspring/decimal"
	"github.com/sqlc-dev/pqtype"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/product"
	domainrepo "github.com/srunas/market-ddd-cqrs-layout/internal/domain/repository"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
	"github.com/srunas/market-ddd-cqrs-layout/internal/infrastructure/repository/sqlcgen"
)

type ProductRepository struct {
	q *sqlcgen.Queries
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{q: sqlcgen.New(stdlib.OpenDBFromPool(pool))}
}

func (r *ProductRepository) Save(ctx context.Context, p *product.Product) error {
	attrBytes, err := json.Marshal(p.Attributes)
	if err != nil {
		return err
	}

	categoryIDs := make([]uuid.UUID, len(p.CategoryIDs))
	for i, id := range p.CategoryIDs {
		categoryIDs[i] = uuid.UUID(id)
	}

	return r.q.SaveProduct(ctx, sqlcgen.SaveProductParams{
		ID:          uuid.UUID(p.ID),
		Name:        p.Name,
		Description: sql.NullString{String: p.Description, Valid: p.Description != ""},
		Price:       p.Price.String(), // decimal → string
		Currency:    string(p.Currency),
		Stock:       int32(p.Stock), //nolint:gosec // stock не может превысить int32
		SellerID:    uuid.UUID(p.SellerID),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Active:      p.Status == product.StatusPublished,
		Attributes:  pqtype.NullRawMessage{RawMessage: attrBytes, Valid: true},
		CategoryIds: categoryIDs,
	})
}

func (r *ProductRepository) Update(ctx context.Context, p *product.Product) error {
	categoryIDs := make([]uuid.UUID, len(p.CategoryIDs))
	for i, id := range p.CategoryIDs {
		categoryIDs[i] = uuid.UUID(id)
	}

	return r.q.UpdateProduct(ctx, sqlcgen.UpdateProductParams{
		ID:          uuid.UUID(p.ID),
		Name:        p.Name,
		Description: sql.NullString{String: p.Description, Valid: p.Description != ""},
		Price:       p.Price.String(),
		Stock:       int32(p.Stock), //nolint:gosec // stock не может превысить int32
		UpdatedAt:   p.UpdatedAt,
		Active:      p.Status == product.StatusPublished,
		CategoryIds: categoryIDs,
	})
}

func (r *ProductRepository) FindByID(ctx context.Context, id types.ProductID) (*product.Product, error) {
	row, err := r.q.FindProductByID(ctx, uuid.UUID(id))
	if err != nil {
		return nil, err
	}
	return toProductDomain(row)
}

func (r *ProductRepository) GetProductList(
	ctx context.Context,
	filter domainrepo.ProductFilter,
) ([]*product.Product, error) {
	var minPrice, maxPrice sql.NullString
	if filter.MinPrice != nil {
		minPrice = sql.NullString{String: filter.MinPrice.String(), Valid: true}
	}
	if filter.MaxPrice != nil {
		maxPrice = sql.NullString{String: filter.MaxPrice.String(), Valid: true}
	}

	var categoryID uuid.NullUUID
	if filter.CategoryID != nil {
		categoryID = uuid.NullUUID{UUID: uuid.UUID(*filter.CategoryID), Valid: true}
	}

	rows, err := r.q.GetProductList(ctx, sqlcgen.GetProductListParams{
		CategoryID: categoryID,
		MinPrice:   minPrice,
		MaxPrice:   maxPrice,
		Limit:      int32(filter.Limit),  //nolint:gosec // limit не может превысить int32
		Offset:     int32(filter.Offset), //nolint:gosec // offset не может превысить int32
	})
	if err != nil {
		return nil, err
	}

	products := make([]*product.Product, 0, len(rows))
	for _, row := range rows {
		p, errProduct := toProductDomain(row)
		if errProduct != nil {
			return nil, errProduct
		}
		products = append(products, p)
	}
	return products, nil
}

func toProductDomain(row sqlcgen.Product) (*product.Product, error) {
	price, err := decimal.NewFromString(row.Price)
	if err != nil {
		return nil, err
	}

	categoryIDs := make([]types.CategoryID, len(row.CategoryIds))
	for i, id := range row.CategoryIds {
		categoryIDs[i] = types.CategoryID(id)
	}

	// JSONB → map[string]interface{}
	attrs := make(map[string]interface{})
	if row.Attributes.Valid {
		_ = json.Unmarshal(row.Attributes.RawMessage, &attrs)
	}

	status := product.StatusDraft
	if row.Active {
		status = product.StatusPublished
	}

	return &product.Product{
		ID:          types.ProductID(row.ID),
		Name:        row.Name,
		Description: row.Description.String,
		Price:       price,
		Currency:    product.Currency(row.Currency),
		Stock:       int(row.Stock),
		SellerID:    types.UserID(row.SellerID),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		CategoryIDs: categoryIDs,
		Attributes:  attrs,
		Status:      status,
	}, nil
}
