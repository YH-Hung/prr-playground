package service

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func TestMetricsService_UserCreatedCounter(t *testing.T) {
	// Create a new registry for testing
	reg := prometheus.NewRegistry()
	
	// Create metrics service with custom registry
	ms := &MetricsService{
		userCreatedCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_user_created_total",
				Help: "Total number of users created",
			},
			[]string{"application", "operation"},
		),
	}
	reg.MustRegister(ms.userCreatedCounter)

	// Test increment
	ms.IncrementUserCreated()
	ms.IncrementUserCreated()

	// Verify metric exists and has correct value
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	var value float64
	for _, mf := range metrics {
		if mf.GetName() == "custom_user_created_total" {
			found = true
			for _, metric := range mf.GetMetric() {
				value += metric.GetCounter().GetValue()
			}
		}
	}

	if !found {
		t.Error("Metric custom_user_created_total not found")
	}
	if value != 2 {
		t.Errorf("Expected counter value 2, got %f", value)
	}
}

func TestMetricsService_UserUpdatedCounter(t *testing.T) {
	reg := prometheus.NewRegistry()
	ms := &MetricsService{
		userUpdatedCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_user_updated_total",
				Help: "Total number of users updated",
			},
			[]string{"application", "operation"},
		),
	}
	reg.MustRegister(ms.userUpdatedCounter)

	ms.IncrementUserUpdated()

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "custom_user_updated_total" {
			found = true
		}
	}

	if !found {
		t.Error("Metric custom_user_updated_total not found")
	}
}

func TestMetricsService_UserDeletedCounter(t *testing.T) {
	reg := prometheus.NewRegistry()
	ms := &MetricsService{
		userDeletedCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_user_deleted_total",
				Help: "Total number of users deleted",
			},
			[]string{"application", "operation"},
		),
	}
	reg.MustRegister(ms.userDeletedCounter)

	ms.IncrementUserDeleted()

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "custom_user_deleted_total" {
			found = true
		}
	}

	if !found {
		t.Error("Metric custom_user_deleted_total not found")
	}
}

func TestMetricsService_UserOperationDuration(t *testing.T) {
	reg := prometheus.NewRegistry()
	ms := &MetricsService{
		userOperationDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "custom_user_operation_duration_seconds",
				Help:    "Duration of user operations in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"application"},
		),
	}
	reg.MustRegister(ms.userOperationDuration)

	stopTimer := ms.StartUserOperationTimer()
	time.Sleep(10 * time.Millisecond)
	stopTimer()

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "custom_user_operation_duration_seconds" {
			found = true
			// Check that histogram has buckets
			if len(mf.GetMetric()) == 0 {
				t.Error("Histogram has no metrics")
			}
		}
	}

	if !found {
		t.Error("Metric custom_user_operation_duration_seconds not found")
	}
}

func TestMetricsService_UserActiveOperations(t *testing.T) {
	reg := prometheus.NewRegistry()
	ms := &MetricsService{
		userActiveOperations: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "custom_user_active_operations",
				Help: "Number of active user operations",
			},
		),
		activeOperationsCount: 0,
	}
	reg.MustRegister(ms.userActiveOperations)

	// Start operation
	ms.mu.Lock()
	ms.activeOperationsCount++
	ms.userActiveOperations.Set(float64(ms.activeOperationsCount))
	ms.mu.Unlock()

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	var value float64
	for _, mf := range metrics {
		if mf.GetName() == "custom_user_active_operations" {
			found = true
			for _, metric := range mf.GetMetric() {
				value = metric.GetGauge().GetValue()
			}
		}
	}

	if !found {
		t.Error("Metric custom_user_active_operations not found")
	}
	if value != 1 {
		t.Errorf("Expected gauge value 1, got %f", value)
	}
}

func TestMetricsService_UserOperationErrors(t *testing.T) {
	reg := prometheus.NewRegistry()
	ms := &MetricsService{
		userOperationErrorsCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_user_operation_errors_total",
				Help: "Total number of user operation errors",
			},
			[]string{"application", "error_type"},
		),
	}
	reg.MustRegister(ms.userOperationErrorsCounter)

	ms.IncrementUserOperationErrors("not_found")
	ms.IncrementUserOperationErrors("duplicate_email")

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	var errorCount int
	for _, mf := range metrics {
		if mf.GetName() == "custom_user_operation_errors_total" {
			found = true
			errorCount = len(mf.GetMetric())
		}
	}

	if !found {
		t.Error("Metric custom_user_operation_errors_total not found")
	}
	if errorCount != 2 {
		t.Errorf("Expected 2 error metric series, got %d", errorCount)
	}
}

func TestMetricsService_ExternalCallDuration(t *testing.T) {
	reg := prometheus.NewRegistry()
	ms := &MetricsService{
		externalCallDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "custom_external_call_duration_seconds",
				Help:    "Duration of external service calls in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"application", "service"},
		),
	}
	reg.MustRegister(ms.externalCallDuration)

	ms.RecordExternalCallDuration("test-service", 100*time.Millisecond)

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "custom_external_call_duration_seconds" {
			found = true
		}
	}

	if !found {
		t.Error("Metric custom_external_call_duration_seconds not found")
	}
}

func TestMetricsService_ExternalCallErrors(t *testing.T) {
	reg := prometheus.NewRegistry()
	ms := &MetricsService{
		externalCallErrorsCounter: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_external_call_errors_total",
				Help: "Total number of external service call failures",
			},
			[]string{"application", "service"},
		),
	}
	reg.MustRegister(ms.externalCallErrorsCounter)

	ms.IncrementExternalCallErrors("test-service")

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	var found bool
	for _, mf := range metrics {
		if mf.GetName() == "custom_external_call_errors_total" {
			found = true
		}
	}

	if !found {
		t.Error("Metric custom_external_call_errors_total not found")
	}
}

func TestMetricsService_AllMetricsRegistered(t *testing.T) {
	ms := NewMetricsService()
	reg := prometheus.NewRegistry()

	// Register all metrics
	reg.MustRegister(ms.userCreatedCounter)
	reg.MustRegister(ms.userUpdatedCounter)
	reg.MustRegister(ms.userDeletedCounter)
	reg.MustRegister(ms.userOperationDuration)
	reg.MustRegister(ms.userActiveOperations)
	reg.MustRegister(ms.userOperationErrorsCounter)
	reg.MustRegister(ms.externalCallDuration)
	reg.MustRegister(ms.externalCallErrorsCounter)

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	expectedMetrics := map[string]bool{
		"custom_user_created_total":              false,
		"custom_user_updated_total":              false,
		"custom_user_deleted_total":               false,
		"custom_user_operation_duration_seconds":  false,
		"custom_user_active_operations":          false,
		"custom_user_operation_errors_total":     false,
		"custom_external_call_duration_seconds":   false,
		"custom_external_call_errors_total":      false,
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

func TestMetricsService_StartUserOperationTimer(t *testing.T) {
	ms := NewMetricsService()
	reg := prometheus.NewRegistry()
	reg.MustRegister(ms.userOperationDuration)
	reg.MustRegister(ms.userActiveOperations)

	// Start timer
	stopTimer := ms.StartUserOperationTimer()

	// Verify active operations increased
	metrics, _ := reg.Gather()
	var activeOps float64
	for _, mf := range metrics {
		if mf.GetName() == "custom_user_active_operations" {
			for _, metric := range mf.GetMetric() {
				activeOps = metric.GetGauge().GetValue()
			}
		}
	}
	if activeOps != 1 {
		t.Errorf("Expected active operations to be 1, got %f", activeOps)
	}

	// Stop timer
	stopTimer()

	// Verify active operations decreased
	metrics, _ = reg.Gather()
	activeOps = 0
	for _, mf := range metrics {
		if mf.GetName() == "custom_user_active_operations" {
			for _, metric := range mf.GetMetric() {
				activeOps = metric.GetGauge().GetValue()
			}
		}
	}
	if activeOps != 0 {
		t.Errorf("Expected active operations to be 0, got %f", activeOps)
	}

	// Verify duration was recorded
	var durationRecorded bool
	for _, mf := range metrics {
		if mf.GetName() == "custom_user_operation_duration_seconds" {
			durationRecorded = true
		}
	}
	if !durationRecorded {
		t.Error("Duration was not recorded")
	}
}

