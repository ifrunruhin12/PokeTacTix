package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Registry is the dedicated Prometheus registry for PokeTacTix metrics.
// Using a custom registry (instead of the default global one) prevents
// duplicate-collection errors when promauto vars are initialized.
var Registry = prometheus.NewRegistry()

// factory wraps Registry so we can define metrics in var blocks cleanly.
var factory = prometheus.WrapRegistererWith(prometheus.Labels{}, Registry)

func newCounterVec(opts prometheus.CounterOpts, labels []string) *prometheus.CounterVec {
	c := prometheus.NewCounterVec(opts, labels)
	Registry.MustRegister(c)
	return c
}

func newHistogramVec(opts prometheus.HistogramOpts, labels []string) *prometheus.HistogramVec {
	h := prometheus.NewHistogramVec(opts, labels)
	Registry.MustRegister(h)
	return h
}

func newGauge(opts prometheus.GaugeOpts) prometheus.Gauge {
	g := prometheus.NewGauge(opts)
	Registry.MustRegister(g)
	return g
}

func init() {
	// Include standard Go runtime and process metrics in our custom registry
	Registry.MustRegister(collectors.NewGoCollector())
	Registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}

// Ensure factory is used to suppress unused import warning
var _ = factory

var (
	// httpRequestsTotal counts all HTTP requests by method, path, and status code
	httpRequestsTotal = newCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// httpRequestDuration tracks request latency by method and path
	httpRequestDuration = newHistogramVec(
		prometheus.HistogramOpts{
			Name:    "poketactix_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// httpActiveRequests tracks currently in-flight requests
	httpActiveRequests = newGauge(
		prometheus.GaugeOpts{
			Name: "poketactix_http_active_requests",
			Help: "Number of currently active HTTP requests",
		},
	)

	// BattleStartTotal counts battle start events by mode
	BattleStartTotal = newCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_battles_started_total",
			Help: "Total number of battles started",
		},
		[]string{"mode"},
	)

	// BattleResultTotal counts battle results
	BattleResultTotal = newCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_battle_results_total",
			Help: "Total number of battle results",
		},
		[]string{"result"}, // win, loss, draw
	)

	// PokemonFetchTotal tracks where pokemon data came from (attempt-based counter)
	PokemonFetchTotal = newCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_pokemon_fetch_total",
			Help: "Total number of pokemon fetch attempts by source",
		},
		[]string{"source"}, // redis, postgres, pokeapi
	)

	// AuthTotal tracks auth events
	AuthTotal = newCounterVec(
		prometheus.CounterOpts{
			Name: "poketactix_auth_total",
			Help: "Total authentication events",
		},
		[]string{"event"}, // login_success, login_failure, register
	)

	// RegisteredUsers tracks total registered users, refreshed periodically from DB
	RegisteredUsers = newGauge(
		prometheus.GaugeOpts{
			Name: "poketactix_registered_users",
			Help: "Total registered users in the database",
		},
	)
)

// MetricsHandler returns an http.Handler that serves metrics from our custom registry.
func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(Registry, promhttp.HandlerOpts{
		EnableOpenMetrics: false,
	})
}

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

		// Use the registered route template (e.g. /api/cards/:id).
		// Fall back to the raw request path if route is unmatched.
		path := c.Route().Path
		if path == "" {
			path = c.Path()
		}
		if path == "" {
			path = "unknown"
		}

		httpRequestsTotal.WithLabelValues(method, path, status).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}
