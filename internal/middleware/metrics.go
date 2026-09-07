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

	// battleStartTotal counts battle start events by mode
	BattleStartTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_battles_started_total",
			Help: "Total number of battles started",
		},
		[]string{"mode"},
	)

	// battleResultTotal counts battle results
	BattleResultTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_battle_results_total",
			Help: "Total number of battle results",
		},
		[]string{"result"}, // win, loss, draw
	)

	// pokemonFetchTotal tracks where pokemon data came from
	PokemonFetchTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_pokemon_fetch_total",
			Help: "Total number of pokemon fetches by source",
		},
		[]string{"source"}, // redis, postgres, pokeapi
	)

	// authTotal tracks auth events
	AuthTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_auth_total",
			Help: "Total authentication events",
		},
		[]string{"event"}, // login_success, login_failure, register
	)

	// activeUsers tracks registered users (set on startup)
	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "poketactix_active_users",
			Help: "Total registered users",
		},
	)
)

// PrometheusMiddleware records HTTP metrics for every request
func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		httpActiveRequests.Inc()

		err := c.Next()

		httpActiveRequests.Dec()
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Response().StatusCode())

		// Normalize path to avoid high cardinality from IDs in URLs
		path := normalizePath(c.Route().Path)

		httpRequestsTotal.WithLabelValues(c.Method(), path, status).Inc()
		httpRequestDuration.WithLabelValues(c.Method(), path).Observe(duration)

		return err
	}
}

// normalizePath replaces dynamic segments to reduce cardinality
// e.g. /api/cards/123 -> /api/cards/:id
func normalizePath(path string) string {
	if path == "" {
		return "unknown"
	}
	return path
}
