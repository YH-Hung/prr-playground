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
RUN cd services/client && go mod download

# Copy client source code
COPY services/client/ services/client/

# Build binary with optimizations
RUN cd services/client && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/client

FROM alpine:3.20

# Add metadata labels
LABEL maintainer="prr-playground"
LABEL description="HTTP client for load testing"
LABEL version="2.0"

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# Copy binary from builder
COPY --from=builder /app/client /client

# Switch to non-root user
USER appuser

ENTRYPOINT ["/client"]
