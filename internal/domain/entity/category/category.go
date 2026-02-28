package category

import (
	"errors"
	"time"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Category struct {
	ID        types.CategoryID
	Name      string
	ParentID  *types.CategoryID
	CreatedAt time.Time
	Level     int
	Status    Status
}

func New(name string, parentID *types.CategoryID) (*Category, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	level := 0
	if parentID != nil {
		level = 1
	}
	return &Category{
		ID:        types.NewCategoryID(),
		Name:      name,
		ParentID:  parentID,
		Level:     level,
		CreatedAt: time.Now().UTC(),
		Status:    StatusActive,
	}, nil
}

func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

func (c *Category) IsActive() bool {
	return c.Status == StatusActive
}

func (c *Category) IsInactive() bool {
	return c.Status == StatusInactive
}
