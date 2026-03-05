package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/srunas/market-ddd-cqrs-layout/internal/migrator"
)

const (
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
	idleTimeout    = 120 * time.Second
	requestTimeout = 5 * time.Second
)

// User represents a system user entity.
type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	// #nosec G117 -- This is a JSON model field, not a hardcoded password
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Role      string    `json:"role"`
}

func main() {
	if err := run(); err != nil {
		log.Printf("Shutting down with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	dbURI := os.Getenv("DB_URI")

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return fmt.Errorf("failed to create pool: %w", err)
	}
	defer pool.Close()

	if pErr := pool.Ping(ctx); pErr != nil {
		return fmt.Errorf("failed to ping db: %w", pErr)
	}

	dbSql := stdlib.OpenDB(*pool.Config().ConnConfig)
	m := migrator.NewMigrator(dbSql, os.Getenv("MIGRATIONS_DIR"))

	if mErr := m.Up(); mErr != nil {
		return fmt.Errorf("migration failed: %w", mErr)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/users", createUserHandler(pool))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	log.Println("Server starting on :8080")
	return server.ListenAndServe()
}

// createUserHandler returns a handler for user creation.
func createUserHandler(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Некорректный JSON", http.StatusBadRequest)
			return
		}

		if user.Role == "" {
			user.Role = "BUYER"
		}

		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		query := `INSERT INTO users (username, email, role) VALUES ($1, $2, $3)`
		_, err := db.Exec(ctx, query, user.Username, user.Email, user.Role)
		if err != nil {
			log.Printf("Ошибка вставки в базу: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintln(w, "Пользователь успешно создан")
	}
}
