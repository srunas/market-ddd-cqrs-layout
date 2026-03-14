package order_service

import (
	"context"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) ListOrders(ctx context.Context, req service.ListOrdersRequest) (
	service.ListOrdersResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	orders, err := s.orderRepo.FindByBuyerID(ctx, req.BuyerID)
	if err != nil {
		return service.ListOrdersResponse{}, err
	}

	summaries := make([]*service.OrderSummary, 0, len(orders))
	for _, order := range orders {
		summaries = append(summaries, &service.OrderSummary{
			ID:     order.ID,
			Status: order.Status,
			Total:  order.Total,
		})
	}

	return service.ListOrdersResponse{OrderSummary: summaries}, nil
}
