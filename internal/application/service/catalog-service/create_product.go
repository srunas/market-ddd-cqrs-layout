package catalog_service

import (
	"context"
	"fmt"

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
		req.CategoryIDs,
	)
	if err != nil {
		return service.CreateProductResponse{}, err
	}

	for _, categoryID := range req.CategoryIDs {
		_, err = s.category.FindByID(ctx, categoryID)
		if err != nil {
			return service.CreateProductResponse{},
				fmt.Errorf("категория %v не найдена: %w", categoryID, err)
		}
	}

	err = s.productRepo.Save(ctx, productNew)

	if err != nil {
		return service.CreateProductResponse{}, err
	}

	return service.CreateProductResponse{ID: productNew.ID}, nil
}
