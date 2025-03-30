# filepath: /go-git-tui/Makefile
# Go-Git-TUI Makefile
.PHONY: build clean lint fmt run gadd gcommit test test-coverage install-tools check-tools install

# Build variables
GO=go
GIT_VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X 'github.com/LaansDole/go-git-tui/cmd.Version=$(GIT_VERSION)' -X 'github.com/LaansDole/go-git-tui/cmd.Commit=$(GIT_COMMIT)' -X 'github.com/LaansDole/go-git-tui/cmd.BuildDate=$(BUILD_DATE)'"

build:
	@mkdir -p bin
	@mkdir -p .build
	$(GO) build $(LDFLAGS) -o .build/ggui .
	@echo "#!/bin/bash\nexec \"\$$(dirname \"\$$0\")/../.build/ggui\" \"\$$@\"" > bin/ggui
	@echo "#!/bin/bash\nexec \"\$$(dirname \"\$$0\")/ggui\" add \"\$$@\"" > bin/gadd
	@echo "#!/bin/bash\nexec \"\$$(dirname \"\$$0\")/ggui\" commit \"\$$@\"" > bin/gcommit
	@chmod +x bin/ggui bin/gadd bin/gcommit

clean:
	rm -rf bin/* .build/*
	rm -rf coverage.out

run:
	./bin/ggui

gadd:
	./bin/ggui add

gcommit:
	./bin/ggui commit

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
	@echo "Installing ggui to $(GOPATH)/bin"
	@mkdir -p $(GOPATH)/bin
	@cp .build/ggui $(GOPATH)/bin/
	@echo "#!/bin/bash\nexec ggui add \"\$$@\"" > $(GOPATH)/bin/gadd
	@echo "#!/bin/bash\nexec ggui commit \"\$$@\"" > $(GOPATH)/bin/gcommit
	@chmod +x $(GOPATH)/bin/ggui $(GOPATH)/bin/gadd $(GOPATH)/bin/gcommit
	@echo "Installation complete. The commands ggui, gadd, and gcommit are now available."

# Release build with proper version information
release:
	@echo "Building release version $(GIT_VERSION)"
	@mkdir -p releases
	$(GO) build $(LDFLAGS) -o releases/ggui .
	@echo "Release binary built: releases/ggui"

# Generate docs
docs:
	@mkdir -p docs
	@if [ -f ".build/ggui" ]; then \
		.build/ggui generate-docs docs; \
	else \
		$(MAKE) build; \
		.build/ggui generate-docs docs; \
	fi