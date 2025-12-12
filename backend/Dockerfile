# Go Backend
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod file only first
COPY go.mod ./

# Download dependencies and generate go.sum
RUN go mod download

# Copy source code
COPY . .

# Ensure dependencies are correct
RUN go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server

# Final stage - minimal image
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=builder /server /app/server

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 3000

CMD ["/app/server"]
