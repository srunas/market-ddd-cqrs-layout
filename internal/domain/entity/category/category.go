package category

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

type Category struct {
	ID        uuid.UUID
	Name      string
	ParentID  *uuid.UUID
	CreatedAt time.Time
	Level     int
	Status    Status
}

func NewCategory(name string, parentID *uuid.UUID) (*Category, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	level := 0
	if parentID != nil {
		level = 1
	}
	return &Category{
		ID:        uuid.New(),
		Name:      name,
		ParentID:  parentID,
		Level:     level,
		CreatedAt: time.Now(),
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
