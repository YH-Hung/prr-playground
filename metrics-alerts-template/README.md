# Metrics and Alerts Template

A comprehensive observability template for web applications featuring Prometheus metrics collection, Grafana dashboards, and Alertmanager alerting. This template provides production-ready monitoring and alerting configurations for REST APIs with database connectivity across multiple technology stacks.

## Overview

This template includes reference implementations for different technology stacks, each demonstrating how to implement comprehensive observability:

- **Spring Boot** (`spring-webapi-db`) - Java/Spring Boot application with PostgreSQL
- **Go/Gin** (`go-webapi-db`) - Go application with MongoDB

Each implementation includes:
- **Prometheus** for metrics collection and alert rule evaluation
- **Grafana** with pre-configured dashboards and datasources
- **Alertmanager** for alert routing and notification management
- **Custom Business Metrics** for application-specific monitoring
- **Pre-configured Alerts** for common failure scenarios

## Project Structure

```
metrics-alerts-template/
├── spring-webapi-db/            # Spring Boot (Java) implementation
│   ├── src/                      # Application source code
│   ├── deployments/              # Prometheus, Grafana, Alertmanager configs
│   ├── docker-compose.metrics.yml
│   └── README.md                 # Spring Boot specific documentation
│
├── go-webapi-db/                 # Go/Gin implementation
│   ├── cmd/                      # Application entry points
│   ├── internal/                 # Application code
│   ├── deployments/              # Prometheus, Grafana, Alertmanager configs
│   ├── docker-compose.metrics.yml
│   └── README.md                 # Go specific documentation
│
└── README.md                     # This file (tech-agnostic overview)
```

## Features

### Metrics Collection

All implementations provide comprehensive metrics covering:

- **Application Metrics**: Startup time, readiness, custom business metrics
- **HTTP Metrics**: Request rate, latency (p50, p95, p99), error rates, status codes
- **Database Metrics**: Connection pool utilization, active/idle connections, operation metrics
- **Runtime Metrics**: Memory usage, CPU, garbage collection (language-specific)
- **System Metrics**: CPU usage, disk space, file descriptors
- **Custom Business Metrics**: Domain-specific operations, external service calls, error tracking

### Alerting

Pre-configured alert rules cover common failure scenarios:

- **High Error Rate**: Alert when error rate exceeds threshold
- **High Latency**: Alert when p95 latency exceeds threshold
- **Database Issues**: Connection failures, pool exhaustion, health check failures
- **Memory Pressure**: Memory usage warnings and critical alerts
- **Service Availability**: Service down, health endpoint failures
- **Resource Exhaustion**: Disk space, connection pool, thread/goroutine saturation
- **Business Metrics**: High business operation error rates, external service failures

### Dashboards

- Pre-configured Grafana dashboards for application monitoring
- Auto-provisioned datasources connecting to Prometheus
- Ready-to-use visualization templates

## Prerequisites

- **Docker and Docker Compose** (for running the monitoring stack)
- **Technology-specific prerequisites** (see individual project READMEs):
  - Spring Boot: Java 25+, Maven 3.6+, PostgreSQL
  - Go: Go 1.21+, MongoDB

## Quick Start

### 1. Choose Your Technology Stack

Select one of the reference implementations:

- **[Spring Boot](spring-webapi-db/README.md)** - Java/Spring Boot with PostgreSQL
- **[Go/Gin](go-webapi-db/README.md)** - Go with MongoDB

### 2. Start the Monitoring Stack

Each implementation includes a Docker Compose file to start the monitoring infrastructure:

```bash
cd <project-directory>
docker compose -f docker-compose.metrics.yml up -d
```

This will start:
- **Prometheus** on `http://localhost:9090`
- **Grafana** on `http://localhost:3000` (admin/admin)
- **Alertmanager** on `http://localhost:9093`
- **Database** (PostgreSQL or MongoDB depending on implementation)

### 3. Configure Prometheus Target

Update `deployments/prometheus/prometheus.yml` to point to your application:

```yaml
scrape_configs:
  - job_name: '<application-name>'
    metrics_path: '/metrics'  # or '/actuator/prometheus' for Spring Boot
    static_configs:
      - targets: ['host.docker.internal:8080']  # Update to your app's host:port
```

### 4. Build and Run Your Application

Follow the technology-specific instructions in each project's README:
- [Spring Boot Quick Start](spring-webapi-db/README.md#quick-start)
- [Go Quick Start](go-webapi-db/README.md#quick-start)

### 5. Verify Metrics Endpoint

Check that metrics are being exposed:

```bash
# For Spring Boot
curl http://localhost:8080/actuator/prometheus

# For Go
curl http://localhost:8080/metrics
```

### 6. Access Grafana

1. Open `http://localhost:3000`
2. Login with `admin` / `admin`
3. Navigate to **Dashboards** to view pre-configured dashboards
4. Navigate to **Explore** to query Prometheus metrics

## Configuration

### Prometheus Configuration

Each project includes a `deployments/prometheus/prometheus.yml` with:

- Global scrape and evaluation intervals
- Scrape configuration for the application
- Alert rule file references

### Alert Rules

Each project includes `deployments/prometheus/alerts.yml` with pre-configured alerts for:

- Service availability
- Error rates and latency
- Database connectivity
- Resource exhaustion
- Business metrics

### Alertmanager Configuration

Each project includes `deployments/alertmanager/alertmanager.yml` for configuring:

- Alert routing and grouping
- Notification channels (webhooks, email, Slack)
- Silence rules

### Grafana Provisioning

Each project includes Grafana provisioning configuration for:

- Prometheus datasource auto-configuration
- Dashboard auto-loading

## Common Metrics Patterns

While metric names may vary by technology stack, the following patterns are consistent:

### HTTP Request Metrics
- Request count (total, by status code, by endpoint)
- Request duration (histogram with p50, p95, p99)
- Error rates (4xx, 5xx)

### Database Metrics
- Connection pool size and utilization
- Active/idle connections
- Operation counts and duration
- Connection errors and timeouts

### Runtime Metrics
- Memory usage (heap/non-heap)
- CPU usage
- Garbage collection (language-specific)
- Thread/goroutine counts

### Business Metrics
- Domain operation counts
- Operation duration
- Error counts by type
- External service call metrics

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
3. Use the metrics documentation from each project as reference

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
- **Slack** - Slack webhook integration

## Troubleshooting

### Prometheus Not Scraping Metrics

1. Verify the application is running and accessible
2. Check the target in Prometheus UI: `http://localhost:9090/targets`
3. Verify the `metrics_path` is correct for your technology stack
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

1. Verify metrics endpoint is enabled and accessible
2. Check metrics endpoint path (varies by technology stack)
3. Test metrics endpoint: `curl http://localhost:8080/<metrics-path>`
4. Review application logs for metric collection errors

For technology-specific troubleshooting, see:
- [Spring Boot Troubleshooting](spring-webapi-db/README.md#troubleshooting)
- [Go Troubleshooting](go-webapi-db/README.md#troubleshooting)

## Production Considerations

### Security
- Enable authentication for Prometheus, Grafana, and Alertmanager
- Use HTTPS for all services
- Restrict network access to monitoring endpoints
- Use secrets management for credentials

### Scalability
- Consider Prometheus federation for multiple instances
- Use remote storage (e.g., Thanos, Cortex) for long-term retention
- Implement alert deduplication and grouping
- Use service discovery for dynamic targets

### High Availability
- Run multiple Prometheus instances for redundancy
- Use Alertmanager clustering for high availability
- Implement Grafana high availability setup
- Backup Grafana dashboards and datasources

### Performance
- Adjust scrape intervals based on metrics volume
- Configure retention policies for Prometheus
- Optimize alert rule evaluation intervals
- Use recording rules for expensive queries

## Project-Specific Documentation

- **[Spring Boot Implementation](spring-webapi-db/README.md)** - Java/Spring Boot specific documentation
- **[Go Implementation](go-webapi-db/README.md)** - Go/Gin specific documentation

## External Documentation

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Alertmanager Documentation](https://prometheus.io/docs/alerting/latest/alertmanager/)

## License

This template is provided as-is for use in your projects.
