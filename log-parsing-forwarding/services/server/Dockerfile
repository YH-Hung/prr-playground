FROM golang:1.22-alpine AS builder

WORKDIR /build

# Copy workspace and all module files for dependency resolution
COPY go.work go.work
COPY services/server/go.mod services/server/
COPY services/client/go.mod services/client/
COPY test/integration/go.mod test/integration/
COPY internal/config/go.mod internal/config/
COPY internal/trace/go.mod internal/trace/
COPY internal/logger/go.mod internal/logger/
COPY internal/retry/go.mod internal/retry/

# Copy internal packages (shared code)
COPY internal/ internal/

# Download dependencies
RUN cd services/server && go mod download

# Copy server source code
COPY services/server/ services/server/

# Build binary
RUN cd services/server && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/server

FROM alpine:3.20

# Add metadata labels
LABEL maintainer="prr-playground"
LABEL description="HTTP server with trace logging"
LABEL version="2.0"

# Install wget for healthcheck
RUN apk add --no-cache wget

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Create log directory
RUN mkdir -p /var/log/app

# Copy binary from builder
COPY --from=builder /app/server /server

# Create entrypoint script to fix permissions
RUN echo '#!/bin/sh' > /entrypoint.sh && \
    echo 'chown -R appuser:appuser /var/log/app 2>/dev/null || true' >> /entrypoint.sh && \
    echo 'exec su-exec appuser /server' >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

# Install su-exec for switching user
RUN apk add --no-cache su-exec

EXPOSE 8080

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/entrypoint.sh"]
