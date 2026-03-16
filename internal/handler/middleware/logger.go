package middleware

import (
	"context"
	"net/http"

	"github.com/not-for-prod/observer/logger"
)

type contextKey string

const loggerKey contextKey = "logger"

func WithLogger(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loggerKey, l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func FromContext(ctx context.Context) logger.Logger {
	l, _ := ctx.Value(loggerKey).(logger.Logger)
	return l
}
