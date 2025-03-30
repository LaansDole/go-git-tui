# filepath: /go-git-tui/Makefile
# Go-Git-TUI Makefile
.PHONY: build clean lint fmt run gadd gcommit test test-coverage install-tools check-tools install

# Build variables
GO=go

build:
	@mkdir -p bin
	@mkdir -p .build
	$(GO) build -o .build/go-git-tui .
	@echo "#!/bin/bash\nexec \"\$$(dirname \"\$$0\")/../.build/go-git-tui\" \"\$$@\"" > bin/go-git-tui
	@echo "#!/bin/bash\nexec \"\$$(dirname \"\$$0\")/go-git-tui\" add \"\$$@\"" > bin/gadd
	@echo "#!/bin/bash\nexec \"\$$(dirname \"\$$0\")/go-git-tui\" commit \"\$$@\"" > bin/gcommit
	@chmod +x bin/go-git-tui bin/gadd bin/gcommit

clean:
	rm -rf bin/* .build/*
	rm -rf coverage.out

run:
	./bin/go-git-tui

gadd:
	./bin/go-git-tui add

gcommit:
	./bin/go-git-tui commit

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
	@echo "Installing go-git-tui to $(GOPATH)/bin"
	@mkdir -p $(GOPATH)/bin
	@cp .build/go-git-tui $(GOPATH)/bin/
	@echo "#!/bin/bash\nexec go-git-tui add \"\$$@\"" > $(GOPATH)/bin/gadd
	@echo "#!/bin/bash\nexec go-git-tui commit \"\$$@\"" > $(GOPATH)/bin/gcommit
	@chmod +x $(GOPATH)/bin/go-git-tui $(GOPATH)/bin/gadd $(GOPATH)/bin/gcommit
	@echo "Installation complete. The commands go-git-tui, gadd, and gcommit are now available."

# Build a release version
release:
	@echo "Building release version"
	@mkdir -p releases
	$(GO) build -o releases/go-git-tui .
	@echo "Release binary built: releases/go-git-tui"

# Generate docs
docs:
	@mkdir -p docs
	@if [ -f ".build/go-git-tui" ]; then \
		.build/go-git-tui generate-docs docs; \
	else \
		$(MAKE) build; \
		.build/go-git-tui generate-docs docs; \
	fi