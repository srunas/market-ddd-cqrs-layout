package catalog_service

import (
	"context"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/product"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) CreateProduct(
	ctx context.Context,
	req service.CreateProductRequest,
) (service.CreateProductResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	productNew, err := product.New(
		req.Name,
		req.Description,
		req.Price,
		req.Currency,
		req.SellerID,
	)
	if err != nil {
		return service.CreateProductResponse{}, err
	}

	for _, catID := range req.CategoryIDs {
		err = productNew.AddCategory(catID)
		if err != nil {
			return service.CreateProductResponse{}, err
		}
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		return s.productRepo.Save(ctx, productNew)
	})

	if err != nil {
		return service.CreateProductResponse{}, err
	}

	return service.CreateProductResponse{ID: productNew.ID}, nil
}
