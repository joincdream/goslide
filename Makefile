.PHONY: all build install test test-race lint complexity fmt clean help setup-tools golden-update bench check

BINARY_NAME=goslide
BIN_DIR=bin
INSTALL_DIR ?= $(HOME)/.local/bin
GO=go

GOPATH ?= $(shell $(GO) env GOPATH)
GOCYCLO ?= $(shell which gocyclo 2>/dev/null || echo $(GOPATH)/bin/gocyclo)
GOCOGNIT ?= $(shell which gocognit 2>/dev/null || echo $(GOPATH)/bin/gocognit)

# Default target
all: build

## setup-tools: Ensure dev tools (gocyclo, gocognit) are installed
setup-tools:
	@if [ ! -x "$$(command -v $(GOCYCLO))" ]; then \
		echo "==> Installing gocyclo to $(GOPATH)/bin..."; \
		$(GO) install github.com/fzipp/gocyclo/cmd/gocyclo@latest; \
	fi
	@if [ ! -x "$$(command -v $(GOCOGNIT))" ]; then \
		echo "==> Installing gocognit to $(GOPATH)/bin..."; \
		$(GO) install github.com/uudashr/gocognit/cmd/gocognit@latest; \
	fi

## build-web: Build frontend Svelte 5 and Tailwind bundle
build-web:
	@if [ -d "web/node_modules" ] && [ -x "$$(command -v npm)" ]; then \
		echo "==> Building web frontend (Svelte 5 + Tailwind)..."; \
		cd web && npm run build; \
	fi

## build: Build static single binary with CGO_ENABLED=0
build: build-web
	@echo "==> Building $(BINARY_NAME) (CGO_ENABLED=0)..."
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/goslide

## install: Build and install binary to ~/.local/bin
install: build
	@echo "==> Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp -f $(BIN_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "==> Successfully installed $(BINARY_NAME) to $(INSTALL_DIR)/$(BINARY_NAME)"

## test: Run unit tests
test:
	@echo "==> Running unit tests..."
	$(GO) test -v -cover ./...

## test-race: Run tests with data race detector
test-race:
	@echo "==> Running tests with race detector..."
	$(GO) test -v -race -cover ./...

## golden-update: Update golden file fixtures with actual test outputs
golden-update:
	@echo "==> Updating golden test fixtures..."
	$(GO) test -v ./internal/renderer/html ./internal/testutil -update

## bench: Run performance benchmarks with memory allocations
bench:
	@echo "==> Running performance benchmarks..."
	$(GO) test -v -bench=. -benchmem -run=^$$ ./cmd/goslide/...

## lint: Run Go official static analysis (go vet)
lint:
	@echo "==> Running go vet..."
	$(GO) vet ./...

## complexity: Run and report code complexity analysis (cyclomatic & cognitive)
complexity: setup-tools
	@echo "================================================================================"
	@echo "📊 Goslide Code Complexity Report (Cyclomatic & Cognitive Analysis)"
	@echo "================================================================================"
	@echo ""
	@echo "🔍 [1] Top 10 Cyclomatic Complexity (Threshold: max 15)"
	@echo "Score  Package    Function                 Location"
	@echo "--------------------------------------------------------------------------------"
	@$(GOCYCLO) -top 10 . | awk '{printf "%-6s %-10s %-24s %s\n", $$1, $$2, $$3, $$4}'
	@echo ""
	@echo "🧠 [2] Top 10 Cognitive Complexity (Threshold: max 20)"
	@echo "Score  Package    Function                 Location"
	@echo "--------------------------------------------------------------------------------"
	@$(GOCOGNIT) -top 10 . | awk '{printf "%-6s %-10s %-24s %s\n", $$1, $$2, $$3, $$4}'
	@echo ""
	@echo "================================================================================"
	@OVER_CYCLO=$$($(GOCYCLO) -over 15 .); \
	OVER_COGNIT=$$($(GOCOGNIT) -over 20 .); \
	if [ -n "$$OVER_CYCLO" ] || [ -n "$$OVER_COGNIT" ]; then \
		echo "❌ Complexity threshold exceeded!"; \
		[ -n "$$OVER_CYCLO" ] && echo "  - Cyclomatic (>15):\n$$OVER_CYCLO"; \
		[ -n "$$OVER_COGNIT" ] && echo "  - Cognitive (>20):\n$$OVER_COGNIT"; \
		exit 1; \
	else \
		echo "✅ Complexity Check Passed: All functions are within safe limits (Cyclo <= 15, Cognit <= 20)."; \
		echo "================================================================================"; \
	fi

## fmt: Format Go source code
fmt:
	@echo "==> Formatting code..."
	$(GO) fmt ./...

## check: Run all code quality checks (fmt, lint, complexity, test-race)
check: fmt lint complexity test-race
	@echo "================================================================================"
	@echo "🎉 All code quality, complexity, and unit tests passed successfully!"
	@echo "================================================================================"

## clean: Remove build artifacts and temporary files
clean:
	@echo "==> Cleaning build artifacts and temporary test outputs..."
	@rm -rf $(BIN_DIR) dist *.out *.test coverage.*
	@rm -f *.pdf *.pptx *.html error.txt .*.swp *.log
	@rm -f examples/demo/*.pdf examples/demo/*.pptx examples/demo/*.html

## help: Display this help menu
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
