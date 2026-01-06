# Spring Boot WebAPI Database Metrics Project

A production-ready Spring Boot REST API with PostgreSQL database connectivity, featuring comprehensive Prometheus metrics collection, Grafana dashboards, and Alertmanager alerting capabilities.

## Features

### Metrics Collection
- **Application Metrics**: Startup time, readiness, custom business metrics
- **HTTP Metrics**: Request rate, latency (p50, p95, p99), error rates, status codes
- **Database Metrics**: Connection pool utilization, active/idle connections, timeouts (HikariCP)
- **JVM Metrics**: Memory usage, GC pauses, thread counts, class loading
- **System Metrics**: CPU usage, disk space, file descriptors
- **Custom Business Metrics**: User operations, external service calls, error tracking

### Alerting
- **High Error Rate**: Alert when error rate exceeds 5%
- **High Latency**: Alert when p95 latency exceeds 1 second
- **Database Issues**: Connection failures, pool exhaustion, health check failures
- **Memory Pressure**: JVM heap usage warnings (85%) and critical alerts (95%)
- **Service Availability**: Service down, health endpoint failures
- **Resource Exhaustion**: Disk space, connection pool, thread pool saturation
- **Business Metrics**: High business operation error rates, external service failures

### Dashboards
- Pre-configured Grafana dashboards for application monitoring
- Auto-provisioned datasources connecting to Prometheus
- Ready-to-use visualization templates

## Prerequisites

- **Java 25+** (or compatible JDK)
- **Maven 3.6+**
- **Docker and Docker Compose**
- **PostgreSQL** (for the Spring Boot application database)

## Project Structure

```
spring-webapi-db/
├── src/                          # Spring Boot application source code
│   ├── main/
│   │   ├── java/
│   │   │   └── org/hle/springwebapidb/
│   │   │       ├── config/       # Metrics and web configuration
│   │   │       ├── controller/   # REST controllers
│   │   │       ├── service/      # Business logic and metrics
│   │   │       ├── repository/   # JPA repositories
│   │   │       ├── entity/       # JPA entities
│   │   │       ├── interceptor/  # HTTP metrics interceptor
│   │   │       └── exception/    # Global exception handler
│   │   └── resources/
│   │       └── application.properties
│   └── test/                     # Unit and integration tests
├── deployments/                  # Deployment configurations
│   ├── prometheus/
│   │   ├── prometheus.yml        # Prometheus configuration
│   │   └── alerts.yml            # Alert rules
│   ├── alertmanager/
│   │   └── alertmanager.yml      # Alert routing and notifications
│   └── grafana/
│       └── provisioning/         # Auto-provisioned dashboards and datasources
│           ├── datasources/
│           │   └── prometheus.yml
│           └── dashboards/
│               └── dashboard.yml
├── docker-compose.metrics.yml    # Docker Compose for monitoring stack
├── METRICS_DOCUMENTATION.md      # Comprehensive metrics documentation
├── HELP.md                       # Spring Boot getting started guide
└── pom.xml                       # Maven project configuration
```

## Quick Start

### 1. Start the Monitoring Stack

Start Prometheus, Grafana, and Alertmanager:

```bash
cd spring-webapi-db
docker compose -f docker-compose.metrics.yml up -d
```

This will start:
- **Prometheus** on `http://localhost:9090`
- **Grafana** on `http://localhost:3000` (admin/admin)
- **Alertmanager** on `http://localhost:9093`

### 2. Configure Prometheus Target

The Spring Boot application needs to be accessible to Prometheus. Update `deployments/prometheus/prometheus.yml` to point to your application:

```yaml
scrape_configs:
  - job_name: 'spring-webapi-db'
    static_configs:
      - targets: ['host.docker.internal:8080']  # Update this to your app's host:port
```

If running the application locally, use `host.docker.internal:8080`. For Docker networks, use the service name and port.

### 3. Build and Run the Spring Boot Application

```bash
# Build the application
mvn clean package

# Run the application
java -jar target/spring-webapi-db-0.0.1-SNAPSHOT.jar
```

Or use Maven:

```bash
mvn spring-boot:run
```

### 4. Verify Metrics Endpoint

Check that metrics are being exposed:

```bash
curl http://localhost:8080/actuator/prometheus
```

### 5. Access Grafana

1. Open `http://localhost:3000`
2. Login with `admin` / `admin`
3. Navigate to **Dashboards** to view pre-configured dashboards
4. Navigate to **Explore** to query Prometheus metrics

## Configuration

### Application Configuration

Edit `src/main/resources/application.properties`:

```properties
# Application name
spring.application.name=spring-webapi-db

# Database configuration
spring.datasource.url=jdbc:postgresql://localhost:5432/postgres
spring.datasource.username=postgres
spring.datasource.password=postgres

# Actuator endpoints
management.endpoints.web.exposure.include=health,metrics,prometheus,info
management.endpoint.health.show-details=always
management.prometheus.metrics.export.enabled=true

# Metrics tags
management.metrics.tags.application=${spring.application.name}
```

### Prometheus Configuration

Edit `deployments/prometheus/prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'spring-webapi-db'
    metrics_path: '/actuator/prometheus'
    static_configs:
      - targets: ['host.docker.internal:8080']
```

### Alert Rules

Edit `deployments/prometheus/alerts.yml` to customize alert thresholds:

```yaml
groups:
  - name: spring_webapi_db_alerts
    rules:
      - alert: HighErrorRate
        expr: |
          (sum(rate(http_server_requests_total{status=~"5.."}[5m])) 
           / sum(rate(http_server_requests_total[5m]))) * 100 > 5
        for: 5m
        annotations:
          summary: "High error rate detected"
```

### Alertmanager Configuration

Edit `deployments/alertmanager/alertmanager.yml` to configure notification channels:

```yaml
receivers:
  - name: 'critical-receiver'
    webhook_configs:
      - url: 'http://your-webhook-url/webhook'
    # email_configs:
    #   - to: 'admin@example.com'
    #     from: 'alertmanager@example.com'
    #     smarthost: 'smtp.example.com:587'
```

## Available Metrics

### Application Metrics
- `application_ready_time_seconds` - Application readiness time
- `application_started_time_seconds` - Application startup time

### Custom Business Metrics
- `custom_user_total` - Total users created
- `custom_user_updated_total` - Total users updated
- `custom_user_deleted_total` - Total users deleted
- `custom_user_active_operations` - Current active operations
- `custom_user_operation_duration_seconds` - Operation duration (p50, p95, p99)
- `custom_user_operation_errors_total` - Operation errors by type
- `custom_external_call_duration_seconds` - External service call duration
- `custom_external_call_errors_total` - External service call failures

### HTTP Metrics
- `http_server_requests_seconds` - Request duration histogram
- `http_server_requests_total` - Total request count
- `http_server_errors_total` - Server error count (5xx)
- `http_server_client_errors_total` - Client error count (4xx)

### Database Metrics (HikariCP)
- `hikaricp_connections` - Total connections in pool
- `hikaricp_connections_active` - Active connections
- `hikaricp_connections_idle` - Idle connections
- `hikaricp_connections_pending` - Threads waiting for connections
- `hikaricp_connections_timeout_total` - Connection acquisition timeouts

### JVM Metrics
- `jvm_memory_used_bytes` - Memory usage by area
- `jvm_memory_max_bytes` - Maximum memory available
- `jvm_gc_pause_seconds` - GC pause duration
- `jvm_gc_overhead` - GC overhead percentage

### System Metrics
- `system_cpu_usage` - System CPU usage
- `process_cpu_usage` - Process CPU usage
- `disk_free_bytes` - Free disk space
- `disk_total_bytes` - Total disk space

For a complete list of all metrics, see [METRICS_DOCUMENTATION.md](METRICS_DOCUMENTATION.md).

## Alert Rules

The template includes pre-configured alerts for:

### Critical Alerts
- **ServiceUnavailable** - Application not responding
- **DatabaseDown** - Database connection failure
- **CriticalMemoryUsage** - JVM heap usage > 95%
- **CriticalDiskSpace** - Disk space < 5%

### Warning Alerts
- **HighErrorRate** - Error rate > 5%
- **HighLatency** - p95 latency > 1 second
- **HighMemoryUsage** - JVM heap usage > 85%
- **ConnectionPoolExhaustion** - Connection pool usage > 90%
- **HighRequestRate** - Request rate spike > 2x baseline
- **HighBusinessOperationErrorRate** - Business error rate > 10%

View all alert rules in `deployments/prometheus/alerts.yml`.

## Example Prometheus Queries

### Error Rate
```promql
sum(rate(http_server_requests_total{status=~"5.."}[5m])) 
/ 
sum(rate(http_server_requests_total[5m])) * 100
```

### p95 Latency
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
sum(rate(http_server_requests_total[5m])) by (uri, method)
```

## Grafana Dashboards

### Accessing Dashboards

1. Navigate to `http://localhost:3000`
2. Login with `admin` / `admin`
3. Go to **Dashboards** → **Browse**
4. Select a dashboard from the list

### Creating Custom Dashboards

1. Go to **Dashboards** → **New** → **New Dashboard**
2. Add panels with Prometheus queries
3. Use the metrics from `METRICS_DOCUMENTATION.md` as reference

## Alertmanager

### Viewing Alerts

Access Alertmanager UI at `http://localhost:9093` to:
- View active alerts
- Silence alerts
- View alert history
- Test notification routes

### Notification Channels

Configure notification channels in `deployments/alertmanager/alertmanager.yml`:
- **Webhooks** - HTTP POST to custom endpoints
- **Email** - SMTP email notifications
- **Slack** - Slack webhook integration (configure `slack_api_url`)

## Development

### Building the Application

```bash
mvn clean package
```

### Running Tests

```bash
mvn test
```

### Running Locally

```bash
mvn spring-boot:run
```

### Accessing Actuator Endpoints

- Health: `http://localhost:8080/actuator/health`
- Metrics: `http://localhost:8080/actuator/metrics`
- Prometheus: `http://localhost:8080/actuator/prometheus`
- Info: `http://localhost:8080/actuator/info`

## Troubleshooting

### Prometheus Not Scraping Metrics

1. Verify the application is running and accessible
2. Check the target in Prometheus UI: `http://localhost:9090/targets`
3. Verify the `metrics_path` is correct (`/actuator/prometheus`)
4. Check network connectivity between Prometheus and the application
5. Review Prometheus logs: `docker compose -f docker-compose.metrics.yml logs prometheus`

### Alerts Not Firing

1. Verify alert rules are loaded: `http://localhost:9090/alerts`
2. Check Prometheus query evaluation: Use the Prometheus query UI
3. Verify Alertmanager is connected: Check `http://localhost:9090/status`
4. Review Alertmanager configuration: `http://localhost:9093/#/status`

### Grafana Not Showing Data

1. Verify Prometheus datasource is configured correctly
2. Check datasource connection: Test in Grafana datasource settings
3. Verify metrics exist: Query Prometheus directly
4. Check time range: Ensure queries match available data

### Application Metrics Not Appearing

1. Verify Actuator is enabled: Check `application.properties`
2. Verify Prometheus export is enabled: `management.prometheus.metrics.export.enabled=true`
3. Check endpoint exposure: `management.endpoints.web.exposure.include=prometheus`
4. Test metrics endpoint: `curl http://localhost:8080/actuator/prometheus`

## Documentation

- [METRICS_DOCUMENTATION.md](METRICS_DOCUMENTATION.md) - Complete metrics reference
- [HELP.md](HELP.md) - Spring Boot getting started guide
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Alertmanager Documentation](https://prometheus.io/docs/alerting/latest/alertmanager/)

## License

This template is provided as-is for use in your projects.

