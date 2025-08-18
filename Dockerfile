# ---------- Build stage for Go backend ----------
# Pick an existing Go tag. 1.22-alpine is safe/stable.
FROM golang:1.24.1-alpine AS builder

# Install git and certs (for go mod)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Go deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go


# ---------- Final runtime image ----------
FROM python:3.11-alpine

# Packages needed by scripts + healthcheck
RUN apk add --no-cache ca-certificates git curl

# Non-root user
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Python deps for scripts
COPY scripts/requirements.txt /app/scripts/requirements.txt
RUN pip install --no-cache-dir -r /app/scripts/requirements.txt

# Copy app binary from builder
COPY --from=builder /app/main /app/main

# Copy scripts and the rest of the repo (if you prefer, limit to what you need)
COPY scripts/ /app/scripts/
COPY . /app/

# Ensure producer script is executable
RUN chmod +x /app/scripts/run_producer.sh

# Permissions
RUN chown -R appuser:appgroup /app
USER appuser

EXPOSE 8080
CMD ["./main"]
