package auth

import (
	"errors"
	"time"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/types"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	ID        types.AuthID
	password  string
	username  string
	AuthAt    time.Time
	CreatedAt time.Time
}

func New(username, plainPassword string) (*Auth, error) {
	if username == "" || plainPassword == "" {
		return nil, errors.New("username and password are required")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Auth{
		ID:        types.NewAuth(),
		username:  username,
		password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
	}, nil
}

func NewFromDB(id types.AuthID, username, password string, authAt, createdAt time.Time) *Auth {
	return &Auth{
		ID:        id,
		password:  password,
		username:  username,
		AuthAt:    authAt,
		CreatedAt: createdAt,
	}
}

func (a *Auth) ValidatePassword(plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.password), []byte(plainPassword)) == nil
}

func (a *Auth) UpdateAuthTime() {
	a.AuthAt = time.Now().UTC()
}

func (a *Auth) Password() string { return a.password }
func (a *Auth) Username() string { return a.username }
