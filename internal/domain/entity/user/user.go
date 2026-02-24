package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "ADMIN"
	UserRoleBuyer UserRole = "BUYER"
)

type User struct {
	ID        uuid.UUID
	Username  string
	Surname   string
	Role      UserRole
	AuthAt    time.Time
	CreatedAt time.Time
	Email     string
	Enabled   bool
}

func NewUser(username, surname, email string, role UserRole) (*User, error) {
	if username == "" || email == "" {
		return nil, fmt.Errorf("Username or email is empty")
	}
	return &User{
		ID:        uuid.New(),
		Username:  username,
		Surname:   surname,
		Role:      role,
		Email:     email,
		Enabled:   true,
		CreatedAt: time.Now(),
	}, nil
}

func (u *User) UpdateRole(newRole UserRole) {
	u.Role = newRole
}

func (u *User) UpdateUsername(newUsername string) {
	u.Username = newUsername
}

func (u *User) UpdateSurname(newSurname string) {
	u.Surname = newSurname
}

func (u *User) UpdateEmail(newEmail string) {
	u.Email = newEmail
}
