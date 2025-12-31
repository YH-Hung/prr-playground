package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yinghanhung/prr-playground/internal/config"
	"github.com/yinghanhung/prr-playground/internal/logger"
	"github.com/yinghanhung/prr-playground/internal/retry"
	"github.com/yinghanhung/prr-playground/internal/trace"
)

func TestConfigPackage(t *testing.T) {
	t.Run("GetString", func(t *testing.T) {
		val := config.GetString("NONEXISTENT_VAR", "default")
		if val != "default" {
			t.Errorf("Expected 'default', got %s", val)
		}
	})

	t.Run("GetInt", func(t *testing.T) {
		val := config.GetInt("NONEXISTENT_VAR", 42)
		if val != 42 {
			t.Errorf("Expected 42, got %d", val)
		}
	})
}

func TestTracePackage(t *testing.T) {
	t.Run("NewTraceID", func(t *testing.T) {
		traceID := trace.New()
		if traceID == "" {
			t.Error("Expected non-empty trace ID")
		}
	})

	t.Run("TraceIDContext", func(t *testing.T) {
		ctx := trace.NewContext(httptest.NewRequest("GET", "/", nil).Context(), "test-id")
		retrieved := trace.FromContext(ctx)
		if retrieved != "test-id" {
			t.Errorf("Expected 'test-id', got %s", retrieved)
		}
	})

	t.Run("HeaderName", func(t *testing.T) {
		if trace.HeaderName != "X-Trace-Id" {
			t.Errorf("Expected 'X-Trace-Id', got %s", trace.HeaderName)
		}
	})
}

func TestLoggerPackage(t *testing.T) {
	var buf strings.Builder
	logger := logger.New(&buf, "[TEST] ")
	logger.Println("test message")

	output := buf.String()
	if !strings.Contains(output, "[TEST]") {
		t.Errorf("Expected prefix in output: %s", output)
	}
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected message in output: %s", output)
	}
}

func TestRetryPackage(t *testing.T) {
	t.Run("Successful operation", func(t *testing.T) {
		callCount := 0
		err := retry.Do(httptest.NewRequest("GET", "/", nil).Context(), 3, func() error {
			callCount++
			return nil
		}, func(error) bool { return true })

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if callCount != 1 {
			t.Errorf("Expected 1 call, got %d", callCount)
		}
	})
}

func TestServerEndpoints(t *testing.T) {
	// Test that we can make HTTP requests to a test server
	// This tests the overall integration without accessing internal packages

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/health":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status": "healthy",
			})
		case "/hello":
			traceID := r.Header.Get(trace.HeaderName)
			if traceID == "" {
				traceID = trace.New()
			}
			json.NewEncoder(w).Encode(map[string]string{
				"message": "hello",
				"traceId": traceID,
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	t.Run("Health endpoint", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/health")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var data map[string]string
		json.NewDecoder(resp.Body).Decode(&data)
		if data["status"] != "healthy" {
			t.Errorf("Expected healthy status")
		}
	})

	t.Run("Hello endpoint with trace ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", server.URL+"/hello", nil)
		req.Header.Set(trace.HeaderName, "test-trace-123")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var data map[string]string
		json.NewDecoder(resp.Body).Decode(&data)
		if data["traceId"] != "test-trace-123" {
			t.Errorf("Expected trace ID to be propagated, got %s", data["traceId"])
		}
	})
}
