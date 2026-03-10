package catalog_service

import (
	"context"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/repository"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) ListProducts(
	ctx context.Context,
	req service.ListProductsRequest,
) (service.ListProductsResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	filter := repository.ProductFilter{
		CategoryID: req.CategoryID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Status:     req.Status,
		Limit:      req.PageSize,
		Offset:     (req.Page - 1) * req.PageSize,
	}
	products, err := s.productRepo.GetProductList(ctx, filter)
	if err != nil {
		return service.ListProductsResponse{}, err
	}

	items := make([]*service.ProductItem, 0, len(products))
	for _, p := range products {
		items = append(items, &service.ProductItem{
			ID:          p.ID,
			Name:        p.Name,
			Price:       p.Price,
			Status:      p.Status,
			CategoryIDs: p.CategoryIDs,
		})
	}

	return service.ListProductsResponse{Items: items}, nil
}
