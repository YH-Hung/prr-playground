# Prometheus Metrics Documentation

This document provides a comprehensive overview of all metrics exposed by the Spring Boot application at the `/actuator/prometheus` endpoint.

## Table of Contents

1. [Application Metrics](#application-metrics)
2. [Custom Business Metrics](#custom-business-metrics)
3. [HTTP Server Metrics](#http-server-metrics)
4. [Database Connection Pool Metrics (HikariCP)](#database-connection-pool-metrics-hikaricp)
5. [JDBC Metrics](#jdbc-metrics)
6. [JVM Metrics](#jvm-metrics)
7. [System Metrics](#system-metrics)
8. [Process Metrics](#process-metrics)
9. [Executor Metrics](#executor-metrics)
10. [Tomcat Session Metrics](#tomcat-session-metrics)
11. [Logging Metrics](#logging-metrics)
12. [Disk Metrics](#disk-metrics)

---

## Application Metrics

### `application_ready_time_seconds`
- **Type**: Gauge
- **Description**: Time taken for the application to be ready to service requests (in seconds)
- **Tags**: 
  - `application`: Application name (spring-webapi-db)
  - `main_application_class`: Main application class name
- **Use Case**: Monitor application startup performance. High values indicate slow startup.

### `application_started_time_seconds`
- **Type**: Gauge
- **Description**: Time taken to start the application (in seconds)
- **Tags**: 
  - `application`: Application name
  - `main_application_class`: Main application class name
- **Use Case**: Track total application startup time. Useful for performance optimization.

---

## Custom Business Metrics

### `custom_user_total`
- **Type**: Counter
- **Description**: Total number of users created since application start
- **Tags**: 
  - `application`: Application name
  - `operation`: Operation type (create)
- **Use Case**: Track user creation rate and total users created. Useful for business analytics.

### `custom_user_updated_total`
- **Type**: Counter
- **Description**: Total number of users updated since application start
- **Tags**: 
  - `application`: Application name
  - `operation`: Operation type (update)
- **Use Case**: Monitor user update operations and track modification patterns.

### `custom_user_deleted_total`
- **Type**: Counter
- **Description**: Total number of users deleted since application start
- **Tags**: 
  - `application`: Application name
  - `operation`: Operation type (delete)
- **Use Case**: Track user deletion rate. Monitor for unusual deletion patterns.

### `custom_user_active_operations`
- **Type**: Gauge
- **Description**: Current number of active user operations being processed
- **Tags**: 
  - `application`: Application name
- **Use Case**: Monitor concurrent user operations. High values may indicate performance issues or bottlenecks.

### `custom_user_operation_duration_seconds`
- **Type**: Summary
- **Description**: Duration of user operations (create, read, update, delete)
- **Tags**: 
  - `application`: Application name
- **Metrics Provided**:
  - `_count`: Total number of operations
  - `_sum`: Sum of all operation durations
  - `_max`: Maximum operation duration
- **Use Case**: Monitor operation performance. Calculate percentiles (p50, p95, p99) for latency analysis.

### `custom_user_operation_errors_total`
- **Type**: Counter
- **Description**: Total number of user operation errors
- **Tags**: 
  - `application`: Application name
  - `error.type`: Type of error (e.g., duplicate_email, not_found, timeout)
- **Use Case**: Track error rates by type. Essential for alerting on high error rates.

### `custom_external_call_duration_seconds`
- **Type**: Summary
- **Description**: Duration of external service calls
- **Tags**: 
  - `application`: Application name
  - `service`: Name of the external service
- **Use Case**: Monitor external service call performance and detect slow external dependencies.

### `custom_external_call_errors_total`
- **Type**: Counter
- **Description**: Total number of external service call failures
- **Tags**: 
  - `application`: Application name
  - `service`: Name of the external service
- **Use Case**: Track external service reliability. Alert when external services are failing.

---

## HTTP Server Metrics

### `http_server_requests_seconds`
- **Type**: Histogram
- **Description**: HTTP request duration in seconds, broken down by method, status, and URI
- **Tags**: 
  - `application`: Application name
  - `method`: HTTP method (GET, POST, PUT, DELETE, etc.)
  - `status`: HTTP status code (200, 404, 500, etc.)
  - `uri`: Request URI (normalized)
  - `outcome`: Request outcome (SUCCESS, CLIENT_ERROR, SERVER_ERROR)
  - `exception`: Exception type if any (none otherwise)
  - `error`: Error type if any
- **Metrics Provided**:
  - `_bucket`: Histogram buckets for percentile calculations
  - `_count`: Total number of requests
  - `_sum`: Sum of all request durations
  - `_max`: Maximum request duration
- **Use Case**: 
  - Calculate request latency percentiles (p50, p95, p99)
  - Monitor request rates by endpoint
  - Track error rates by status code
  - Identify slow endpoints

### `http_server_requests_active_seconds`
- **Type**: Histogram
- **Description**: Time spent processing active HTTP requests
- **Tags**: Same as `http_server_requests_seconds`
- **Use Case**: Monitor concurrent request processing time.

### `http_server_requests_total`
- **Type**: Counter
- **Description**: Total number of HTTP requests (from custom interceptor)
- **Tags**: 
  - `method`: HTTP method
  - `uri`: Request URI
  - `status`: HTTP status code
- **Use Case**: Track total request count per endpoint.

### `http_server_errors_total`
- **Type**: Counter
- **Description**: Total number of HTTP server errors (5xx status codes)
- **Tags**: 
  - `method`: HTTP method
  - `uri`: Request URI
  - `status`: HTTP status code
- **Use Case**: Monitor server error rate. Critical for alerting.

### `http_server_client_errors_total`
- **Type**: Counter
- **Description**: Total number of HTTP client errors (4xx status codes)
- **Tags**: 
  - `method`: HTTP method
  - `uri`: Request URI
  - `status`: HTTP status code
- **Use Case**: Track client error patterns. May indicate API usage issues.

---

## Database Connection Pool Metrics (HikariCP)

### `hikaricp_connections`
- **Type**: Gauge
- **Description**: Total number of connections in the pool
- **Tags**: 
  - `application`: Application name
  - `pool`: Connection pool name (HikariPool-1)
- **Use Case**: Monitor total connection pool size.

### `hikaricp_connections_active`
- **Type**: Gauge
- **Description**: Number of active (in-use) connections
- **Tags**: Same as above
- **Use Case**: Monitor connection pool utilization. High values may indicate connection leaks or high load.

### `hikaricp_connections_idle`
- **Type**: Gauge
- **Description**: Number of idle connections available for use
- **Tags**: Same as above
- **Use Case**: Track available connections. Low values may indicate connection pool exhaustion.

### `hikaricp_connections_max`
- **Type**: Gauge
- **Description**: Maximum number of connections allowed in the pool
- **Tags**: Same as above
- **Use Case**: Monitor pool capacity configuration.

### `hikaricp_connections_min`
- **Type**: Gauge
- **Description**: Minimum number of idle connections maintained
- **Tags**: Same as above
- **Use Case**: Track minimum pool size configuration.

### `hikaricp_connections_pending`
- **Type**: Gauge
- **Description**: Number of threads waiting for a connection
- **Tags**: Same as above
- **Use Case**: **Critical metric**. Non-zero values indicate connection pool exhaustion. Alert immediately.

### `hikaricp_connections_timeout_total`
- **Type**: Counter
- **Description**: Total number of connection acquisition timeouts
- **Tags**: Same as above
- **Use Case**: **Critical metric**. Indicates connection pool is exhausted. Alert on any non-zero value.

### `hikaricp_connections_acquire_seconds`
- **Type**: Summary
- **Description**: Time taken to acquire a connection from the pool
- **Tags**: Same as above
- **Metrics Provided**: `_count`, `_sum`, `_max`
- **Use Case**: Monitor connection acquisition performance. High values indicate pool contention.

### `hikaricp_connections_creation_seconds`
- **Type**: Summary
- **Description**: Time taken to create new database connections
- **Tags**: Same as above
- **Metrics Provided**: `_count`, `_sum`, `_max`
- **Use Case**: Monitor database connection establishment performance.

### `hikaricp_connections_usage_seconds`
- **Type**: Summary
- **Description**: Time connections are held/used
- **Tags**: Same as above
- **Metrics Provided**: `_count`, `_sum`, `_max`
- **Use Case**: Monitor how long connections are held. High values may indicate slow queries or connection leaks.

---

## JDBC Metrics

### `jdbc_connections_active`
- **Type**: Gauge
- **Description**: Current number of active JDBC connections
- **Tags**: 
  - `application`: Application name
  - `name`: DataSource name (dataSource)
- **Use Case**: Monitor active database connections.

### `jdbc_connections_idle`
- **Type**: Gauge
- **Description**: Number of idle JDBC connections
- **Tags**: Same as above
- **Use Case**: Track available database connections.

### `jdbc_connections_max`
- **Type**: Gauge
- **Description**: Maximum number of JDBC connections allowed
- **Tags**: Same as above
- **Use Case**: Monitor connection pool maximum capacity.

### `jdbc_connections_min`
- **Type**: Gauge
- **Description**: Minimum number of idle JDBC connections
- **Tags**: Same as above
- **Use Case**: Track minimum connection pool size.

---

## JVM Metrics

### `jvm_memory_used_bytes`
- **Type**: Gauge
- **Description**: Amount of memory currently used by the JVM
- **Tags**: 
  - `application`: Application name
  - `area`: Memory area (heap, nonheap)
  - `id`: Memory pool ID (e.g., G1 Eden Space, G1 Old Gen, Metaspace, CodeCache)
- **Use Case**: 
  - Monitor memory usage by area
  - Detect memory leaks
  - Track heap vs non-heap memory usage
  - **Critical for alerting**: Alert when heap usage exceeds 85%

### `jvm_memory_max_bytes`
- **Type**: Gauge
- **Description**: Maximum amount of memory available for the JVM
- **Tags**: Same as above
- **Use Case**: Calculate memory usage percentage: `used / max * 100`

### `jvm_memory_committed_bytes`
- **Type**: Gauge
- **Description**: Amount of memory committed (guaranteed to be available) by the JVM
- **Tags**: Same as above
- **Use Case**: Monitor memory allocation. Committed memory is reserved by the OS.

### `jvm_gc_pause_seconds`
- **Type**: Summary
- **Description**: Time spent in garbage collection pauses
- **Tags**: 
  - `application`: Application name
  - `action`: GC action (end_of_minor_gc, end_of_major_gc)
  - `cause`: GC cause
- **Metrics Provided**: `_count`, `_sum`, `_max`
- **Use Case**: 
  - Monitor GC performance
  - Detect excessive GC pauses
  - **Alert**: High GC pause time indicates memory pressure

### `jvm_gc_live_data_size_bytes`
- **Type**: Gauge
- **Description**: Size of long-lived heap memory after GC
- **Tags**: 
  - `application`: Application name
- **Use Case**: Monitor long-lived object memory usage.

### `jvm_gc_max_data_size_bytes`
- **Type**: Gauge
- **Description**: Maximum size of long-lived heap memory pool
- **Tags**: Same as above
- **Use Case**: Track maximum old generation heap size.

### `jvm_gc_memory_allocated_bytes_total`
- **Type**: Counter
- **Description**: Total memory allocated in young generation between GCs
- **Tags**: Same as above
- **Use Case**: Monitor memory allocation rate.

### `jvm_gc_memory_promoted_bytes_total`
- **Type**: Counter
- **Description**: Total memory promoted from young to old generation
- **Tags**: Same as above
- **Use Case**: Track object promotion rate. High values may indicate memory pressure.

### `jvm_gc_overhead`
- **Type**: Gauge
- **Description**: Percentage of CPU time used by GC (0-1 range)
- **Tags**: Same as above
- **Use Case**: **Critical metric**. High values (>0.1) indicate excessive GC overhead. Alert threshold.

### `jvm_memory_usage_after_gc`
- **Type**: Gauge
- **Description**: Percentage of long-lived heap used after last GC (0-1 range)
- **Tags**: 
  - `application`: Application name
  - `area`: Memory area (heap)
  - `pool`: Memory pool (long-lived)
- **Use Case**: Monitor heap usage after GC. High values indicate memory pressure.

### `jvm_classes_loaded_classes`
- **Type**: Gauge
- **Description**: Current number of classes loaded in the JVM
- **Tags**: 
  - `application`: Application name
- **Use Case**: Monitor class loading. Unusual increases may indicate class loader leaks.

### `jvm_classes_loaded_count_classes_total`
- **Type**: Counter
- **Description**: Total number of classes loaded since JVM start
- **Tags**: Same as above
- **Use Case**: Track class loading rate.

### `jvm_classes_unloaded_classes_total`
- **Type**: Counter
- **Description**: Total number of classes unloaded since JVM start
- **Tags**: Same as above
- **Use Case**: Monitor class unloading (rare in most applications).

### `jvm_compilation_time_ms_total`
- **Type**: Counter
- **Description**: Total time spent in JIT compilation (milliseconds)
- **Tags**: 
  - `application`: Application name
  - `compiler`: Compiler name
- **Use Case**: Monitor JIT compilation overhead.

### `jvm_buffer_count_buffers`
- **Type**: Gauge
- **Description**: Number of buffers in buffer pools
- **Tags**: 
  - `application`: Application name
  - `id`: Buffer pool ID (direct, mapped)
- **Use Case**: Monitor direct memory usage.

### `jvm_buffer_memory_used_bytes`
- **Type**: Gauge
- **Description**: Memory used by buffer pools
- **Tags**: Same as above
- **Use Case**: Track direct memory consumption.

### `jvm_buffer_total_capacity_bytes`
- **Type**: Gauge
- **Description**: Total capacity of buffer pools
- **Tags**: Same as above
- **Use Case**: Monitor buffer pool capacity.

### `jvm_info`
- **Type**: Gauge
- **Description**: JVM version information
- **Tags**: 
  - `application`: Application name
  - `version`: JVM version
  - `vendor`: JVM vendor
  - `runtime`: Runtime name
- **Use Case**: Track JVM version and vendor information.

---

## System Metrics

### `system_cpu_count`
- **Type**: Gauge
- **Description**: Number of CPU cores available to the JVM
- **Tags**: 
  - `application`: Application name
- **Use Case**: Reference metric for CPU usage calculations.

### `system_cpu_usage`
- **Type**: Gauge
- **Description**: Recent CPU usage of the entire system (0-1 range)
- **Tags**: Same as above
- **Use Case**: Monitor system-wide CPU usage. High values may indicate system resource constraints.

### `system_load_average_1m`
- **Type**: Gauge
- **Description**: System load average over 1 minute
- **Tags**: Same as above
- **Use Case**: Monitor system load. Values > CPU count indicate system overload.

---

## Process Metrics

### `process_cpu_usage`
- **Type**: Gauge
- **Description**: Recent CPU usage of the Java process (0-1 range)
- **Tags**: 
  - `application`: Application name
- **Use Case**: Monitor application CPU usage. High values indicate CPU-intensive operations.

### `process_cpu_time_ns_total`
- **Type**: Counter
- **Description**: Total CPU time used by the Java process (nanoseconds)
- **Tags**: Same as above
- **Use Case**: Track cumulative CPU time. Useful for calculating average CPU usage over time.

### `process_uptime_seconds`
- **Type**: Gauge
- **Description**: JVM uptime in seconds
- **Tags**: Same as above
- **Use Case**: Monitor application uptime. Track for availability metrics.

### `process_start_time_seconds`
- **Type**: Gauge
- **Description**: Process start time (Unix epoch timestamp)
- **Tags**: Same as above
- **Use Case**: Calculate uptime and track restart events.

### `process_files_open_files`
- **Type**: Gauge
- **Description**: Number of open file descriptors
- **Tags**: Same as above
- **Use Case**: Monitor file descriptor usage. High values may indicate file handle leaks.

### `process_files_max_files`
- **Type**: Gauge
- **Description**: Maximum number of file descriptors allowed
- **Tags**: Same as above
- **Use Case**: Track file descriptor limits. Alert when approaching maximum.

---

## Executor Metrics

### `executor_active_threads`
- **Type**: Gauge
- **Description**: Number of threads actively executing tasks
- **Tags**: 
  - `application`: Application name
  - `name`: Executor name (applicationTaskExecutor)
- **Use Case**: Monitor thread pool utilization.

### `executor_pool_size_threads`
- **Type**: Gauge
- **Description**: Current number of threads in the pool
- **Tags**: Same as above
- **Use Case**: Track thread pool size.

### `executor_pool_core_threads`
- **Type**: Gauge
- **Description**: Core number of threads in the pool
- **Tags**: Same as above
- **Use Case**: Monitor core thread pool configuration.

### `executor_pool_max_threads`
- **Type**: Gauge
- **Description**: Maximum number of threads allowed in the pool
- **Tags**: Same as above
- **Use Case**: Track maximum thread pool size.

### `executor_queued_tasks`
- **Type**: Gauge
- **Description**: Number of tasks queued for execution
- **Tags**: Same as above
- **Use Case**: **Critical metric**. Non-zero values indicate thread pool saturation. Alert threshold.

### `executor_queue_remaining_tasks`
- **Type**: Gauge
- **Description**: Remaining queue capacity
- **Tags**: Same as above
- **Use Case**: Monitor queue capacity. Low values indicate potential saturation.

### `executor_completed_tasks_total`
- **Type**: Counter
- **Description**: Total number of completed tasks
- **Tags**: Same as above
- **Use Case**: Track task completion rate.

---

## Tomcat Session Metrics

### `tomcat_sessions_active_current_sessions`
- **Type**: Gauge
- **Description**: Current number of active HTTP sessions
- **Tags**: 
  - `application`: Application name
- **Use Case**: Monitor session count. Useful for session management and memory planning.

### `tomcat_sessions_active_max_sessions`
- **Type**: Gauge
- **Description**: Maximum number of concurrent sessions
- **Tags**: Same as above
- **Use Case**: Track peak session count.

### `tomcat_sessions_alive_max_seconds`
- **Type**: Gauge
- **Description**: Maximum session alive time in seconds
- **Tags**: Same as above
- **Use Case**: Monitor longest-lived session.

### `tomcat_sessions_created_sessions_total`
- **Type**: Counter
- **Description**: Total number of sessions created
- **Tags**: Same as above
- **Use Case**: Track session creation rate.

### `tomcat_sessions_expired_sessions_total`
- **Type**: Counter
- **Description**: Total number of expired sessions
- **Tags**: Same as above
- **Use Case**: Monitor session expiration rate.

### `tomcat_sessions_rejected_sessions_total`
- **Type**: Counter
- **Description**: Total number of rejected sessions (when max sessions reached)
- **Tags**: Same as above
- **Use Case**: **Critical metric**. Non-zero values indicate session limit reached. Alert immediately.

---

## Logging Metrics

### `logback_events_total`
- **Type**: Counter
- **Description**: Total number of log events by level
- **Tags**: 
  - `application`: Application name
  - `level`: Log level (trace, debug, info, warn, error)
- **Use Case**: 
  - Monitor log volume by level
  - Track error log rate
  - **Alert**: High error log count indicates application issues

---

## Disk Metrics

### `disk_free_bytes`
- **Type**: Gauge
- **Description**: Free disk space in bytes
- **Tags**: 
  - `application`: Application name
  - `path`: Disk path being monitored
- **Use Case**: 
  - Monitor available disk space
  - **Critical for alerting**: Alert when free space < 10% of total

### `disk_total_bytes`
- **Type**: Gauge
- **Description**: Total disk space in bytes
- **Tags**: Same as above
- **Use Case**: Calculate disk usage percentage: `(total - free) / total * 100`

---

## Key Metrics for Alerting

### Critical Metrics (Alert Immediately)
1. **`hikaricp_connections_pending`** > 0 - Connection pool exhausted
2. **`hikaricp_connections_timeout_total`** > 0 - Connection acquisition failures
3. **`jvm_memory_used_bytes / jvm_memory_max_bytes`** > 0.95 - Critical memory pressure
4. **`tomcat_sessions_rejected_sessions_total`** > 0 - Session limit reached
5. **`executor_queued_tasks`** > 0 - Thread pool saturation

### Warning Metrics (Monitor Closely)
1. **`jvm_memory_used_bytes / jvm_memory_max_bytes`** > 0.85 - High memory usage
2. **`jvm_gc_overhead`** > 0.1 - Excessive GC overhead
3. **`http_server_requests_seconds`** p95 > 1s - High latency
4. **Error rate** > 5% - High error rate
5. **`disk_free_bytes / disk_total_bytes`** < 0.1 - Low disk space
6. **`hikaricp_connections_active / hikaricp_connections_max`** > 0.9 - Connection pool near exhaustion

### Business Metrics to Track
1. **`custom_user_operation_errors_total`** - Business operation error rate
2. **`custom_external_call_errors_total`** - External service reliability
3. **`custom_user_operation_duration_seconds`** - Business operation performance
4. **Request rate spikes** - Sudden traffic increases

---

## Accessing Metrics

### Prometheus Endpoint
```
GET http://localhost:8080/actuator/prometheus
```

### Metrics Endpoint (JSON)
```
GET http://localhost:8080/actuator/metrics
```

### Individual Metric
```
GET http://localhost:8080/actuator/metrics/{metricName}
```

### Health Endpoint
```
GET http://localhost:8080/actuator/health
```

---

## Metric Naming Conventions

- **Counters**: End with `_total` (e.g., `custom_user_total`)
- **Gauges**: No suffix (e.g., `jvm_memory_used_bytes`)
- **Histograms**: End with `_seconds` or `_bytes` (e.g., `http_server_requests_seconds`)
- **Summaries**: End with `_seconds` with `_count`, `_sum`, `_max` variants

All metrics are tagged with `application="spring-webapi-db"` for filtering and aggregation.

---

## Example Prometheus Queries

### Calculate Error Rate
```promql
sum(rate(http_server_requests_seconds_count{status=~"5.."}[5m])) 
/ 
sum(rate(http_server_requests_seconds_count[5m])) * 100
```

### Calculate p95 Latency
```promql
histogram_quantile(0.95, 
  sum(rate(http_server_requests_seconds_bucket[5m])) by (le, uri, method)
)
```

### Memory Usage Percentage
```promql
(jvm_memory_used_bytes{area="heap"} / jvm_memory_max_bytes{area="heap"}) * 100
```

### Connection Pool Utilization
```promql
(hikaricp_connections_active / hikaricp_connections_max) * 100
```

### Request Rate by Endpoint
```promql
sum(rate(http_server_requests_seconds_count[5m])) by (uri, method)
```

---

## Notes

- All time-based metrics are in seconds unless otherwise specified
- All byte-based metrics are in bytes
- Histogram buckets allow percentile calculations (p50, p95, p99)
- Counters are cumulative and reset on application restart
- Gauges represent current state and can increase or decrease
- Summary metrics provide count, sum, and max for duration calculations

