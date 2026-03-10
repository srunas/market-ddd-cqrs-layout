package catalog_service

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/repository"
)

type Implementation struct {
	category    repository.Category
	productRepo repository.Product
	txManager   trm.Manager
}

func NewImplementation(
	category repository.Category,
	productRepo repository.Product,
	txManager trm.Manager,
) *Implementation {
	return &Implementation{
		category:    category,
		productRepo: productRepo,
		txManager:   txManager,
	}
}
