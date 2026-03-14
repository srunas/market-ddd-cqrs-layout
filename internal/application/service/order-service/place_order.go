package order_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/order"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) PlaceOrder(ctx context.Context, req service.PlaceOrderRequest) (
	service.PlaceOrderResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	cartEntity, err := s.cartRepo.FindByBuyerID(ctx, req.BuyerID)
	if err != nil {
		return service.PlaceOrderResponse{}, fmt.Errorf("корзина не найдена: %w", err)
	}

	if len(cartEntity.Items) == 0 {
		return service.PlaceOrderResponse{}, errors.New("корзина пуста: %w")
	}

	orderEntity := order.New(req.BuyerID, req.Currency, req.PaymentMethod)

	for _, cartItem := range cartEntity.Items {
		productEntity, productErr := s.productRepo.FindByID(ctx, cartItem.ProductID)
		if productErr != nil {
			return service.PlaceOrderResponse{}, fmt.Errorf(
				"товар %v не найден: %w", cartItem.ProductID, productErr)
		}

		err = orderEntity.AddItem(cartItem.ProductID, cartItem.Quantity, productEntity.Price)
		if err != nil {
			return service.PlaceOrderResponse{}, fmt.Errorf("ошибка дабавления товара в заказ: %w", err)
		}
	}

	err = s.orderRepo.Save(ctx, orderEntity)
	if err != nil {
		return service.PlaceOrderResponse{}, err
	}

	return service.PlaceOrderResponse{Order: orderEntity}, nil
}
