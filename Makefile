# filepath: /go-git-tui/Makefile
# Go-Git-TUI Makefile
.PHONY: build clean lint fmt run gadd gcommit test test-coverage install-tools check-tools install

# Build variables
GO=go
GIT_VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X 'go-git-tui/cmd.Version=$(GIT_VERSION)' -X 'go-git-tui/cmd.Commit=$(GIT_COMMIT)' -X 'go-git-tui/cmd.BuildDate=$(BUILD_DATE)'"

build:
	@mkdir -p bin
	@mkdir -p .build
	$(GO) build $(LDFLAGS) -o .build/git-tui .
	@echo "#!/bin/bash\nexec \"\$$(dirname \"\$$0\")/../.build/git-tui\" \"\$$@\"" > bin/git-tui
	@echo "#!/bin/bash\nexec \"\$$(dirname \"\$$0\")/git-tui\" add \"\$$@\"" > bin/gadd
	@echo "#!/bin/bash\nexec \"\$$(dirname \"\$$0\")/git-tui\" commit \"\$$@\"" > bin/gcommit
	@chmod +x bin/git-tui bin/gadd bin/gcommit

clean:
	rm -rf bin/* .build/*
	rm -rf coverage.out

run:
	./bin/git-tui

gadd:
	./bin/git-tui add

gcommit:
	./bin/git-tui commit

test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

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
		echo "golangci-lint not installed. Run 'make install-tools' to install it."; \
	fi

# Check if tools are installed
check-tools:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Run 'make install-tools' to install it."; exit 1; }
	@command -v mockgen >/dev/null 2>&1 || { echo "mockgen not installed. Run 'make install-tools' to install it."; exit 1; }

# Install development tools
install-tools:
	@echo "Installing required development tools..."
	@echo "Getting goimports..."
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang/mock/mockgen@latest
	@echo "Tools installed successfully"

# Quick check (format + lint)
check: fmt lint

# Generate mocks for testing
generate-mocks:
	go generate ./...

# Install the application globally
install: build
	@echo "Installing git-tui to $(GOPATH)/bin"
	@mkdir -p $(GOPATH)/bin
	@cp .build/git-tui $(GOPATH)/bin/
	@echo "#!/bin/bash\nexec git-tui add \"\$$@\"" > $(GOPATH)/bin/gadd
	@echo "#!/bin/bash\nexec git-tui commit \"\$$@\"" > $(GOPATH)/bin/gcommit
	@chmod +x $(GOPATH)/bin/git-tui $(GOPATH)/bin/gadd $(GOPATH)/bin/gcommit
	@echo "Installation complete. The commands git-tui, gadd, and gcommit are now available."

# Release build with proper version information
release:
	@echo "Building release version $(GIT_VERSION)"
	@mkdir -p releases
	$(GO) build $(LDFLAGS) -o releases/git-tui .
	@echo "Release binary built: releases/git-tui"

# Generate docs
docs:
	@mkdir -p docs
	@if [ -f ".build/git-tui" ]; then \
		.build/git-tui generate-docs docs; \
	else \
		$(MAKE) build; \
		.build/git-tui generate-docs docs; \
	fi