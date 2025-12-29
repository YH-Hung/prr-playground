package trace

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	traceID := New()
	if traceID == "" {
		t.Error("New() returned empty trace ID")
	}

	// Generate another to ensure uniqueness
	traceID2 := New()
	if traceID == traceID2 {
		t.Error("New() generated duplicate trace IDs")
	}
}

func TestFromContext(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "trace ID present",
			ctx:      NewContext(context.Background(), "test-trace-id"),
			expected: "test-trace-id",
		},
		{
			name:     "trace ID not present",
			ctx:      context.Background(),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromContext(tt.ctx)
			if got != tt.expected {
				t.Errorf("FromContext() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewContext(t *testing.T) {
	traceID := "test-trace-id-123"
	ctx := NewContext(context.Background(), traceID)

	got := FromContext(ctx)
	if got != traceID {
		t.Errorf("NewContext/FromContext = %v, want %v", got, traceID)
	}
}
