# syntax=docker/dockerfile:1.4

# ==========================================
# STAGE 1: Build the Go Backend
# ==========================================
FROM golang:1.26.3-bookworm AS builder
WORKDIR /app

# Cache Go modules for faster rebuilds
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

# Build a highly optimized, static Go binary
RUN --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -trimpath -ldflags "-s -w" -o /app/blan-server ./main.go

# ==========================================
# STAGE 2: The Final Production Image
# ==========================================
FROM debian:bookworm-slim
WORKDIR /app

# Install standard C++ runtime libraries (required by blan)
RUN apt-get update && apt-get install -y --no-install-recommends libc6 libstdc++6 ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

# Create a non-root user for security (prevents container escape vulnerabilities)
RUN addgroup --system app && adduser --system --ingroup app appuser

# Copy Go API and the C++ binary
COPY --from=builder /app/blan-server /app/blan-server
COPY blan /app/blan

# Ensure both binaries are executable
RUN chmod +x /app/blan /app/blan-server

# Create required directories and transfer ownership to the secure appuser
RUN mkdir -p /app/workspace /app/strata_cache_data \
	&& chown -R appuser:app /app

USER appuser

EXPOSE 8080
CMD ["./blan-server"]