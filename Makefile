# filepath: /go-git-tui/Makefile
.PHONY: build clean lint fmt run-gadd run-gcommit test install-tools check-tools

build:
	go build -o bin/gadd ./cmd/gadd
	go build -o bin/gcommit ./cmd/gcommit

clean:
	rm -rf bin/*

gadd:
	./bin/gadd

gcommit:
	./bin/gcommit

test:
	go test ./... -v

# Enhanced fmt to include imports
fmt:
	@command -v gofmt >/dev/null 2>&1 || { echo "gofmt not installed, please run 'make install-tools'"; exit 1; }
	gofmt -s -w .
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w .; \
	else \
		echo "goimports not installed, falling back to gofmt only. Run 'make install-tools' to install goimports"; \
	fi

lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "Using golangci-lint for linting..."; \
		golangci-lint run --config .golangci.yml ./...; \
	else \
		echo "golangci-lint not installed."; \
	fi

# Check if tools are installed
check-tools:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Run 'make install-tools' to install it."; exit 1; }

# Install development tools
install-tools:
	@echo "Installing required development tools..."
	@echo "Getting goimports..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

# Quick check (format + lint)
check: fmt lint