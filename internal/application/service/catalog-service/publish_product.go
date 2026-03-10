package catalog_service

import (
	"context"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

// PublishProduct переводит товар из статуса DRAFT в PUBLISHED.
// Публикация возможна только если у товара есть остаток на складе
// и привязана хотя бы одна категория.
func (s *Implementation) PublishProduct(
	ctx context.Context,
	req service.PublishProductRequest,
) (service.PublishProductResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	product, err := s.productRepo.FindByID(ctx, req.ID)

	if err != nil {
		return service.PublishProductResponse{}, err
	}

	err = product.Publish()
	if err != nil {
		return service.PublishProductResponse{}, err
	}

	err = s.productRepo.Update(ctx, product)

	if err != nil {
		return service.PublishProductResponse{}, err
	}

	return service.PublishProductResponse{}, nil
}
