# MCP Code Review Server - Makefile
# Go-based MCP server for AI-powered code review

.PHONY: help build install clean test test-unit test-integration test-contract test-all \
        coverage lint fmt vet check run dev docker-build docker-run

# Default target
.DEFAULT_GOAL := help

# Variables
BINARY_NAME := mcp-code-review
MAIN_PATH := ./cmd/mcp-code-review
BUILD_DIR := ./build
DIST_DIR := ./dist
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Go commands
GO := go
GOFLAGS := -v
LDFLAGS := -s -w
GO_BUILD := $(GO) build $(GOFLAGS)
GO_TEST := $(GO) test $(GOFLAGS)
GO_INSTALL := $(GO) install $(GOFLAGS)

# Detect OS for platform-specific builds
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

##@ Help

help: ## Display this help message
	@echo "$(GREEN)MCP Code Review Server - Makefile Commands$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(YELLOW)<target>$(NC)\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(GREEN)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Building

build: ## Build the binary for current platform
	@echo "$(GREEN)Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Binary built: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

build-optimized: ## Build optimized binary (stripped, smaller size)
	@echo "$(GREEN)Building optimized $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO_BUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME) | awk '{print "$(GREEN)✓ Optimized binary: " $$9 " (" $$5 ")$(NC)"}'

build-all: ## Build for all platforms (linux, darwin, windows)
	@echo "$(GREEN)Building for all platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	@echo "$(YELLOW)Building for linux/amd64...$(NC)"
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "$(YELLOW)Building for linux/arm64...$(NC)"
	GOOS=linux GOARCH=arm64 $(GO_BUILD) -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "$(YELLOW)Building for darwin/amd64...$(NC)"
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@echo "$(YELLOW)Building for darwin/arm64...$(NC)"
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "$(YELLOW)Building for windows/amd64...$(NC)"
	GOOS=windows GOARCH=amd64 $(GO_BUILD) -ldflags="$(LDFLAGS)" -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "$(GREEN)✓ All binaries built in $(DIST_DIR)/$(NC)"
	@ls -lh $(DIST_DIR)

install: ## Install binary to $GOPATH/bin
	@echo "$(GREEN)Installing $(BINARY_NAME) to $$GOPATH/bin...$(NC)"
	$(GO_INSTALL) $(MAIN_PATH)
	@echo "$(GREEN)✓ Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)$(NC)"

##@ Testing

test: test-all ## Run all tests (alias for test-all)

test-all: ## Run all test suites (contract, integration, unit)
	@echo "$(GREEN)Running all tests...$(NC)"
	$(GO_TEST) ./tests/...
	@echo "$(GREEN)✓ All tests passed$(NC)"

test-unit: ## Run unit tests only
	@echo "$(GREEN)Running unit tests...$(NC)"
	$(GO_TEST) ./tests/unit/...
	@echo "$(GREEN)✓ Unit tests passed$(NC)"

test-integration: ## Run integration tests (requires API keys)
	@echo "$(GREEN)Running integration tests...$(NC)"
	@echo "$(YELLOW)Note: Requires ANTHROPIC_API_KEY, OPENAI_API_KEY, GOOGLE_API_KEY$(NC)"
	$(GO_TEST) ./tests/integration/...
	@echo "$(GREEN)✓ Integration tests passed$(NC)"

test-contract: ## Run contract tests
	@echo "$(GREEN)Running contract tests...$(NC)"
	$(GO_TEST) ./tests/contract/...
	@echo "$(GREEN)✓ Contract tests passed$(NC)"

test-verbose: ## Run all tests with verbose output
	@echo "$(GREEN)Running all tests (verbose)...$(NC)"
	$(GO_TEST) -v ./tests/...

test-short: ## Run tests with -short flag (skip slow tests)
	@echo "$(GREEN)Running short tests...$(NC)"
	$(GO_TEST) -short ./tests/...
	@echo "$(GREEN)✓ Short tests passed$(NC)"

coverage: ## Generate test coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	$(GO_TEST) -coverprofile=$(COVERAGE_FILE) ./...
	@$(GO) tool cover -func=$(COVERAGE_FILE) | tail -n 1
	@echo "$(GREEN)✓ Coverage report: $(COVERAGE_FILE)$(NC)"

coverage-html: coverage ## Generate and open HTML coverage report
	@echo "$(GREEN)Generating HTML coverage report...$(NC)"
	@$(GO) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)✓ Opening coverage report in browser...$(NC)"
	@open $(COVERAGE_HTML) 2>/dev/null || xdg-open $(COVERAGE_HTML) 2>/dev/null || echo "$(YELLOW)Open $(COVERAGE_HTML) manually$(NC)"

##@ Code Quality

fmt: ## Format Go code with gofmt
	@echo "$(GREEN)Formatting code...$(NC)"
	@$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	@$(GO) vet ./...
	@echo "$(GREEN)✓ go vet passed$(NC)"

lint: ## Run golangci-lint
	@echo "$(GREEN)Running golangci-lint...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✓ Linting passed$(NC)"; \
	else \
		echo "$(RED)✗ golangci-lint not installed. Run: brew install golangci-lint$(NC)"; \
		exit 1; \
	fi

lint-fix: ## Run golangci-lint with auto-fix
	@echo "$(GREEN)Running golangci-lint with auto-fix...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --fix ./...; \
		echo "$(GREEN)✓ Auto-fixes applied$(NC)"; \
	else \
		echo "$(RED)✗ golangci-lint not installed. Run: brew install golangci-lint$(NC)"; \
		exit 1; \
	fi

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
	@echo "$(GREEN)✓ All checks passed!$(NC)"

##@ Running

run: build ## Build and run the server
	@echo "$(GREEN)Starting MCP server...$(NC)"
	@echo "$(YELLOW)Note: Ensure API keys are set (ANTHROPIC_API_KEY, OPENAI_API_KEY, or GOOGLE_API_KEY)$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME)

dev: ## Run server in development mode (rebuild on changes - requires entr)
	@echo "$(GREEN)Starting development server...$(NC)"
	@echo "$(YELLOW)Watching for changes (press Ctrl+C to stop)...$(NC)"
	@if command -v entr > /dev/null; then \
		find . -name '*.go' | entr -r sh -c 'make build && $(BUILD_DIR)/$(BINARY_NAME)'; \
	else \
		echo "$(RED)✗ entr not installed. Run: brew install entr$(NC)"; \
		echo "$(YELLOW)Falling back to single run...$(NC)"; \
		make run; \
	fi

##@ Dependencies

deps: ## Download dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@$(GO) mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

deps-tidy: ## Tidy dependencies
	@echo "$(GREEN)Tidying dependencies...$(NC)"
	@$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies tidied$(NC)"

deps-verify: ## Verify dependencies
	@echo "$(GREEN)Verifying dependencies...$(NC)"
	@$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies verified$(NC)"

deps-upgrade: ## Upgrade all dependencies
	@echo "$(GREEN)Upgrading dependencies...$(NC)"
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies upgraded$(NC)"

##@ Cleanup

clean: ## Remove build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_FILE) $(COVERAGE_HTML)
	@$(GO) clean
	@echo "$(GREEN)✓ Cleaned$(NC)"

clean-cache: ## Clean Go build cache
	@echo "$(GREEN)Cleaning Go cache...$(NC)"
	@$(GO) clean -cache -testcache -modcache
	@echo "$(GREEN)✓ Cache cleaned$(NC)"

##@ Docker (Optional)

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t mcp-code-review:latest .
	@echo "$(GREEN)✓ Docker image built: mcp-code-review:latest$(NC)"

docker-run: ## Run Docker container
	@echo "$(GREEN)Running Docker container...$(NC)"
	@docker run --rm -it \
		-e ANTHROPIC_API_KEY="$$ANTHROPIC_API_KEY" \
		-e OPENAI_API_KEY="$$OPENAI_API_KEY" \
		-e GOOGLE_API_KEY="$$GOOGLE_API_KEY" \
		mcp-code-review:latest

##@ Information

info: ## Show build information
	@echo "$(GREEN)Build Information$(NC)"
	@echo "  Binary name:    $(BINARY_NAME)"
	@echo "  Main path:      $(MAIN_PATH)"
	@echo "  Build dir:      $(BUILD_DIR)"
	@echo "  Dist dir:       $(DIST_DIR)"
	@echo "  GOOS:           $(GOOS)"
	@echo "  GOARCH:         $(GOARCH)"
	@echo "  Go version:     $$($(GO) version)"
	@echo "  GOPATH:         $$(go env GOPATH)"

version: ## Show version information
	@if [ -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "$(GREEN)Installed version:$(NC)"; \
		$(BUILD_DIR)/$(BINARY_NAME) --version 2>/dev/null || echo "$(YELLOW)Version flag not implemented$(NC)"; \
	else \
		echo "$(YELLOW)Binary not built. Run: make build$(NC)"; \
	fi

##@ CI/CD

ci: deps check coverage ## Run full CI pipeline (deps, checks, coverage)
	@echo "$(GREEN)✓ CI pipeline completed successfully$(NC)"

pre-commit: fmt vet test-short ## Run pre-commit checks (fast)
	@echo "$(GREEN)✓ Pre-commit checks passed$(NC)"

##@ Quick Commands

quick: build test-short ## Quick build and test (no integration tests)
	@echo "$(GREEN)✓ Quick build and test completed$(NC)"

all: clean build test lint coverage ## Clean, build, test, lint, and coverage
	@echo "$(GREEN)✓ Full build pipeline completed$(NC)"
