# PRR Playground

A collection of observability and monitoring examples demonstrating best practices for log aggregation, metrics collection, and alerting.

## Projects

### [log-parsing-forwarding](./log-parsing-forwarding/)

A Go-based example demonstrating log aggregation with Vector and Elasticsearch. Features include:
- Go HTTP server with structured JSON logging and trace IDs
- Vector-based log parsing and stateful aggregation
- Elasticsearch integration for log storage and querying
- Go monorepo structure with workspace support

See [log-parsing-forwarding/README.md](./log-parsing-forwarding/README.md) for detailed documentation.

### [metrics-alerts-template](./metrics-alerts-template/)

A comprehensive observability template for Spring Boot applications featuring:
- Prometheus metrics collection and alerting
- Grafana dashboards and visualization
- Alertmanager for alert routing and notifications
- Pre-configured alerts for common failure scenarios
- Custom business metrics instrumentation

See [metrics-alerts-template/README.md](./metrics-alerts-template/README.md) for detailed documentation.

## Getting Started

Each project contains its own README with specific setup instructions. Navigate to the project directory and follow the documentation for that project.
