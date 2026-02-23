package user

import (
	"net/smtp"
	"time"

	"github.com/google/uuid"
)

type ID string

func (id ID) String() string {
	return string(id)
}

type User struct {
	ID       ID
	Username string
	Surname  string
	//role		UserRole (enum)
	//Auth        Auth
	createdAt time.Time
	Email     string
}

func New(clientID client.ID, title, description string) *Review {
	return &Review{
		ID:          ID(uuid.NewString()),
		ClientID:    clientID,
		Title:       title,
		Description: description,
	}
}
