package identity_service

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/repository"
)

type Implementation struct {
	userRepo  repository.User
	authRepo  repository.Auth
	txManager trm.Manager
	jwtSecret []byte
}

func NewImplementation(
	userRepo repository.User,
	authRepo repository.Auth,
	txManager trm.Manager,
	jwtSecret []byte,
) *Implementation {
	return &Implementation{
		userRepo:  userRepo,
		authRepo:  authRepo,
		txManager: txManager,
		jwtSecret: jwtSecret,
	}
}
