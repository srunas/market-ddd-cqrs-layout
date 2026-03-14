package cart_service

import (
	"context"
	"fmt"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) RemoveItem(ctx context.Context, req service.RemoveItemRequest) (
	service.RemoveItemResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	_, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return service.RemoveItemResponse{}, fmt.Errorf("товар не найден: %w", err)
	}

	cartEntity, err := s.cartRepo.FindByBuyerID(ctx, req.BuyerID)
	if err != nil {
		return service.RemoveItemResponse{}, fmt.Errorf("корзина не найдена: %w", err)
	}

	cartEntity.RemoveItem(req.ProductID)
	err = s.cartRepo.Update(ctx, cartEntity)
	if err != nil {
		return service.RemoveItemResponse{}, err
	}

	return service.RemoveItemResponse{}, nil
}
