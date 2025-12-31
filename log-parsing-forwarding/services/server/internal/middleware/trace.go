// Package middleware provides HTTP middleware for the server.
package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/yinghanhung/prr-playground/internal/trace"
	"github.com/yinghanhung/prr-playground/services/server/internal/metrics"
)

// StatusRecorder wraps http.ResponseWriter to capture the status code.
type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

// WriteHeader captures the status code before writing it.
func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type logEntry struct {
	TraceID   string `json:"traceId"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	LatencyMs int64  `json:"latencyMs"`
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

// Trace returns middleware that adds trace ID to requests and logs them.
func Trace(stdoutLogger *log.Logger, fileLogger *log.Logger, collector *metrics.Collector, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get or generate trace ID
		traceID := r.Header.Get(trace.HeaderName)
		if traceID == "" {
			traceID = trace.New()
		}

		// Add trace ID to context
		ctx := trace.NewContext(r.Context(), traceID)
		rec := &StatusRecorder{ResponseWriter: w, Status: http.StatusOK}

		// Process request
		next.ServeHTTP(rec, r.WithContext(ctx))

		// Record metrics
		latency := time.Since(start)
		collector.RecordRequest()
		if rec.Status >= 400 {
			collector.RecordError()
		}
		collector.RecordLatency(latency)

		// Log request
		logJSON(stdoutLogger, fileLogger, logEntry{
			TraceID:   traceID,
			Method:    r.Method,
			Path:      r.URL.Path,
			Status:    rec.Status,
			LatencyMs: latency.Milliseconds(),
			Message:   "request completed",
		})
	})
}
