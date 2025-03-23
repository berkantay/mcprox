.PHONY: build clean test run fmt lint help generate

BINARY_NAME=mcprox
GOFLAGS=-ldflags="-s -w"
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GOVERSION=$(shell go version | awk '{print $$3}')
BUILD_DIR=./build
OUTPUT=$(BUILD_DIR)/$(BINARY_NAME)

# Colors for terminal output
YELLOW=\033[0;33m
GREEN=\033[0;32m
NC=\033[0m # No Color

help: ## Display this help message
	@echo "$(YELLOW)McProx - Model Context Protocol Generator$(NC)"
	@echo ""
	@echo "$(YELLOW)Usage:$(NC)"
	@awk 'BEGIN {FS = ":.*##"; printf "  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

clean: ## Remove build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -rf dist/
	go clean
	@echo "$(GREEN)Clean complete$(NC)"

fmt: ## Format code using gofmt
	@echo "$(YELLOW)Formatting code...$(NC)"
	gofmt -s -w .
	@echo "$(GREEN)Formatting complete$(NC)"

lint: ## Run linters
	@echo "$(YELLOW)Running linters...$(NC)"
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Running golangci-lint..."; \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi
	@echo "$(GREEN)Linting complete$(NC)"

test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	go test ./... -v
	@echo "$(GREEN)Tests complete$(NC)"

build: clean ## Build the binary
	@echo "$(YELLOW)Building $(BINARY_NAME) version $(VERSION) with $(GOVERSION)...$(NC)"
	mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) -o $(OUTPUT) -v
	@echo "$(GREEN)Build complete: $(OUTPUT)$(NC)"


run:
	go run cmd/mcprox/main.go