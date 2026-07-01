# Makefile for building and running the Go server application

# Define variables
APP_NAME=my-go-server
BINARY_NAME=server
CMD_PATH=./cmd/server

# Get version info
VERSION ?= $(shell git describe --tags --always --dirty)
GIT_COMMIT ?= $(shell git rev-parse HEAD)
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION ?= $(shell go version | awk '{print $$3}')

# Build flags
LD_FLAGS := -ldflags "\
	-X '$(APP_NAME)/version.version=$(VERSION)' \
	-X '$(APP_NAME)/version.gitCommit=$(GIT_COMMIT)' \
	-X '$(APP_NAME)/version.buildDate=$(BUILD_DATE)' \
	-X '$(APP_NAME)/version.goVersion=$(GO_VERSION)'"

# Define targets
.PHONY: build run test docker-build docker-run clean

# Build the Go server binary
build:
	go build $(LD_FLAGS) -o bin/$(BINARY_NAME) $(CMD_PATH)

# Run the Go server
run:
	go run $(CMD_PATH)/main.go --port 8080

# Run tests
test:
	go test ./...

# Build the Docker image
docker-build:
	docker build -t my-go-server:latest .

# Run the Docker container
docker-run:
	docker run -p 8000:8000 my-go-server:latest

# Clean up build artifacts
clean:
	rm -rf bin/
