package middleware

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metricsMiddleware struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func newMetricsMiddleware() *metricsMiddleware {
	return &metricsMiddleware{
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Общее количество HTTP запросов",
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Длительность HTTP запросов в секундах",
				Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
			},
			[]string{"method", "path"},
		),
	}
}

func WithMetrics() func(http.Handler) http.Handler {
	m := newMetricsMiddleware()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/metrics" {
				next.ServeHTTP(w, r)
				return
			}

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			timer := prometheus.NewTimer(m.requestDuration.WithLabelValues(r.Method, r.URL.Path))
			next.ServeHTTP(rw, r)
			timer.ObserveDuration()

			m.requestsTotal.WithLabelValues(
				r.Method,
				r.URL.Path,
				strconv.Itoa(rw.statusCode),
			).Inc()
		})
	}
}
