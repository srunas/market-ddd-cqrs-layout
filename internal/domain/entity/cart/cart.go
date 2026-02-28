package cart

import (
	"time"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

type Cart struct {
	ID        types.CartID
	BuyerID   types.UserID
	CreatedAt time.Time
	Items     []CartItem
}

type CartItem struct {
	ProductID types.ProductID
	Quantity  int64
}

func New(buyerID types.UserID) *Cart {
	return &Cart{
		ID:        types.NewCartID(),
		BuyerID:   buyerID,
		CreatedAt: time.Now().UTC(),
		Items:     []CartItem{},
	}
}

func (c *Cart) AddItem(productID types.ProductID, quantity int64) {
	if quantity <= 0 {
		return
	}

	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			c.Items[i].Quantity += quantity
			return
		}
	}
	c.Items = append(c.Items, CartItem{
		ProductID: productID,
		Quantity:  quantity,
	})
}

func (c *Cart) RemoveItem(productID types.ProductID) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			return
		}
	}
}

func (c *Cart) DecreaseQuantity(productID types.ProductID, quantity int64) {
	if quantity <= 0 {
		return
	}
	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			c.Items[i].Quantity -= quantity
			if c.Items[i].Quantity <= 0 {
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			}
			return
		}
	}
}
