package auth

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	ID        uuid.UUID
	Password  string
	Username  string
	AuthAt    time.Time
	CreatedAt time.Time
}

func NewAuth(username, plainPassword string) (*Auth, error) {
	if username == "" || plainPassword == "" {
		return nil, fmt.Errorf("username and password are required")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Auth{
		ID:        uuid.New(),
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}, nil
}

func (a *Auth) ValidatePassword(plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(plainPassword)) == nil
}

func (a *Auth) UpdateAuthTime() {
	a.AuthAt = time.Now()
}
