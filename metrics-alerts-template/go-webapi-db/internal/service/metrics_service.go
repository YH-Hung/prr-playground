package service

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsService struct {
	userCreatedCounter         *prometheus.CounterVec
	userUpdatedCounter         *prometheus.CounterVec
	userDeletedCounter         *prometheus.CounterVec
	userOperationDuration      *prometheus.HistogramVec
	userActiveOperations       prometheus.Gauge
	userOperationErrorsCounter *prometheus.CounterVec
	externalCallDuration       *prometheus.HistogramVec
	externalCallErrorsCounter  *prometheus.CounterVec
	activeOperationsCount      int64
	mu                         sync.Mutex
}

func NewMetricsService() *MetricsService {
	return &MetricsService{
		userCreatedCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_user_created_total",
				Help: "Total number of users created",
			},
			[]string{"application", "operation"},
		),
		userUpdatedCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_user_updated_total",
				Help: "Total number of users updated",
			},
			[]string{"application", "operation"},
		),
		userDeletedCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_user_deleted_total",
				Help: "Total number of users deleted",
			},
			[]string{"application", "operation"},
		),
		userOperationDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "custom_user_operation_duration_seconds",
				Help:    "Duration of user operations in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"application"},
		),
		userActiveOperations: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "custom_user_active_operations",
				Help: "Number of active user operations",
			},
		),
		userOperationErrorsCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_user_operation_errors_total",
				Help: "Total number of user operation errors",
			},
			[]string{"application", "error_type"},
		),
		externalCallDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "custom_external_call_duration_seconds",
				Help:    "Duration of external service calls in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"application", "service"},
		),
		externalCallErrorsCounter: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "custom_external_call_errors_total",
				Help: "Total number of external service call failures",
			},
			[]string{"application", "service"},
		),
	}
}

func (m *MetricsService) IncrementUserCreated() {
	m.userCreatedCounter.WithLabelValues("go-webapi-db", "create").Inc()
}

func (m *MetricsService) IncrementUserUpdated() {
	m.userUpdatedCounter.WithLabelValues("go-webapi-db", "update").Inc()
}

func (m *MetricsService) IncrementUserDeleted() {
	m.userDeletedCounter.WithLabelValues("go-webapi-db", "delete").Inc()
}

func (m *MetricsService) IncrementUserOperationErrors(errorType string) {
	m.userOperationErrorsCounter.WithLabelValues("go-webapi-db", errorType).Inc()
}

func (m *MetricsService) StartUserOperationTimer() func() {
	m.mu.Lock()
	m.activeOperationsCount++
	m.userActiveOperations.Set(float64(m.activeOperationsCount))
	m.mu.Unlock()

	start := time.Now()
	return func() {
		duration := time.Since(start).Seconds()
		m.userOperationDuration.WithLabelValues("go-webapi-db").Observe(duration)

		m.mu.Lock()
		m.activeOperationsCount--
		m.userActiveOperations.Set(float64(m.activeOperationsCount))
		m.mu.Unlock()
	}
}

func (m *MetricsService) RecordExternalCallDuration(service string, duration time.Duration) {
	m.externalCallDuration.WithLabelValues("go-webapi-db", service).Observe(duration.Seconds())
}

func (m *MetricsService) IncrementExternalCallErrors(service string) {
	m.externalCallErrorsCounter.WithLabelValues("go-webapi-db", service).Inc()
}
