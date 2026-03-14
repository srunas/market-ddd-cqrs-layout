package cart_service

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/repository"
)

type Implementation struct {
	cartRepo    repository.Cart
	productRepo repository.Product
	txManager   trm.Manager
}

func NewImplementation(
	cartRepo repository.Cart,
	productRepo repository.Product,
	txManager trm.Manager,
) *Implementation {
	return &Implementation{
		cartRepo:    cartRepo,
		productRepo: productRepo,
		txManager:   txManager,
	}
}
