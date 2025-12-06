# Multi-stage build for wave-backend
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# Minimal runtime image
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

# Provide a default for DB password; override in runtime env

USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
