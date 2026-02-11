.DEFAULT_GOAL := help

.PHONY: help test build lint all clean

help:
	@echo "Available targets:"
	@echo "  make test    - Run Go tests"
	@echo "  make build   - Build Go binary (bin/rooda)"
	@echo "  make lint    - Run go vet and shellcheck (if available)"
	@echo "  make all     - Run lint, test, and build"
	@echo "  make clean   - Remove build artifacts"
	@echo "  make help    - Show this help message"

test:
	@echo "Running Go tests..."
	go test ./...

build:
	@echo "Building Go binary..."
	go build -o bin/rooda ./cmd/rooda

lint:
	@echo "Running go vet..."
	go vet ./...
	@echo "Checking for shellcheck..."
	@if command -v shellcheck >/dev/null 2>&1 && [ -f archive/src/rooda.sh ]; then \
		echo "Running shellcheck on archived bash script..."; \
		shellcheck archive/src/rooda.sh; \
	elif [ ! -f archive/src/rooda.sh ]; then \
		echo "SKIP: archive/src/rooda.sh not found"; \
	else \
		echo "WARN: shellcheck not installed, skipping bash lint"; \
	fi

all: lint test build

clean:
	rm -rf bin/
