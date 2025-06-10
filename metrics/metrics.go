package metrics

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP metrics
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Ticket metrics
	TicketOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ticket_operations_total",
			Help: "Total number of ticket operations",
		},
		[]string{"operation", "status"},
	)

	TicketStatusGauge = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ticket_status_total",
			Help: "Total number of tickets by status",
		},
		[]string{"status"},
	)

	// Error metrics
	ErrorTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_total",
			Help: "Total number of errors",
		},
		[]string{"type"},
	)
)

// PrometheusMiddleware returns a gin middleware for Prometheus metrics
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		timer := prometheus.NewTimer(HttpRequestDuration.WithLabelValues(c.Request.Method, path))
		c.Next()
		timer.ObserveDuration()

		status := c.Writer.Status()
		HttpRequestsTotal.WithLabelValues(c.Request.Method, path, string(status)).Inc()
	}
}
