package cart_service

import (
	"context"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) GetCart(ctx context.Context, req service.GetCartRequest) (
	service.GetCartResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	cartEntity, err := s.cartRepo.FindByBuyerID(ctx, req.BuyerID)
	if err != nil {
		return service.GetCartResponse{}, err
	}

	items := make([]service.CartItem, 0, len(cartEntity.Items))
	for _, item := range cartEntity.Items {
		items = append(items, service.CartItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return service.GetCartResponse{
		ID:    cartEntity.ID,
		Items: items,
	}, nil
}
