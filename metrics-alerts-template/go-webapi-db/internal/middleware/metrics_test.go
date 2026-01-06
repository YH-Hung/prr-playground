package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(MetricsMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/test/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "test error"})
	})
	r.GET("/test/client-error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
	})
	r.GET("/api/users/123", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"id": "123"})
	})
	r.GET("/api/users/email/test@example.com", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"email": "test@example.com"})
	})
	return r
}

func TestMetricsMiddleware_HTTPRequestsTotal(t *testing.T) {
	// Reset metrics
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Gather metrics
	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "http_server_requests_total" {
			found = true
			var count float64
			for _, metric := range mf.GetMetric() {
				count += metric.GetCounter().GetValue()
			}
			if count == 0 {
				t.Error("Counter should be incremented")
			}
		}
	}

	if !found {
		t.Error("Metric http_server_requests_total not found")
	}
}

func TestMetricsMiddleware_HTTPRequestDuration(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "http_server_requests_seconds" {
			found = true
		}
	}

	if !found {
		t.Error("Metric http_server_requests_seconds not found")
	}
}

func TestMetricsMiddleware_HTTPServerErrors(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/test/error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	var count float64
	for _, mf := range metrics {
		if mf.GetName() == "http_server_errors_total" {
			found = true
			for _, metric := range mf.GetMetric() {
				count += metric.GetCounter().GetValue()
			}
		}
	}

	if !found {
		t.Error("Metric http_server_errors_total not found")
	}
	if count == 0 {
		t.Error("Server error counter should be incremented for 5xx status")
	}
}

func TestMetricsMiddleware_HTTPClientErrors(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/test/client-error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	var count float64
	for _, mf := range metrics {
		if mf.GetName() == "http_server_client_errors_total" {
			found = true
			for _, metric := range mf.GetMetric() {
				count += metric.GetCounter().GetValue()
			}
		}
	}

	if !found {
		t.Error("Metric http_server_client_errors_total not found")
	}
	if count == 0 {
		t.Error("Client error counter should be incremented for 4xx status")
	}
}

func TestMetricsMiddleware_AllMetricsRegistered(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	r := setupRouter()

	// Make requests to trigger all metric types
	req1, _ := http.NewRequest("GET", "/test", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	req2, _ := http.NewRequest("GET", "/test/error", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	req3, _ := http.NewRequest("GET", "/test/client-error", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	expectedMetrics := map[string]bool{
		"http_server_requests_total":        false,
		"http_server_requests_seconds":      false,
		"http_server_errors_total":           false,
		"http_server_client_errors_total":    false,
	}

	for _, mf := range metrics {
		if _, exists := expectedMetrics[mf.GetName()]; exists {
			expectedMetrics[mf.GetName()] = true
		}
	}

	for metricName, found := range expectedMetrics {
		if !found {
			t.Errorf("Expected metric %s not found", metricName)
		}
	}
}

func TestSanitizeURI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "numeric ID",
			input:    "/api/users/123",
			expected: "/api/users/{id}",
		},
		{
			name:     "email pattern",
			input:    "/api/users/email/test@example.com",
			expected: "/api/users/email/{email}",
		},
		{
			name:     "status pattern",
			input:    "/api/users/status/ACTIVE",
			expected: "/api/users/status/{status}",
		},
		{
			name:     "status count pattern",
			input:    "/api/users/status/ACTIVE/count",
			expected: "/api/users/status/{status}/count",
		},
		{
			name:     "external service pattern",
			input:    "/api/users/external/test-service",
			expected: "/api/users/external/{serviceName}",
		},
		{
			name:     "test error pattern",
			input:    "/api/users/test/error",
			expected: "/api/users/test/error",
		},
		{
			name:     "regular path",
			input:    "/health",
			expected: "/health",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeURI(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeURI(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMetricsMiddleware_LabelValues(t *testing.T) {
	prometheus.DefaultRegisterer = prometheus.NewRegistry()
	
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/api/users/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "http_server_requests_total" {
			found = true
			for _, metric := range mf.GetMetric() {
				labels := metric.GetLabel()
				var method, uri, status string
				for _, label := range labels {
					switch label.GetName() {
					case "method":
						method = label.GetValue()
					case "uri":
						uri = label.GetValue()
					case "status":
						status = label.GetValue()
					}
				}
				if method != "GET" {
					t.Errorf("Expected method GET, got %s", method)
				}
				if uri != "/api/users/{id}" {
					t.Errorf("Expected URI /api/users/{id}, got %s", uri)
				}
				if status != "200" {
					t.Errorf("Expected status 200, got %s", status)
				}
			}
		}
	}

	if !found {
		t.Error("Metric with expected labels not found")
	}
}

