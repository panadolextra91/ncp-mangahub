# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Install build dependencies for CGO (SQLite needs it)
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the server binary with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o server cmd/server/main.go

# Stage 2: Final lightweight image
FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .
# Copy static assets if any (MangaHub has some embedded, but let's be safe)
COPY --from=builder /app/internal/interfaces/http/static ./internal/interfaces/http/static

# Expose all 5 protocol ports
EXPOSE 8080 9090 9191 50052 50052/udp

# Set environment variables
ENV PORT=8080
ENV TCP_PORT=9090
ENV UDP_PORT=9191
ENV GRPC_PORT=50052
ENV DB_PATH=/app/data/mangahub.db

# Create data directory
RUN mkdir -p /app/data

# Run the server
CMD ["./server"]
