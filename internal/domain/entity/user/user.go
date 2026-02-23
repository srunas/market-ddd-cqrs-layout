package user

import (
	"github.com/google/uuid"
	"time"
)


type UserRole string
const (
	UserRoleAdmin UserRole = "ADMIN"
	UserRoleBuyer UserRole = "BUYER"
)

type User struct {
	ID uuid.UUID
	Username string
	Surname string
	Role UserRole
	AuthAt time.Time
	CreatedAt time.Time
	Email   string
	Enabled bool
}


func NewUser(username, surname, email string, role UserRole) (*User, error) {
	if username == ""
}