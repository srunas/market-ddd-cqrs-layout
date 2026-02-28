// Package user содержит доменную модель пользователя и связанные бизнес-правила.
package user

import (
	"errors"
	"time"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
)

// Role UserRole определяет роль пользователя в системе.
type Role string

const (

	// RoleBuyer — роль обычного покупателя.
	RoleBuyer Role = "BUYER"
)

// User представляет зарегистрированного пользователя системы.
type User struct {
	ID        types.UserID
	Username  string
	Surname   string
	Role      Role
	AuthAt    time.Time // Время последней успешной аутентификации
	CreatedAt time.Time
	Email     string
	Enabled   bool
}

// New создаёт нового пользователя с минимальной валидацией.
func New(username, surname, email string, role Role) (*User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}

	return &User{
		ID:        types.NewUserID(),
		Username:  username,
		Surname:   surname,
		Role:      role,
		Email:     email,
		Enabled:   true,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// UpdateRole меняет роль пользователя.
func (u *User) UpdateRole(newRole Role) {
	u.Role = newRole
	// В будущем здесь может быть логика: проверка допустимых переходов, событие RoleChanged и т.д.
}

// UpdateUsername обновляет имя пользователя.
func (u *User) UpdateUsername(newUsername string) {
	if newUsername != "" {
		u.Username = newUsername
	}
}

// UpdateSurname обновляет фамилию пользователя.
func (u *User) UpdateSurname(newSurname string) {
	u.Surname = newSurname // можно добавить валидацию длины и т.д.
}

// UpdateEmail обновляет email пользователя.
func (u *User) UpdateEmail(newEmail string) {
	if newEmail != "" {
		u.Email = newEmail
	}
}
