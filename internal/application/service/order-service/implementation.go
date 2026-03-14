package order_service

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/repository"
)

type Implementation struct {
	orderRepo   repository.Order
	cartRepo    repository.Cart
	productRepo repository.Product
	txManager   trm.Manager
}

func NewImplementation(
	orderRepo repository.Order,
	cartRepo repository.Cart,
	productRepo repository.Product,
	txManager trm.Manager,
) *Implementation {
	return &Implementation{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		txManager:   txManager,
	}
}
