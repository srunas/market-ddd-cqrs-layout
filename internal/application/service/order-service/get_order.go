package order_service

import (
	"context"
	"fmt"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) GetOrder(ctx context.Context, req service.GetOrderRequest) (
	service.GetOrderResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	orderEntity, err := s.orderRepo.FindByID(ctx, req.OrderID)
	if err != nil {
		return service.GetOrderResponse{}, fmt.Errorf("заказ не найден: %w", err)
	}

	return service.GetOrderResponse{Order: orderEntity}, nil
}
