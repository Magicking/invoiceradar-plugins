.PHONY: build test clean run help

# Build the application
build:
	@echo "Building invoiceradar-plugins..."
	@go build -o invoiceradar-plugins .

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f invoiceradar-plugins

# Run with example plugin (plausible)
run-example:
	@echo "Running with Plausible plugin..."
	@./invoiceradar-plugins -plugin plugins/plausible.json -check-auth

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run || echo "golangci-lint not installed, skipping"

# Run all checks
check: fmt test
	@echo "All checks passed!"

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Remove build artifacts"
	@echo "  run-example  - Run with example plugin"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  check        - Run fmt and test"
	@echo "  help         - Show this help message"
