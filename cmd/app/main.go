package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	trmsqldb "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	zaploger "github.com/not-for-prod/observer/logger/zap"
	"github.com/not-for-prod/observer/tracer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cart_service "github.com/srunas/market-ddd-cqrs-layout/internal/application/service/cart-service"
	catalog_service "github.com/srunas/market-ddd-cqrs-layout/internal/application/service/catalog-service"
	identity_service "github.com/srunas/market-ddd-cqrs-layout/internal/application/service/identity-service"
	order_service "github.com/srunas/market-ddd-cqrs-layout/internal/application/service/order-service"
	"github.com/srunas/market-ddd-cqrs-layout/internal/handler"
	"github.com/srunas/market-ddd-cqrs-layout/internal/handler/middleware"
	infra "github.com/srunas/market-ddd-cqrs-layout/internal/infrastructure/repository"
	"github.com/srunas/market-ddd-cqrs-layout/internal/migrator"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	readTimeout     = 10 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 120 * time.Second
	shutdownTimeout = 5 * time.Second
)

func main() {
	if err := run(); err != nil {
		// fmt.Printf запрещён линтером — пишем напрямую в stderr
		_, _ = fmt.Fprintf(os.Stderr, "завершение с ошибкой: %v\n", err)
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

	// — logger в виде JSON —
	logger := zaploger.NewLogger()
	defer func() {
		// Sync сбрасывает буфер — ошибку игнорируем намеренно (возникает на некоторых ОС)
		_ = logger.Sync()
	}()

	// — Трейсер —
	tracerProvider := tracer.NewProvider(
		tracer.WithServiceName("marketplace"),
		tracer.WithHost(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
	)
	if err = tracerProvider.Start(ctx); err != nil {
		return fmt.Errorf("ошибка запуска трейсера: %w", err)
	}
	defer func() {
		if stopErr := tracerProvider.Stop(ctx); stopErr != nil {
			logger.Error("ошибка остановки трейсера", "error", stopErr)
		}
	}()

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

	mux := buildRouter(pool, txManager)

	h := otelhttp.NewHandler(
		middleware.WithLogger(logger)(middleware.WithRequestLogging()(middleware.WithMetrics()(mux))),
		"http",
	)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      h,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if listenErr := server.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			logger.Error("ошибка сервера", "error", listenErr)
		}
	}()

	logger.Info("сервер запущен на :8080")

	<-quit
	logger.Info("Получен сигнал остановки")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownCancel()

	if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
		return fmt.Errorf("ошибка graceful shutdown: %w", shutdownErr)
	}

	logger.Info("сервер остановлен")
	return nil
}

func buildRouter(pool *pgxpool.Pool, txManager *manager.Manager) *http.ServeMux {
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

	_ = catalogSvc
	_ = cartSvc
	_ = orderSvc

	// — HTTP маршруты —
	mux := http.NewServeMux()

	identityHandler := handler.NewIdentityHandler(identitySvc)
	mux.HandleFunc("POST /api/v1/register", identityHandler.Register)
	mux.HandleFunc("POST /api/v1/login", identityHandler.Login)

	mux.HandleFunc("/healthz/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/healthz/ready", func(w http.ResponseWriter, r *http.Request) {
		if pingErr := pool.Ping(r.Context()); pingErr != nil {
			http.Error(w, "db unavailable", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/metrics", promhttp.Handler())

	return mux
}
