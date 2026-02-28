// cmd/app/main.go
package main

import (
	"log"
	"time"

	"github.com/srunas/market-ddd-cqrs-layout/internal/domain/entity/user"
)

func main() {
	log.Println("Магазин запущен")

	u, err := user.New("ivan123", "Иванов", "ivan@example.com", user.RoleBuyer)
	if err != nil {
		log.Printf("Ошибка создания пользователя: %v\n", err)
		return
	}

	log.Printf("Создан пользователь:\n")
	log.Printf("  ID:       %s\n", u.ID)
	log.Printf("  Username: %s\n", u.Username)
	log.Printf("  Email:    %s\n", u.Email)
	log.Printf("  Created:  %s\n", u.CreatedAt.Format(time.RFC3339))
}
