# ============================================================
# Burnless Makefile
# ============================================================

BINARY_NAME   := burnless
AGENT_NAME    := burnless-agent
MODULE        := github.com/burnless/burnless
VERSION       := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS       := -ldflags "-X main.version=$(VERSION)"
GOFLAGS       :=
BUILD_DIR     := ./bin
GO            := go

.PHONY: all build build-agent build-all test test-race lint fmt vet clean setup help

## ── build ───────────────────────────────────────────────────

all: build

build: ## Build the burnless CLI
	@echo "→ Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/burnless

build-agent: ## Build the burn rate agent
	@echo "→ Building $(AGENT_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(AGENT_NAME) ./cmd/agent

build-all: build build-agent ## Build all binaries

cross-build: ## Build for Linux, macOS, Windows (amd64 + arm64)
	@echo "→ Cross-compiling..."
	GOOS=linux   GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64   ./cmd/burnless
	GOOS=linux   GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64   ./cmd/burnless
	GOOS=darwin  GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64  ./cmd/burnless
	GOOS=darwin  GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64  ./cmd/burnless
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/burnless

## ── test ────────────────────────────────────────────────────

test: ## Run unit tests
	@echo "→ Running tests..."
	$(GO) test ./... -count=1

test-race: ## Run tests with race detector
	@echo "→ Running tests with race detector..."
	$(GO) test -race ./... -count=1

test-coverage: ## Run tests with coverage report
	@echo "→ Running tests with coverage..."
	$(GO) test ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "→ Coverage report: coverage.html"

test-integration: ## Run integration tests (requires Docker)
	@echo "→ Running integration tests..."
	$(GO) test ./tests/integration/... -tags=integration -count=1

## ── quality ─────────────────────────────────────────────────

lint: ## Run golangci-lint
	@echo "→ Running linter..."
	golangci-lint run ./...

fmt: ## Format all Go files
	@echo "→ Formatting..."
	$(GO) fmt ./...
	goimports -w .

vet: ## Run go vet
	@echo "→ Running go vet..."
	$(GO) vet ./...

## ── dev ─────────────────────────────────────────────────────

run: build ## Build and run the CLI
	$(BUILD_DIR)/$(BINARY_NAME)

setup: ## Install dev dependencies
	@echo "→ Installing dev tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/goreleaser/goreleaser@latest
	@echo "✓ Dev tools installed"

generate: ## Run go generate
	$(GO) generate ./...

tidy: ## Tidy go.mod and go.sum
	$(GO) mod tidy

## ── release ─────────────────────────────────────────────────

snapshot: ## Build a snapshot release (local, no publish)
	goreleaser release --snapshot --clean

release: ## Create and publish a release (requires GITHUB_TOKEN)
	goreleaser release --clean

## ── clean ───────────────────────────────────────────────────

clean: ## Remove build artifacts
	@echo "→ Cleaning..."
	rm -rf $(BUILD_DIR) dist/ coverage.out coverage.html

## ── help ────────────────────────────────────────────────────

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
	  | sort \
	  | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
