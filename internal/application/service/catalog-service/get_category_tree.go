package catalog_service

import (
	"context"

	"github.com/google/uuid"
	"github.com/not-for-prod/observer/tracer/prospan"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/service"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

func (s *Implementation) GetCategoryTree(ctx context.Context) (service.GetCategoryTreeResponse, error) {
	ctx, span := prospan.Start(ctx)
	defer span.End()

	allCategories, err := s.category.FindAll(ctx)
	if err != nil {
		return service.GetCategoryTreeResponse{}, err
	}

	nodes := make(map[types.CategoryID]*service.CategoryNode)

	for _, category := range allCategories {
		nodes[category.ID] = &service.CategoryNode{
			ID:       uuid.UUID(category.ID).String(),
			Name:     category.Name,
			Level:    category.Level,
			Children: []*service.CategoryNode{},
		}
	}

	var roots []*service.CategoryNode

	for _, category := range allCategories {
		node := nodes[category.ID]

		if category.ParentID == nil {
			roots = append(roots, node)
		} else {
			parentNode := nodes[*category.ParentID]
			parentNode.Children = append(parentNode.Children, node)
		}
	}

	return service.GetCategoryTreeResponse{Roots: roots}, nil
}
