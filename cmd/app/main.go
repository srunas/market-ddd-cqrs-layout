package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	trmsqldb "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	cart_service "github.com/srunas/market-ddd-cqrs-layout/internal/application/service/cart-service"
	catalog_service "github.com/srunas/market-ddd-cqrs-layout/internal/application/service/catalog-service"
	identity_service "github.com/srunas/market-ddd-cqrs-layout/internal/application/service/identity-service"
	order_service "github.com/srunas/market-ddd-cqrs-layout/internal/application/service/order-service"
	infra "github.com/srunas/market-ddd-cqrs-layout/internal/infrastructure/repository"
	"github.com/srunas/market-ddd-cqrs-layout/internal/migrator"
)

const (
	readTimeout  = 10 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 120 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Printf("завершение с ошибкой: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// — База данных —
	pool, err := pgxpool.New(ctx, os.Getenv("DB_URI"))
	if err != nil {
		return fmt.Errorf("ошибка создания пула: %w", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// — Миграции —
	dbSQL := stdlib.OpenDB(*pool.Config().ConnConfig)
	m := migrator.NewMigrator(dbSQL, os.Getenv("MIGRATIONS_DIR"))
	if err = m.Up(); err != nil {
		return fmt.Errorf("ошибка миграции: %w", err)
	}

	// — Transaction manager —
	txManager, err := manager.New(trmsqldb.NewDefaultFactory(dbSQL))
	if err != nil {
		return fmt.Errorf("ошибка создания tx manager: %w", err)
	}

	// — Репозитории —
	userRepo := infra.NewUserRepository(pool)
	authRepo := infra.NewAuthRepository(pool)
	productRepo := infra.NewProductRepository(pool)
	categoryRepo := infra.NewCategoryRepository(pool)
	cartRepo := infra.NewCartRepository(pool)
	orderRepo := infra.NewOrderRepository(pool)

	// — Сервисы —
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	identitySvc := identity_service.NewImplementation(userRepo, authRepo, txManager, jwtSecret)
	catalogSvc := catalog_service.NewImplementation(categoryRepo, productRepo, txManager)
	cartSvc := cart_service.NewImplementation(cartRepo, productRepo, txManager)
	orderSvc := order_service.NewImplementation(orderRepo, cartRepo, productRepo, txManager)

	_ = identitySvc
	_ = catalogSvc
	_ = cartSvc
	_ = orderSvc

	// — HTTP сервер —
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	log.Println("сервер запущен на :8080")
	return server.ListenAndServe()
}
