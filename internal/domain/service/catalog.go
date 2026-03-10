package service

import (
	"context"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type CategoryNode struct {
	ID       string
	Name     string
	Level    int
	Children []*CategoryNode
}

type CreateCategoryRequest struct {
	Name     string
	ParentID *types.CategoryID
}

type CreateCategoryResponse struct {
	ID types.CategoryID
}

type GetCategoryTreeResponse struct {
	Roots []*CategoryNode
}

type Catalog interface {
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (CreateCategoryResponse, error)
	GetCategoryTree(ctx context.Context) (GetCategoryTreeResponse, error)
}
