# Makefile for Vajean project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint
GOIMPORTS=goimports

# Binary name
BINARY_NAME=Vajean

## Docker
DOCKER_IMAGE_NAME=$(BINARY_NAME)-app
DOCKER_IMAGE_TAG=latest
DOCKERFILE=Dockerfile

# Build directory
BUILD_DIR=build

# Main package path
MAIN_PACKAGE=.

.PHONY: all build test clean lint deps help goimports docker-build docker-buildx docker-run docker-clean

all: test goimports fmt build

# Build the project
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run linter
lint:
	$(GOLINT) run

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run goimports
goimports:
	@which $(GOIMPORTS) > /dev/null || go install golang.org/x/tools/cmd/goimports@latest
	$(GOIMPORTS) -w ./

# Verify dependencies
verify:
	$(GOMOD) verify

# Build Docker image
docker-build:
	@echo "Building Docker image with APP_NAME=$(BINARY_NAME)..."
	docker build --build-arg APP_NAME=$(BINARY_NAME) -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) -f $(DOCKERFILE) .

docker-buildx:
	@echo "Building multi-arch Docker image with APP_NAME=$(BINARY_NAME)..."
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg APP_NAME=$(BINARY_NAME) \
		-t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) \
		-f $(DOCKERFILE) \
		--load \
		.

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

# Remove Docker image
docker-clean:
	@echo "Removing Docker image..."
	docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) || true
run:
	$(GOCMD) run $(MAIN_PACKAGE)

# Show help
help:
	@echo "Make targets:"
	@echo "  all          - Run tests and build"
	@echo "  build        - Build the binary"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  goimports    - Run goimports to format code and update imports"
	@echo "  verify       - Verify dependencies"
	@echo "  docker-build   - Build the Docker image"
	@echo "  docker-buildx  - Build the multi-arch Docker image"
	@echo "  docker-run     - Run the Docker container"
	@echo "  docker-clean   - Remove the Docker image"
	@echo "  help         - Show this help"
