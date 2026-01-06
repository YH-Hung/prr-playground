package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go-webapi-db/internal/config"
	"go-webapi-db/internal/health"
	"go-webapi-db/internal/metrics"
	"go-webapi-db/internal/middleware"
	"go-webapi-db/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// TestMetricsEndpoint_AllMetricsExported verifies that all expected metrics are exported via /metrics endpoint
func TestMetricsEndpoint_AllMetricsExported(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Reset Prometheus registry
	reg := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = reg
	prometheus.DefaultGatherer = reg

	cfg := config.Load()

	// Setup router similar to main.go
	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.MetricsMiddleware())

	// Create services
	metricsService := service.NewMetricsService()

	// Create mock handlers - we don't need full functionality for metrics testing
	// Just need to ensure routes exist to generate HTTP metrics
	router.GET("/api/users/:id", func(c *gin.Context) {
		c.JSON(200, gin.H{"id": c.Param("id")})
	})

	// Create mock health handler (nil db is OK for testing - it will just report DOWN)
	healthHandler := health.NewHealthHandler(nil)

	// Register routes
	router.GET("/health", healthHandler.HealthCheck)
	router.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))

	// Make some requests to generate metrics
	req1, _ := http.NewRequest("GET", "/api/users/123", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	req2, _ := http.NewRequest("GET", "/health", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// Exercise metrics service
	metricsService.IncrementUserCreated()
	metricsService.IncrementUserUpdated()
	metricsService.IncrementUserDeleted()
	metricsService.IncrementUserOperationErrors("test_error")
	metricsService.IncrementExternalCallErrors("test_service")

	// Request metrics endpoint
	req, _ := http.NewRequest("GET", cfg.Metrics.Path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()

	// Expected HTTP metrics
	expectedHTTPMetrics := []string{
		"http_server_requests_total",
		"http_server_requests_seconds",
	}

	// Expected custom business metrics
	expectedCustomMetrics := []string{
		"custom_user_created_total",
		"custom_user_updated_total",
		"custom_user_deleted_total",
		"custom_user_operation_errors_total",
		"custom_external_call_errors_total",
	}

	// Expected system metrics (Go runtime)
	expectedSystemMetrics := []string{
		"go_goroutines",
		"go_memstats",
	}

	// Check HTTP metrics
	for _, metric := range expectedHTTPMetrics {
		if !strings.Contains(body, metric) {
			t.Errorf("Expected HTTP metric '%s' not found in metrics output", metric)
		}
	}

	// Check custom business metrics
	for _, metric := range expectedCustomMetrics {
		if !strings.Contains(body, metric) {
			t.Errorf("Expected custom metric '%s' not found in metrics output", metric)
		}
	}

	// Check system metrics (at least some Go runtime metrics should be present)
	foundSystemMetric := false
	for _, metric := range expectedSystemMetrics {
		if strings.Contains(body, metric) {
			foundSystemMetric = true
			break
		}
	}
	if !foundSystemMetric {
		t.Error("Expected at least one Go runtime metric (go_goroutines or go_memstats) not found")
	}
}

// TestMetricsEndpoint_Format verifies metrics are in Prometheus format
func TestMetricsEndpoint_Format(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reg := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = reg
	prometheus.DefaultGatherer = reg

	cfg := config.Load()
	router := gin.New()
	router.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))

	req, _ := http.NewRequest("GET", cfg.Metrics.Path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body := w.Body.String()

	// Check Prometheus format: metric_name{labels} value
	if !strings.Contains(body, "# TYPE") && !strings.Contains(body, "# HELP") {
		t.Error("Metrics output should contain TYPE or HELP comments")
	}

	// Check that we have at least one metric line
	lines := strings.Split(body, "\n")
	hasMetricLine := false
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		if strings.Contains(line, " ") || strings.Contains(line, "{") {
			hasMetricLine = true
			break
		}
	}
	if !hasMetricLine {
		t.Error("Metrics output should contain at least one metric line")
	}
}

// TestMetricsEndpoint_ContentType verifies correct content type
func TestMetricsEndpoint_ContentType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reg := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = reg
	prometheus.DefaultGatherer = reg

	cfg := config.Load()
	router := gin.New()
	router.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))

	req, _ := http.NewRequest("GET", cfg.Metrics.Path, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	contentType := w.Header().Get("Content-Type")
	expectedContentType := "text/plain; version=0.0.4; charset=utf-8"

	if contentType != expectedContentType {
		t.Errorf("Expected Content-Type '%s', got '%s'", expectedContentType, contentType)
	}
}

// TestAllExpectedMetrics verifies comprehensive list of expected metrics
func TestAllExpectedMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reg := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = reg
	prometheus.DefaultGatherer = reg

	cfg := config.Load()
	router := gin.New()
	router.Use(middleware.MetricsMiddleware())
	router.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Exercise all metric types
	metricsService := service.NewMetricsService()
	metricsService.IncrementUserCreated()
	metricsService.IncrementUserUpdated()
	metricsService.IncrementUserDeleted()
	stopTimer := metricsService.StartUserOperationTimer()
	stopTimer()
	metricsService.IncrementUserOperationErrors("test")
	metricsService.RecordExternalCallDuration("test-service", time.Duration(0))
	metricsService.IncrementExternalCallErrors("test-service")

	// Initialize MongoDB metrics (set connection pool config)
	metrics.SetConnectionPoolConfig("go-webapi-db", "test_db", 10, 5)

	// Record some MongoDB operations to ensure metrics are created
	metrics.RecordOperation("go-webapi-db", "test_db", "find", "users", 10*time.Millisecond, nil)
	metrics.RecordOperation("go-webapi-db", "test_db", "insert", "users", 5*time.Millisecond, nil)
	metrics.RecordPing("go-webapi-db", "test_db", 2*time.Millisecond)

	// Make HTTP request
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Get metrics
	req2, _ := http.NewRequest("GET", cfg.Metrics.Path, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	body := w2.Body.String()

	// Comprehensive list of expected metrics
	expectedMetrics := map[string]string{
		// HTTP metrics
		"http_server_requests_total":      "HTTP request counter",
		"http_server_requests_seconds":    "HTTP request duration histogram",
		"http_server_errors_total":        "HTTP server errors counter",
		"http_server_client_errors_total": "HTTP client errors counter",

		// Custom business metrics
		"custom_user_created_total":              "User creation counter",
		"custom_user_updated_total":              "User update counter",
		"custom_user_deleted_total":              "User deletion counter",
		"custom_user_operation_duration_seconds": "User operation duration histogram",
		"custom_user_active_operations":          "Active operations gauge",
		"custom_user_operation_errors_total":     "User operation errors counter",
		"custom_external_call_duration_seconds":  "External call duration histogram",
		"custom_external_call_errors_total":      "External call errors counter",

		// Health metrics
		"health_status": "Health status gauge",

		// MongoDB metrics (these are created when operations are performed)
		"mongodb_connections_max":            "MongoDB max connections gauge",
		"mongodb_connections_min":            "MongoDB min connections gauge",
		"mongodb_operations_total":           "MongoDB operations counter",
		"mongodb_operation_duration_seconds": "MongoDB operation duration histogram",
		"mongodb_ping_duration_seconds":      "MongoDB ping duration histogram",
	}

	missingMetrics := []string{}
	for metricName, description := range expectedMetrics {
		if !strings.Contains(body, metricName) {
			missingMetrics = append(missingMetrics, metricName+" ("+description+")")
		}
	}

	if len(missingMetrics) > 0 {
		t.Errorf("Missing expected metrics:\n%s", strings.Join(missingMetrics, "\n"))
	}
}

// TestMetricsLabels verifies metrics have correct labels
func TestMetricsLabels(t *testing.T) {
	gin.SetMode(gin.TestMode)

	reg := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = reg
	prometheus.DefaultGatherer = reg

	cfg := config.Load()
	router := gin.New()
	router.Use(middleware.MetricsMiddleware())
	router.GET(cfg.Metrics.Path, gin.WrapH(promhttp.Handler()))
	router.GET("/api/users/123", func(c *gin.Context) {
		c.JSON(200, gin.H{"id": "123"})
	})

	metricsService := service.NewMetricsService()
	metricsService.IncrementUserCreated()
	metricsService.IncrementUserOperationErrors("not_found")

	req, _ := http.NewRequest("GET", "/api/users/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req2, _ := http.NewRequest("GET", cfg.Metrics.Path, nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	body := w2.Body.String()

	// Check HTTP metrics have method, uri, status labels
	if strings.Contains(body, "http_server_requests_total") {
		if !strings.Contains(body, "method=") || !strings.Contains(body, "uri=") || !strings.Contains(body, "status=") {
			t.Error("http_server_requests_total should have method, uri, and status labels")
		}
	}

	// Check custom metrics have application label
	if strings.Contains(body, "custom_user_created_total") {
		if !strings.Contains(body, "application=") {
			t.Error("custom_user_created_total should have application label")
		}
	}

	// Check error metrics have error_type label
	if strings.Contains(body, "custom_user_operation_errors_total") {
		if !strings.Contains(body, "error_type=") {
			t.Error("custom_user_operation_errors_total should have error_type label")
		}
	}
}
