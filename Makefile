.PHONY: build test test-unit test-cover lint clean install run help

BINARY_NAME=haft
BUILD_DIR=bin
MAIN_PATH=./cmd/haft
COVERAGE_FILE=coverage.out

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run: build
	@$(BUILD_DIR)/$(BINARY_NAME)

install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(MAIN_PATH)

test:
	@echo "Running all tests..."
	@go test ./... -v

test-unit:
	@echo "Running unit tests..."
	@go test ./internal/... -v

test-cover:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=$(COVERAGE_FILE)
	@go tool cover -func=$(COVERAGE_FILE)

test-cover-html: test-cover
	@go tool cover -html=$(COVERAGE_FILE)

lint:
	@echo "Running linter..."
	@golangci-lint run ./...

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVERAGE_FILE)

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  run          - Build and run"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  test         - Run all tests"
	@echo "  test-unit    - Run unit tests only"
	@echo "  test-cover   - Run tests with coverage"
	@echo "  lint         - Run golangci-lint"
	@echo "  fmt          - Format code"
	@echo "  vet          - Run go vet"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Download and tidy dependencies"
