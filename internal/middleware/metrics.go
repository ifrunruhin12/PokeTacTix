package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpRequestsTotal counts all HTTP requests by method, path, and status code
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration tracks request latency by method and path
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "poketactix_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// httpActiveRequests tracks currently in-flight requests
	httpActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "poketactix_http_active_requests",
			Help: "Number of currently active HTTP requests",
		},
	)

	// BattleStartTotal counts battle start events by mode
	BattleStartTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_battles_started_total",
			Help: "Total number of battles started",
		},
		[]string{"mode"},
	)

	// BattleResultTotal counts battle results
	BattleResultTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_battle_results_total",
			Help: "Total number of battle results",
		},
		[]string{"result"}, // win, loss, draw
	)

	// PokemonFetchTotal tracks where pokemon data came from (attempt-based counter)
	PokemonFetchTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_pokemon_fetch_total",
			Help: "Total number of pokemon fetch attempts by source",
		},
		[]string{"source"}, // redis, postgres, pokeapi
	)

	// AuthTotal tracks auth events
	AuthTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_auth_total",
			Help: "Total authentication events",
		},
		[]string{"event"}, // login_success, login_failure, register
	)

	// RegisteredUsers tracks registered users (set on startup)
	RegisteredUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "poketactix_registered_users",
			Help: "Total registered users",
		},
	)
)

// PrometheusMiddleware records HTTP metrics for every request
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip OPTIONS preflight requests — they're CORS noise, not real traffic
		if c.Method() == fiber.MethodOptions {
			return c.Next()
		}

		start := time.Now()
		method := c.Method() // capture before c.Next() to avoid race
		httpActiveRequests.Inc()
		defer httpActiveRequests.Dec()

		err := c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())

		path := c.Route().Path
		if path == "" {
			path = "unknown"
		}

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}
