package catalog_service

import (
	"context"
	"fmt"

	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/category"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
)

func (s *Implementation) CreateCategory(
	ctx context.Context,
	req service.CreateCategoryRequest,
) (service.CreateCategoryResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	var parent *category.Category
	var err error

	if req.ParentID != nil {
		parent, err = s.category.FindByID(ctx, *req.ParentID)
		if err != nil {
			return service.CreateCategoryResponse{}, err
		}
	}

	newCategory, err := category.New(req.Name, parent)
	if err != nil {
		return service.CreateCategoryResponse{}, fmt.Errorf("ошибка создания категорий: %w", err)
	}

	err = s.txManager.Do(ctx, func(ctx context.Context) error {
		return s.category.Save(ctx, newCategory)
	})

	if err != nil {
		return service.CreateCategoryResponse{}, err
	}

	return service.CreateCategoryResponse{
		ID: newCategory.ID,
	}, nil
}
