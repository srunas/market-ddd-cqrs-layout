package types

import (
	"github.com/google/uuid"
)

type (
	AuthID     uuid.UUID
	UserID     uuid.UUID
	CartID     uuid.UUID
	OrderID    uuid.UUID
	ProductID  uuid.UUID
	CategoryID uuid.UUID
)

func NewAuth() AuthID {
	return AuthID(uuid.New())
}

func NewUserID() UserID {
	return UserID(uuid.New())
}

func NewCartID() CartID {
	return CartID(uuid.New())
}

func NewOrderID() OrderID {
	return OrderID(uuid.New())
}

func NewProductID() ProductID {
	return ProductID(uuid.New())
}

func NewCategoryID() CategoryID {
	return CategoryID(uuid.New())
}
