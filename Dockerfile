# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /src

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with proper version information
ARG VERSION=0.1.7
ARG BUILD_DATE=$(date -u +"%Y-%m-%d")
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE}" -o term-rex

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates ncurses-terminfo

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /src/term-rex /app/term-rex

# Set terminal environment
ENV TERM=xterm-256color

# Set executable permissions
RUN chmod +x /app/term-rex

# Create a non-root user to run the application
RUN addgroup -S termrex && adduser -S termrex -G termrex
USER termrex

# Command to run
ENTRYPOINT ["/app/term-rex"]
