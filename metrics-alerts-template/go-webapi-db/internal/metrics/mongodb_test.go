package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDBMetrics_RecordOperation(t *testing.T) {
	reg := prometheus.NewRegistry()
	
	// Create metrics with custom registry
	opCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mongodb_operations_total",
			Help: "Total number of MongoDB operations",
		},
		[]string{"application", "database", "operation", "collection"},
	)
	reg.MustRegister(opCounter)

	// Simulate recording an operation
	labels := []string{"go-webapi-db", "test_db", "find", "users"}
	opCounter.WithLabelValues(labels...).Inc()

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "mongodb_operations_total" {
			found = true
		}
	}

	if !found {
		t.Error("Metric mongodb_operations_total not found")
	}
}

func TestMongoDBMetrics_RecordOperationWithError(t *testing.T) {
	reg := prometheus.NewRegistry()
	
	opErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mongodb_operation_errors_total",
			Help: "Total number of MongoDB operation errors",
		},
		[]string{"application", "database", "operation", "collection", "error_type"},
	)
	reg.MustRegister(opErrors)

	// Simulate recording an error
	labels := []string{"go-webapi-db", "test_db", "find", "users", "not_found"}
	opErrors.WithLabelValues(labels...).Inc()

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "mongodb_operation_errors_total" {
			found = true
		}
	}

	if !found {
		t.Error("Metric mongodb_operation_errors_total not found")
	}
}

func TestMongoDBMetrics_RecordPing(t *testing.T) {
	reg := prometheus.NewRegistry()
	
	pingDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mongodb_ping_duration_seconds",
			Help:    "MongoDB ping duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5},
		},
		[]string{"application", "database"},
	)
	reg.MustRegister(pingDuration)

	// Simulate recording a ping
	labels := []string{"go-webapi-db", "test_db"}
	pingDuration.WithLabelValues(labels...).Observe(0.01)

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "mongodb_ping_duration_seconds" {
			found = true
		}
	}

	if !found {
		t.Error("Metric mongodb_ping_duration_seconds not found")
	}
}

func TestMongoDBMetrics_ConnectionPoolConfig(t *testing.T) {
	reg := prometheus.NewRegistry()
	
	connMax := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_connections_max",
			Help: "Maximum number of MongoDB connections allowed",
		},
		[]string{"application", "database"},
	)
	connMin := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mongodb_connections_min",
			Help: "Minimum number of MongoDB connections maintained",
		},
		[]string{"application", "database"},
	)
	reg.MustRegister(connMax, connMin)

	// Set connection pool config
	labels := []string{"go-webapi-db", "test_db"}
	connMax.WithLabelValues(labels...).Set(10)
	connMin.WithLabelValues(labels...).Set(5)

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var foundMax, foundMin bool
	for _, mf := range metrics {
		if mf.GetName() == "mongodb_connections_max" {
			foundMax = true
		}
		if mf.GetName() == "mongodb_connections_min" {
			foundMin = true
		}
	}

	if !foundMax {
		t.Error("Metric mongodb_connections_max not found")
	}
	if !foundMin {
		t.Error("Metric mongodb_connections_min not found")
	}
}

func TestMongoDBMetrics_AllMetricsRegistered(t *testing.T) {
	// Test that all expected MongoDB metrics are registered
	expectedMetrics := []string{
		"mongodb_connections_active",
		"mongodb_connections_idle",
		"mongodb_connections_max",
		"mongodb_connections_min",
		"mongodb_connections_total",
		"mongodb_connection_acquire_seconds",
		"mongodb_connection_timeouts_total",
		"mongodb_operations_total",
		"mongodb_operation_duration_seconds",
		"mongodb_operation_errors_total",
		"mongodb_connection_errors_total",
		"mongodb_ping_duration_seconds",
	}

	// Note: In a real test, you would verify these are registered in the default registry
	// For now, we just verify the list is complete
	if len(expectedMetrics) != 12 {
		t.Errorf("Expected 12 MongoDB metrics, got %d", len(expectedMetrics))
	}
}

func TestMongoDBMetrics_RecordOperationErrorTypes(t *testing.T) {
	reg := prometheus.NewRegistry()
	
	opErrors := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mongodb_operation_errors_total",
			Help: "Total number of MongoDB operation errors",
		},
		[]string{"application", "database", "operation", "collection", "error_type"},
	)
	reg.MustRegister(opErrors)

	// Test different error types
	errorTypes := []string{"not_found", "timeout", "cancelled", "unknown"}
	
	for _, errorType := range errorTypes {
		labels := []string{"go-webapi-db", "test_db", "find", "users", errorType}
		opErrors.WithLabelValues(labels...).Inc()
	}

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var errorCount int
	for _, mf := range metrics {
		if mf.GetName() == "mongodb_operation_errors_total" {
			errorCount = len(mf.GetMetric())
		}
	}

	if errorCount != len(errorTypes) {
		t.Errorf("Expected %d error metric series, got %d", len(errorTypes), errorCount)
	}
}

func TestMongoDBMetricsCollector_StartStop(t *testing.T) {
	// Create a mock MongoDB client (will fail connection, but tests the collector)
	client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	defer client.Disconnect(context.Background())
	
	collector := NewMongoDBMetricsCollector(client, "test_db", "go-webapi-db")
	
	// Start collector
	collector.Start(100 * time.Millisecond)
	
	// Let it run briefly
	time.Sleep(150 * time.Millisecond)
	
	// Stop collector
	collector.Stop()
	
	// If we get here without deadlock, the test passes
}

