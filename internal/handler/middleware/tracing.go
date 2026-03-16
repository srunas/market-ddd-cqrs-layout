package middleware

import (
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// responseWriter оборачивает http.ResponseWriter чтобы перехватить статус код.
type responseWriter struct {
	http.ResponseWriter

	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// WithRequestLogging логирует каждый входящий запрос с trace_id, методом, путём и статус кодом.
// Должен применяться ПОСЛЕ WithLogger, чтобы логгер уже был в контексте.
func WithRequestLogging() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			l := FromContext(r.Context())
			if l == nil {
				return
			}

			// Извлекаем trace_id из span контекста OpenTelemetry
			span := trace.SpanFromContext(r.Context())
			traceID := span.SpanContext().TraceID().String()

			l.Info("входящий запрос",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.statusCode,
				"duration_ms", time.Since(start).Milliseconds(),
				"trace_id", traceID,
			)
		})
	}
}
