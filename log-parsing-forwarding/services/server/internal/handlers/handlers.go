// Package handlers provides HTTP request handlers for the server.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yinghanhung/prr-playground/internal/trace"
	"github.com/yinghanhung/prr-playground/services/server/internal/metrics"
)

type logEntry struct {
	TraceID   string `json:"traceId"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	LatencyMs int64  `json:"latencyMs,omitempty"`
	Message   string `json:"message"`
}

func logJSON(stdoutLogger *log.Logger, fileLogger *log.Logger, entry logEntry) {
	b, err := json.Marshal(entry)
	if err != nil {
		stdoutLogger.Printf(`{"message":"failed to marshal log","error":"%v"}\n`, err)
		fileLogger.Printf(`{"message":"failed to marshal log","error":"%v"}\n`, err)
		return
	}
	stdoutLogger.Println(string(b))
	fileLogger.Printf("%s\n", string(b))
}

// Hello returns a handler for the main hello endpoint.
func Hello(stdoutLogger *log.Logger, fileLogger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		traceID := trace.FromContext(r.Context())
		resp := map[string]string{
			"message": "hello",
			"traceId": traceID,
			"path":    r.URL.Path,
		}
		time.Sleep(50 * time.Millisecond) // simulate work

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logJSON(stdoutLogger, fileLogger, logEntry{
				TraceID: traceID,
				Method:  r.Method,
				Path:    r.URL.Path,
				Status:  http.StatusInternalServerError,
				Message: "failed to encode response",
			})
			return
		}

		logJSON(stdoutLogger, fileLogger, logEntry{
			TraceID: traceID,
			Method:  r.Method,
			Path:    r.URL.Path,
			Status:  http.StatusOK,
			Message: "handler finished",
		})
	}
}

// Health returns a handler for the health check endpoint.
func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "prr-playground-server",
		})
	}
}

// Metrics returns a handler for the metrics endpoint.
func Metrics(collector *metrics.Collector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats := collector.GetStats()

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "# HELP http_requests_total Total number of HTTP requests\n")
		fmt.Fprintf(w, "# TYPE http_requests_total counter\n")
		fmt.Fprintf(w, "http_requests_total %d\n", stats.RequestCount)
		fmt.Fprintf(w, "# HELP http_errors_total Total number of HTTP errors (4xx, 5xx)\n")
		fmt.Fprintf(w, "# TYPE http_errors_total counter\n")
		fmt.Fprintf(w, "http_errors_total %d\n", stats.ErrorCount)
		fmt.Fprintf(w, "# HELP http_request_duration_ms Average request latency in milliseconds\n")
		fmt.Fprintf(w, "# TYPE http_request_duration_ms gauge\n")
		fmt.Fprintf(w, "http_request_duration_ms %d\n", stats.AvgLatencyMs)
	}
}
