BINARY_NAME=whatidid
BUILD_DIR=build

.PHONY: all build clean test run setup help dev-run

# Put help first so it's the default target
help:
	@echo "Available commands:"
	@echo "  help     - Show this help message"
	@echo "  build    - Build the application"
	@echo "  clean    - Remove build artifacts and database"
	@echo "  test     - Run all tests"
	@echo "  run      - Build and run the application"
	@echo "  setup    - Initial setup (creates config/wid.yaml)"
	@echo "  fetch    - Fetch new events"
	@echo "  install  - Install globally with config in ~/.config/whatidid/"
	@echo "  dev-deps - Install development dependencies"
	@echo "  lint     - Run the linter"
	@echo "  dev-run  - Quick run without building (usage: make dev-run args='fetch')"

all: clean build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f whatidid.db

test:
	@echo "Running tests..."
	@go test -v ./...

run: build
	@echo "Running..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

setup:
	@echo "Setting up development environment..."
	@if [ ! -f config/wid.yaml ]; then \
		cp config/dummy.yaml config/wid.yaml; \
		echo "Created config/wid.yaml from template"; \
	fi

fetch: build
	@echo "Fetching events..."
	@./$(BUILD_DIR)/$(BINARY_NAME) fetch

install:
	@echo "Installing to ${GOBIN:-${GOPATH}/bin}..."
	@go install
	@echo "Creating configuration directories..."
	@mkdir -p ~/.config/whatidid
	@mkdir -p ~/.local/share/whatidid
	@if [ -f config/wid.yaml ]; then \
		cp config/wid.yaml ~/.config/whatidid/wid.yaml; \
		echo "Copied existing config/wid.yaml to ~/.config/whatidid/wid.yaml"; \
	else \
		cp config/dummy.yaml ~/.config/whatidid/wid.yaml; \
		echo "Created new ~/.config/whatidid/wid.yaml from template"; \
		echo "Edit this file to configure whatidid."; \
	fi
	@echo "\nInstallation complete!"
	@echo "Binary location: ${GOBIN:-${GOPATH}/bin}/whatidid"
	@echo "Config location: ~/.config/whatidid/wid.yaml"
	@echo "Database location: ~/.local/share/whatidid/whatidid.db"
	@echo "\nMake sure ${GOBIN:-${GOPATH}/bin} is in your PATH"

# Development helper targets
dev-deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	@echo "Linting..."
	@golangci-lint run

dev-run:
	@go run main.go $(args)

.DEFAULT_GOAL := help
