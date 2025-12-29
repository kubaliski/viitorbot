# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o viitorbot ./cmd/viitorbot

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/viitorbot .

# Create non-root user
RUN addgroup -g 1000 viitor && \
    adduser -D -u 1000 -G viitor viitor && \
    chown -R viitor:viitor /app

USER viitor

CMD ["./viitorbot"]