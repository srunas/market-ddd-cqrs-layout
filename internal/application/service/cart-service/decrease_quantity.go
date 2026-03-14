package cart_service

import (
	"context"
	"fmt"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) DecreaseQuantity(ctx context.Context, req service.DecreaseQuantityRequest) (
	service.DecreaseQuantityResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	_, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return service.DecreaseQuantityResponse{}, fmt.Errorf("товар не найден: %w", err)
	}

	cartEntity, err := s.cartRepo.FindByBuyerID(ctx, req.BuyerID)
	if err != nil {
		return service.DecreaseQuantityResponse{}, fmt.Errorf("корзина не найдена: %w", err)
	}

	cartEntity.DecreaseQuantity(req.ProductID, req.Quantity)
	err = s.cartRepo.Update(ctx, cartEntity)
	if err != nil {
		return service.DecreaseQuantityResponse{}, err
	}

	return service.DecreaseQuantityResponse{}, nil
}
