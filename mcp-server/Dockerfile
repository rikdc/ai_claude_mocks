# Multi-stage Docker build for Mockery MCP Server

# Build stage
FROM golang:1.24.2-alpine AS builder

# Install build dependencies
RUN apk --no-cache add ca-certificates git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mockery-mcp-server ./cmd/server/main.go

# Install mockery
RUN go install github.com/vektra/mockery/v2@v2.53.0

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates git

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binaries from build stage
COPY --from=builder /go/bin/mockery /usr/local/bin/mockery
COPY --from=builder /build/mockery-mcp-server /usr/local/bin/mockery-mcp-server

# Set permissions
RUN chmod +x /usr/local/bin/mockery /usr/local/bin/mockery-mcp-server

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set default command
CMD ["/usr/local/bin/mockery-mcp-server", "-addr", ":8080"]