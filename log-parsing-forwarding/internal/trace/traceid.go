// Package trace provides distributed tracing utilities with trace ID generation and propagation.
package trace

import (
	"context"

	"github.com/google/uuid"
)

// HeaderName is the standard HTTP header name for trace IDs.
const HeaderName = "X-Trace-Id"

type ctxKey string

const traceIDKey ctxKey = "traceId"

// New generates a new trace ID using UUID v4.
func New() string {
	return uuid.New().String()
}

// FromContext retrieves the trace ID from the context, or returns an empty string if not found.
func FromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value(traceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// NewContext creates a new context with the given trace ID.
func NewContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}
