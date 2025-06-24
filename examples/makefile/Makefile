# Makefile for Mock Generation Example
# This demonstrates the simplest possible approach to reliable mock generation

.PHONY: help install-mockery generate-mocks clean-mocks test build all

# Default target
all: clean-mocks generate-mocks test ## Clean, generate mocks, and run tests

# Help target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Install mockery if not present
install-mockery: ## Install mockery v2.53+ if not present
	@if ! command -v mockery >/dev/null 2>&1; then \
		echo "Installing mockery..."; \
		go install github.com/vektra/mockery/v2@latest; \
	else \
		echo "Mockery already installed: $$(mockery --version)"; \
	fi

# Generate all mocks
generate-mocks: install-mockery ## Generate mocks for all interfaces
	@echo "Generating mocks..."
	mockery --config=.mockery.yaml
	@echo "Mock generation complete. Generated files:"
	@ls -la mocks/

# Clean generated mocks
clean-mocks: ## Remove all generated mock files
	@echo "Cleaning generated mocks..."
	rm -rf mocks/
	@echo "Mock files cleaned."

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	go test ./internal/service/... -v

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test ./internal/service/... -cover -v

# Build everything to verify compilation
build: ## Build all packages to verify compilation
	@echo "Building all packages..."
	go build ./internal/...
	go build ./mocks/...
	@echo "Build successful."

# Verify mock generation and run tests
verify: generate-mocks build test ## Generate mocks, build, and test (full verification)
	@echo "Verification complete: mocks generated, code compiled, tests passed."

# Development workflow
dev: clean-mocks generate-mocks test ## Development workflow: clean, generate, test

# Show mockery version
mockery-version: install-mockery ## Show installed mockery version
	mockery --version

# Show generated mock files
list-mocks: ## List all generated mock files
	@if [ -d "mocks" ]; then \
		echo "Generated mock files:"; \
		find mocks -name "*.go" -type f | sort; \
	else \
		echo "No mock files found. Run 'make generate-mocks' first."; \
	fi

# Validate mockery configuration
validate-config: ## Validate .mockery.yaml configuration
	@if [ -f ".mockery.yaml" ]; then \
		echo "Mockery configuration found:"; \
		cat .mockery.yaml; \
	else \
		echo "No .mockery.yaml found. Creating default configuration..."; \
		$(MAKE) create-config; \
	fi

# Create default mockery configuration
create-config: ## Create default .mockery.yaml configuration
	@echo "Creating default .mockery.yaml..."
	@cat > .mockery.yaml << 'EOF'\
with-expecter: true\
packages:\
  github.com/example/makefile-mocks/internal/interfaces:\
    config:\
      all: true\
      output: "./mocks"\
      outpkg: "mocks"\
      filename: "Mock{{.InterfaceName}}.go"\
EOF
	@echo ".mockery.yaml created."

# Debug mockery execution (dry run)
debug-mockery: install-mockery ## Show what mockery would generate (dry run)
	@echo "Mockery dry run (showing what would be generated):"
	mockery --config=.mockery.yaml --dry-run

# Quick test of specific interface
test-user-service: ## Run only UserService tests
	go test ./internal/service/... -run TestUserService -v

# Generate and test in one command
quick: ## Quick workflow: generate mocks and test
	@$(MAKE) generate-mocks && $(MAKE) test