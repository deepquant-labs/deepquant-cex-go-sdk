.PHONY: help build test test-coverage lint fmt vet clean deps check-deps security

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the project
	@echo "Building..."
	go build -v ./...

build-examples: ## Build example applications
	@echo "Building examples..."
	go build -v ./examples/...

# Test targets
test: ## Run tests
	@echo "Running tests..."
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-short: ## Run short tests
	@echo "Running short tests..."
	go test -short -v ./...

# Code quality targets
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Security targets
security: ## Run security checks
	@echo "Running security checks..."
	gosec ./...
	nancy sleuth --input go.list

# Dependency targets
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

check-deps: ## Check for dependency updates
	@echo "Checking for dependency updates..."
	go list -u -m all

tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning..."
	go clean -cache -testcache -modcache
	rm -f coverage.out coverage.html

# Development targets
dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/sonatype-nexus-community/nancy@latest

check: fmt vet lint test ## Run all checks (format, vet, lint, test)

ci: deps check test-coverage ## Run CI pipeline locally

# Release targets
release-dry: ## Dry run release
	@echo "Dry run release..."
	git tag -l
	@echo "Next version should be: $(shell git describe --tags --abbrev=0 | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g')"

release: ## Create a new release (use VERSION=x.y.z)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make release VERSION=x.y.z"; exit 1; fi
	@echo "Creating release $(VERSION)..."
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)

# Documentation targets
docs: ## Generate documentation
	@echo "Generating documentation..."
	go doc -all ./...

docs-serve: ## Serve documentation locally
	@echo "Serving documentation at http://localhost:6060"
	godoc -http=:6060