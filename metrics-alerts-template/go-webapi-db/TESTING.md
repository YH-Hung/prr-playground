# Testing Guide

This document describes the test suite for verifying that all expected metrics are properly exported.

## Test Files

### Unit Tests

1. **`internal/service/metrics_service_test.go`**
   - Tests all custom business metrics
   - Verifies counters, histograms, and gauges are properly registered
   - Tests metric values and labels

2. **`internal/middleware/metrics_test.go`**
   - Tests HTTP metrics middleware
   - Verifies HTTP request metrics (counters, histograms, error metrics)
   - Tests URI sanitization
   - Verifies label values

3. **`internal/health/health_test.go`**
   - Tests health check metrics
   - Verifies health status gauge
   - Tests component labels

### Integration Tests

4. **`cmd/server/metrics_integration_test.go`**
   - Comprehensive test verifying all metrics are exported via `/metrics` endpoint
   - Tests Prometheus format compliance
   - Verifies content type
   - Tests metric labels

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Run Specific Test Package

```bash
# Test metrics service
go test ./internal/service/...

# Test middleware
go test ./internal/middleware/...

# Test health checks
go test ./internal/health/...

# Test integration
go test ./cmd/server/...
```

### Run Specific Test Function

```bash
go test -v ./internal/service -run TestMetricsService_AllMetricsRegistered
```

## Expected Metrics Coverage

The test suite verifies the following metrics are exported:

### HTTP Metrics
- ✅ `http_server_requests_total` - Request counter
- ✅ `http_server_requests_seconds` - Request duration histogram
- ✅ `http_server_errors_total` - Server errors (5xx)
- ✅ `http_server_client_errors_total` - Client errors (4xx)

### Custom Business Metrics
- ✅ `custom_user_created_total` - User creation counter
- ✅ `custom_user_updated_total` - User update counter
- ✅ `custom_user_deleted_total` - User deletion counter
- ✅ `custom_user_operation_duration_seconds` - Operation duration histogram
- ✅ `custom_user_active_operations` - Active operations gauge
- ✅ `custom_user_operation_errors_total` - Operation errors counter
- ✅ `custom_external_call_duration_seconds` - External call duration histogram
- ✅ `custom_external_call_errors_total` - External call errors counter

### Health Metrics
- ✅ `health_status` - Health status gauge with component label

### System Metrics (Go Runtime)
- ✅ `go_goroutines` - Goroutine count
- ✅ `go_memstats_*` - Memory statistics
- ✅ `process_*` - Process metrics (via prometheus.NewProcessCollector)

## Test Coverage

To generate test coverage report:

```bash
go test -cover ./...
```

To generate detailed coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Continuous Integration

These tests should be run in CI/CD pipelines to ensure:
1. All metrics are properly registered
2. Metrics follow Prometheus naming conventions
3. Labels are correctly applied
4. Metrics endpoint returns valid Prometheus format

## Troubleshooting

### Tests Fail Due to Missing Dependencies

Run `go mod tidy` to download dependencies:

```bash
go mod tidy
```

### Tests Fail Due to Registry Conflicts

The tests use isolated Prometheus registries to avoid conflicts. If you see registry-related errors, ensure tests are using `prometheus.NewRegistry()` instead of the default registry.

### Metrics Not Found in Integration Tests

Ensure that:
1. Metrics are registered before making HTTP requests
2. The metrics service is initialized
3. HTTP requests are made to trigger middleware metrics

