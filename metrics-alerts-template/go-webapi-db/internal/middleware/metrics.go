package middleware

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "uri", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_server_requests_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "uri", "status"},
	)

	httpServerErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_errors_total",
			Help: "Total number of HTTP server errors (5xx)",
		},
		[]string{"method", "uri", "status"},
	)

	httpClientErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_server_client_errors_total",
			Help: "Total number of HTTP client errors (4xx)",
		},
		[]string{"method", "uri", "status"},
	)
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		uri := sanitizeURI(path)

		// Record metrics
		httpRequestsTotal.WithLabelValues(method, uri, status).Inc()
		httpRequestDuration.WithLabelValues(method, uri, status).Observe(duration)

		// Record error metrics
		statusCode := c.Writer.Status()
		if statusCode >= 500 {
			httpServerErrors.WithLabelValues(method, uri, status).Inc()
		} else if statusCode >= 400 {
			httpClientErrors.WithLabelValues(method, uri, status).Inc()
		}
	}
}

func sanitizeURI(uri string) string {
	// Replace path variables with placeholders for better metric aggregation
	if strings.HasPrefix(uri, "/api/users/") {
		parts := strings.Split(strings.TrimPrefix(uri, "/api/users/"), "/")
		if len(parts) > 0 {
			// Check if it's a numeric ID
			if _, err := strconv.Atoi(parts[0]); err == nil {
				return "/api/users/{id}"
			}
			// Check for specific patterns
			if parts[0] == "email" && len(parts) > 1 {
				return "/api/users/email/{email}"
			}
			if parts[0] == "status" && len(parts) > 1 {
				if len(parts) > 2 && parts[2] == "count" {
					return "/api/users/status/{status}/count"
				}
				return "/api/users/status/{status}"
			}
			if parts[0] == "external" && len(parts) > 1 {
				return "/api/users/external/{serviceName}"
			}
			if parts[0] == "test" && len(parts) > 1 {
				return "/api/users/test/" + parts[1]
			}
		}
	}
	return uri
}

