package cart_service

import (
	"context"
	"fmt"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/cart"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) AddItem(ctx context.Context,
	req service.AddItemRequest) (service.AddItemResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	_, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return service.AddItemResponse{}, fmt.Errorf("товар не найден: %w", err)
	}

	cartEntity, err := s.cartRepo.FindByBuyerID(ctx, req.BuyerID)
	if err != nil {
		newCart := cart.New(req.BuyerID)
		newCart.AddItem(req.ProductID, req.Quantity)
		err = s.cartRepo.Save(ctx, newCart)
		if err != nil {
			return service.AddItemResponse{}, err
		}
		return service.AddItemResponse{}, nil
	}

	cartEntity.AddItem(req.ProductID, req.Quantity)
	err = s.cartRepo.Update(ctx, cartEntity)
	if err != nil {
		return service.AddItemResponse{}, err
	}
	return service.AddItemResponse{}, nil
}
