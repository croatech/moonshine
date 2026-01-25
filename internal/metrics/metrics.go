package metrics

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "websocket_connections_active",
			Help: "Number of active WebSocket connections",
		},
	)

	FightsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "moonshine_fights_total",
			Help: "Total number of fights started",
		},
	)

	FightDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "moonshine_fight_duration_seconds",
			Help:    "Fight duration in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
		},
	)

	PlayersOnline = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "moonshine_players_online",
			Help: "Number of players currently online",
		},
	)
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start).Seconds()
			status := c.Response().Status
			path := c.Path()
			method := c.Request().Method

			HttpRequestsTotal.WithLabelValues(method, path, string(rune(status))).Inc()
			HttpRequestDuration.WithLabelValues(method, path).Observe(duration)

			return err
		}
	}
}
