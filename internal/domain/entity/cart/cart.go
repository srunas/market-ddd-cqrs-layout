package cart

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID
	BuyerID   uuid.UUID
	CreatedAt time.Time
	Items     []CartItem
}

type CartItem struct {
	ProductID uuid.UUID
	Quantity  int
}

func NewCart(buyerID uuid.UUID, createdAt time.Time) *Cart {
	return &Cart{
		ID:        uuid.New(),
		BuyerID:   buyerID,
		CreatedAt: createdAt,
		Items:     []CartItem{},
	}
}

func (c *Cart) AddItem(productID uuid.UUID, quantity int) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items[i].Quantity += quantity
			return
		}
	}
	c.Items = append(c.Items, CartItem{ProductID: productID, Quantity: quantity})
}

func (c *Cart) RemoveItem(productID uuid.UUID, quantity int) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}
