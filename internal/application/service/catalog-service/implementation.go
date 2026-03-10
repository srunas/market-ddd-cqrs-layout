package catalog_service

import (
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/repository"
)

type Implementation struct {
	category    repository.Category
	productRepo repository.Product
	txManager   *manager.Manager
}

func NewImplementation(
	category repository.Category,
	productRepo repository.Product,
	txManager *manager.Manager,
) *Implementation {
	return &Implementation{
		category:    category,
		productRepo: productRepo,
		txManager:   txManager,
	}
}
