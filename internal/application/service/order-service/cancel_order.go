package order_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) CancelOrder(ctx context.Context, req service.CancelOrderRequest) (
	service.CancelOrderResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	orderEntity, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return service.CancelOrderResponse{}, fmt.Errorf("заказ не найден: %w", err)
	}

	if orderEntity.BuyerID != req.BuyerID {
		return service.CancelOrderResponse{}, errors.New("доступ запрещён")
	}

	if err = orderEntity.Cancel(); err != nil {
		return service.CancelOrderResponse{}, err
	}

	if err = s.orderRepo.Update(ctx, orderEntity); err != nil {
		return service.CancelOrderResponse{}, err
	}

	return service.CancelOrderResponse{}, nil
}
