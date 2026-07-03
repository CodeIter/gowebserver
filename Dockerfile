# Dockerfile for building and running the Go server application

# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Stage 2: Run
FROM alpine:3.19

WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates

# Copy the built binary and external files from the builder stage
COPY --from=builder /app/server /app/server
COPY --from=builder /app/external /app/external

# Expose the port the application will run on
EXPOSE 8000

# Healthcheck for Docker/K8s
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

# Set the entrypoint and default command
ENTRYPOINT ["/app/server"]
CMD ["--host", "0.0.0.0", "--port", "8000", "--external", "/app/external"]
