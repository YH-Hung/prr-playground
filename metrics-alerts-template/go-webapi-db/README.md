# Go WebAPI MongoDB Metrics Project

A production-ready Go web API built with Gin framework, MongoDB, and comprehensive Prometheus metrics and alerting capabilities. This project mirrors the Spring Boot `spring-webapi-db` project's metrics and alerting features.

## Features

### Metrics Collection
- **HTTP Metrics**: Request rate, latency (p50, p95, p99), error rates, status codes
- **Custom Business Metrics**: User operations (create, update, delete), operation duration, error tracking
- **System Metrics**: Go runtime metrics (memory, GC, goroutines), process metrics (CPU, uptime, file descriptors)
- **Database Metrics**: MongoDB connection pool metrics, operation metrics, health check status
- **External Service Metrics**: External service call duration and error tracking

### Alerting
- **High Error Rate**: Alert when error rate exceeds 5%
- **High Latency**: Alert when p95 latency exceeds 1 second
- **Database Issues**: Connection failures, health check failures
- **Memory Pressure**: Memory usage warnings (85%) and critical alerts (95%)
- **Service Availability**: Service down, health endpoint failures
- **Resource Exhaustion**: Goroutine count, file descriptor usage
- **Business Metrics**: High business operation error rates, external service failures

### API Endpoints
- `POST /api/users` - Create user
- `GET /api/users/{id}` - Get user by ID
- `GET /api/users` - Get all users
- `GET /api/users/email/{email}` - Get user by email
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user
- `GET /api/users/status/{status}` - Get users by status
- `GET /api/users/status/{status}/count` - Count users by status
- `GET /api/users/external/{serviceName}` - Simulate external service call
- `GET /api/users/test/error` - Test error endpoint
- `GET /api/users/test/slow` - Test slow response endpoint
- `GET /metrics` - Prometheus metrics endpoint
- `GET /health` - Health check endpoint

## Prerequisites

- **Go 1.21+**
- **Docker and Docker Compose**
- **MongoDB** (included in Docker Compose)

## Quick Start

### 1. Start the Full Stack (Application + Monitoring)

Start MongoDB, the Go application, Prometheus, Grafana, and Alertmanager:

```bash
cd go-webapi-db
docker compose -f docker-compose.metrics.yml up -d
```

**Note**: If you're using an older version of Docker, you may need to use `docker-compose` (with hyphen) instead of `docker compose`.

This will start:
- **MongoDB** on `localhost:27017`
- **Go WebAPI** on `http://localhost:8080`
- **Prometheus** on `http://localhost:9090`
- **Grafana** on `http://localhost:3000` (admin/admin)
- **Alertmanager** on `http://localhost:9093`

### 2. Running Locally (Without Docker)

#### Start MongoDB

```bash
# Using Docker
docker run -d -p 27017:27017 --name mongodb mongo:7

# Check if MongoDB is running
docker ps | grep mongodb

# View MongoDB logs if needed
docker logs mongodb

# Stop MongoDB container (when done)
# docker stop mongodb

# Remove MongoDB container (when done)
# docker rm mongodb

# Or use your local MongoDB instance
```

#### Build and Run

```bash
cd go-webapi-db

# Download dependencies and generate go.sum
go mod tidy

# Build
go build -o bin/server ./cmd/server

# Or use the build script
./build.sh

# Run
./bin/server
```

**Note**: 
- If you encounter compilation errors about missing `go.sum`, run `go mod tidy` first to download dependencies and generate the checksum file
- The `bin/` directory will be created automatically by the build script, or create it manually with `mkdir -p bin` before building

Or run directly:

```bash
go run ./cmd/server
```

#### Environment Variables

```bash
export SERVER_PORT=8080
export MONGODB_URI=mongodb://localhost:27017
export MONGODB_DATABASE=go_webapi_db
export METRICS_PATH=/metrics
```

### 3. Start Only the Monitoring Stack

If you're running the application locally:

```bash
cd go-webapi-db
docker compose -f docker-compose.metrics.yml up prometheus alertmanager grafana -d
```

Update `deployments/prometheus/prometheus.yml` to use `host.docker.internal:8080` instead of `go-webapi-db:8080`:

```yaml
scrape_configs:
  - job_name: 'go-webapi-db'
    static_configs:
      - targets: ['host.docker.internal:8080']  # For local app
```

## Project Structure

```
go-webapi-db/
├── cmd/
│   └── server/
│       ├── main.go                      # Application entry point
│       └── metrics_integration_test.go  # Integration tests
├── internal/
│   ├── config/
│   │   └── config.go                    # Configuration management
│   ├── handler/
│   │   └── user_handler.go              # HTTP handlers
│   ├── middleware/
│   │   ├── metrics.go                    # Prometheus HTTP metrics middleware
│   │   ├── metrics_test.go               # Middleware tests
│   │   └── recovery.go                   # Panic recovery middleware
│   ├── model/
│   │   └── user.go                      # User data model
│   ├── repository/
│   │   ├── user_repository.go           # MongoDB repository layer
│   │   ├── instrumented_repository.go    # Instrumented repository with metrics
│   │   └── interface.go                 # Repository interface
│   ├── service/
│   │   ├── metrics_service.go           # Custom business metrics service
│   │   ├── metrics_service_test.go      # Metrics service tests
│   │   └── user_service.go              # Business logic layer
│   ├── metrics/
│   │   ├── mongodb.go                   # MongoDB metrics collector
│   │   └── mongodb_test.go             # MongoDB metrics tests
│   └── health/
│       ├── health.go                     # Health check handlers
│       └── health_test.go               # Health check tests
├── deployments/
│   ├── prometheus/
│   │   ├── prometheus.yml               # Prometheus configuration
│   │   └── alerts.yml                   # Alert rules
│   ├── alertmanager/
│   │   └── alertmanager.yml             # Alertmanager configuration
│   └── grafana/
│       └── provisioning/                # Grafana provisioning configs
│           ├── datasources/
│           │   └── prometheus.yml
│           └── dashboards/
│               └── dashboard.yml
├── docker-compose.metrics.yml           # Docker Compose for full stack
├── Dockerfile                            # Application Dockerfile
├── build.sh                              # Build script
├── go.mod                                # Go module dependencies
├── go.sum                                # Go module checksums
├── .gitignore                            # Git ignore rules
├── .dockerignore                         # Docker ignore rules
├── README.md                              # This file
└── TESTING.md                             # Testing documentation
```

## Metrics Documentation

### HTTP Server Metrics

#### `http_server_requests_total`
- **Type**: Counter
- **Description**: Total number of HTTP requests
- **Labels**: `method`, `uri`, `status`
- **Use Case**: Track request count per endpoint

#### `http_server_requests_seconds`
- **Type**: Histogram
- **Description**: HTTP request duration in seconds
- **Labels**: `method`, `uri`, `status`
- **Use Case**: Calculate latency percentiles (p50, p95, p99)

#### `http_server_errors_total`
- **Type**: Counter
- **Description**: Total number of HTTP server errors (5xx)
- **Labels**: `method`, `uri`, `status`
- **Use Case**: Monitor server error rate

#### `http_server_client_errors_total`
- **Type**: Counter
- **Description**: Total number of HTTP client errors (4xx)
- **Labels**: `method`, `uri`, `status`
- **Use Case**: Track client error patterns

### Custom Business Metrics

#### `custom_user_created_total`
- **Type**: Counter
- **Description**: Total number of users created
- **Labels**: `application`, `operation`
- **Use Case**: Track user creation rate

#### `custom_user_updated_total`
- **Type**: Counter
- **Description**: Total number of users updated
- **Labels**: `application`, `operation`
- **Use Case**: Monitor user update operations

#### `custom_user_deleted_total`
- **Type**: Counter
- **Description**: Total number of users deleted
- **Labels**: `application`, `operation`
- **Use Case**: Track user deletion rate

#### `custom_user_operation_duration_seconds`
- **Type**: Histogram
- **Description**: Duration of user operations in seconds
- **Labels**: `application`
- **Use Case**: Monitor operation performance, calculate percentiles

#### `custom_user_active_operations`
- **Type**: Gauge
- **Description**: Number of active user operations
- **Labels**: None
- **Use Case**: Monitor concurrent operations

#### `custom_user_operation_errors_total`
- **Type**: Counter
- **Description**: Total number of user operation errors
- **Labels**: `application`, `error_type`
- **Use Case**: Track error rates by type

#### `custom_external_call_duration_seconds`
- **Type**: Histogram
- **Description**: Duration of external service calls in seconds
- **Labels**: `application`, `service`
- **Use Case**: Monitor external service performance

#### `custom_external_call_errors_total`
- **Type**: Counter
- **Description**: Total number of external service call failures
- **Labels**: `application`, `service`
- **Use Case**: Track external service reliability

### System Metrics (Go Runtime)

#### `go_memstats_heap_inuse_bytes`
- **Type**: Gauge
- **Description**: Bytes of heap memory in use
- **Use Case**: Monitor memory usage

#### `go_memstats_heap_sys_bytes`
- **Type**: Gauge
- **Description**: Bytes of heap memory obtained from OS
- **Use Case**: Calculate memory usage percentage

#### `go_memstats_gc_duration_seconds`
- **Type**: Summary
- **Description**: GC pause duration
- **Use Case**: Monitor GC performance

#### `go_goroutines`
- **Type**: Gauge
- **Description**: Number of goroutines
- **Use Case**: Monitor goroutine count

### Health Metrics

#### `health_status`
- **Type**: Gauge
- **Description**: Health status of components (1 = healthy, 0 = unhealthy)
- **Labels**: `component`
- **Use Case**: Monitor component health

### MongoDB Metrics

#### `mongodb_connections_active`
- **Type**: Gauge
- **Description**: Number of active MongoDB connections
- **Labels**: `application`, `database`
- **Use Case**: Monitor connection pool utilization

#### `mongodb_connections_idle`
- **Type**: Gauge
- **Description**: Number of idle MongoDB connections
- **Labels**: `application`, `database`
- **Use Case**: Track available connections

#### `mongodb_connections_max`
- **Type**: Gauge
- **Description**: Maximum number of MongoDB connections allowed
- **Labels**: `application`, `database`
- **Use Case**: Monitor pool capacity configuration

#### `mongodb_connections_min`
- **Type**: Gauge
- **Description**: Minimum number of MongoDB connections maintained
- **Labels**: `application`, `database`
- **Use Case**: Track minimum pool size configuration

#### `mongodb_connections_total`
- **Type**: Gauge
- **Description**: Total number of MongoDB connections in the pool
- **Labels**: `application`, `database`
- **Use Case**: Monitor total connection pool size

#### `mongodb_connection_acquire_seconds`
- **Type**: Histogram
- **Description**: Time taken to acquire a MongoDB connection
- **Labels**: `application`, `database`
- **Use Case**: Monitor connection acquisition performance

#### `mongodb_connection_timeouts_total`
- **Type**: Counter
- **Description**: Total number of MongoDB connection acquisition timeouts
- **Labels**: `application`, `database`
- **Use Case**: **Critical metric**. Indicates connection pool exhaustion

#### `mongodb_operations_total`
- **Type**: Counter
- **Description**: Total number of MongoDB operations
- **Labels**: `application`, `database`, `operation`, `collection`
- **Use Case**: Track operation count by type (find, insert, update, delete, count)

#### `mongodb_operation_duration_seconds`
- **Type**: Histogram
- **Description**: Duration of MongoDB operations in seconds
- **Labels**: `application`, `database`, `operation`, `collection`
- **Use Case**: Monitor operation performance, calculate percentiles

#### `mongodb_operation_errors_total`
- **Type**: Counter
- **Description**: Total number of MongoDB operation errors
- **Labels**: `application`, `database`, `operation`, `collection`, `error_type`
- **Use Case**: Track error rates by operation type and error type

#### `mongodb_connection_errors_total`
- **Type**: Counter
- **Description**: Total number of MongoDB connection errors
- **Labels**: `application`, `database`, `error_type`
- **Use Case**: Track connection errors

#### `mongodb_ping_duration_seconds`
- **Type**: Histogram
- **Description**: MongoDB ping duration in seconds
- **Labels**: `application`, `database`
- **Use Case**: Monitor database connectivity latency

## Alert Rules

The project includes comprehensive alert rules in `deployments/prometheus/alerts.yml`:

### HTTP Alerts
- **HighErrorRate**: Error rate > 5% for 5 minutes
- **HighLatency**: p95 latency > 1s for 5 minutes
- **HighRequestRate**: Request rate spike > 2x baseline

### Database Alerts
- **DatabaseDown**: Database connection failure (service down, health check failing, or no active connections)
- **DatabaseHealthCheckFailure**: Health check reporting DOWN for 1 minute
- **MongoDBConnectionPoolExhaustion**: Connection pool usage > 90% for 5 minutes
- **MongoDBConnectionTimeouts**: Connection acquisition timeouts detected (> 0/sec for 2 minutes)
- **HighMongoDBOperationErrorRate**: Operation error rate > 5% for 5 minutes
- **HighMongoDBOperationLatency**: p95 operation latency > 1s for 5 minutes
- **MongoDBConnectionErrors**: Connection errors > 0.1/sec for 5 minutes

### System Alerts
- **HighMemoryUsage**: Memory usage > 85% for 5 minutes
- **CriticalMemoryUsage**: Memory usage > 95% for 2 minutes
- **HighGCPauseTime**: GC pause time > 0.1s per second for 5 minutes
- **HighGoroutineCount**: Goroutine count > 1000 for 5 minutes
- **HighFileDescriptors**: File descriptor usage > 90% for 5 minutes

### Business Metrics Alerts
- **HighBusinessOperationErrorRate**: Error rate > 10% for 5 minutes
- **ExternalServiceCallFailures**: External service failures > 0.1/sec for 5 minutes

### Service Alerts
- **ServiceUnavailable**: Service not responding (up metric == 0 for 1 minute)
- **HealthEndpointDown**: Health endpoint reporting issues for 2 minutes

## Testing

### Test API Endpoints

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","name":"Test User","status":"ACTIVE"}'

# Get all users
curl http://localhost:8080/api/users

# Get user by ID (replace {id} with actual MongoDB ObjectID)
curl http://localhost:8080/api/users/507f1f77bcf86cd799439011

# Update user (replace {id} with actual MongoDB ObjectID)
curl -X PUT http://localhost:8080/api/users/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Name"}'

# Delete user (replace {id} with actual MongoDB ObjectID)
curl -X DELETE http://localhost:8080/api/users/507f1f77bcf86cd799439011

# Health check
curl http://localhost:8080/health

# Metrics endpoint
curl http://localhost:8080/metrics
```

### Test Error Scenarios

```bash
# Trigger error endpoint (for testing alerts)
curl http://localhost:8080/api/users/test/error

# Trigger slow response (for testing latency alerts)
curl http://localhost:8080/api/users/test/slow
```

### Verify Metrics

```bash
# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Query metrics in Prometheus
# Visit http://localhost:9090 and try:
# - rate(http_server_requests_total[5m])
# - histogram_quantile(0.95, sum(rate(http_server_requests_seconds_bucket[5m])) by (le, uri, method))

# Check if alerts are loaded
curl http://localhost:9090/api/v1/rules

# View active alerts
curl http://localhost:9090/api/v1/alerts
```

## Prometheus Queries

### Calculate Error Rate
```promql
sum(rate(http_server_requests_total{status=~"5.."}[5m])) 
/ 
sum(rate(http_server_requests_total[5m])) * 100
```

### Calculate p95 Latency
```promql
histogram_quantile(0.95, 
  sum(rate(http_server_requests_seconds_bucket[5m])) by (le, uri, method)
)
```

### Memory Usage Percentage
```promql
(go_memstats_heap_inuse_bytes / go_memstats_heap_sys_bytes) * 100
```

### Request Rate by Endpoint
```promql
sum(rate(http_server_requests_total[5m])) by (uri, method)
```

### MongoDB Operation Rate
```promql
sum(rate(mongodb_operations_total[5m])) by (operation, collection)
```

### MongoDB Connection Pool Utilization
```promql
(mongodb_connections_active / mongodb_connections_max) * 100
```

## Development

### Building

```bash
# Create bin directory if it doesn't exist
mkdir -p bin

# Build the application
go build -o bin/server ./cmd/server

# Or use the build script (creates bin directory automatically)
./build.sh
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test package
go test ./internal/service/...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Formatting

```bash
# Format all Go files
go fmt ./...

# Check which files would be formatted (using gofmt directly)
gofmt -l .
```

### Linting

```bash
golangci-lint run
```

## Docker

### Build Docker Image

```bash
# Ensure go.sum exists (required by Dockerfile)
go mod tidy

# Build Docker image
docker build -t go-webapi-db .
```

**Note**: The Dockerfile requires `go.sum` to exist. Run `go mod tidy` first if you haven't already.

### Run Container

```bash
docker run -p 8080:8080 \
  -e SERVER_PORT=8080 \
  -e MONGODB_URI=mongodb://host.docker.internal:27017 \
  -e MONGODB_DATABASE=go_webapi_db \
  -e METRICS_PATH=/metrics \
  go-webapi-db
```

**Note**: If MongoDB is running in Docker, use `mongodb://mongodb:27017` instead of `host.docker.internal:27017` when containers are on the same network.

## Troubleshooting

### Application won't start
- Check MongoDB is running and accessible
- Verify environment variables are set correctly
- Check logs: `docker logs go-webapi-db`

### Metrics not appearing in Prometheus
- Verify Prometheus can reach the application: `curl http://go-webapi-db:8080/metrics` (from within Docker network) or `curl http://localhost:8080/metrics` (from host)
- Check Prometheus targets: `http://localhost:9090/targets`
- Verify network connectivity in Docker Compose
- Check Prometheus logs: `docker logs prometheus`
- Ensure the application is running: `docker ps | grep go-webapi-db`

### Alerts not firing
- Check alert rules are loaded: `http://localhost:9090/alerts` or `curl http://localhost:9090/api/v1/rules`
- Verify Alertmanager is configured correctly: `curl http://localhost:9093/api/v1/status`
- Check alert evaluation interval in Prometheus config (default: 15s)
- Verify metrics exist: Query the metrics used in alert expressions in Prometheus UI
- Check Prometheus logs for rule evaluation errors: `docker logs prometheus`

## License

This project is provided as-is for demonstration purposes.

## References

- [Gin Web Framework](https://gin-gonic.com/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/)
- [Prometheus Client Library](https://github.com/prometheus/client_golang)
- [Prometheus Documentation](https://prometheus.io/docs/)

