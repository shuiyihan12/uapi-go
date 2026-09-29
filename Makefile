SHELL := /usr/bin/env bash

.DEFAULT_GOAL := help

GO ?= go
BIN_DIR ?= bin
ENV_FILE ?= .env
VERSION ?= latest

COMMANDS := generate generator daemon

.PHONY: help deps wsdl clean-binaries build test lint check all run start dev release clean

help: ## Show available commands
	@printf '%s\n' 'Usage: make <target>' '' 'Targets:'
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_.-]+:.*## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Tidy and download Go dependencies
	$(GO) mod tidy
	$(GO) mod download

wsdl: ## Regenerate Go structs from the archived WSDL/XSD files
	./scripts/generate-air-xsd.sh

clean-binaries: ## Remove binaries that are outside the build command list
	@mkdir -p $(BIN_DIR)
	@find $(BIN_DIR) -maxdepth 1 -type f \
		! -name generate ! -name generator ! -name daemon -delete

build: deps clean-binaries ## Build all command binaries into bin/
	@mkdir -p $(BIN_DIR)
	@for command in $(COMMANDS); do \
		echo "[INFO] Building cmd/$$command"; \
		$(GO) build -o "$(BIN_DIR)/$$command" "./cmd/$$command"; \
	done

test: deps ## Run the repository test suite
	$(GO) test ./test/...

lint: ## Check formatting and run go vet
	@files="$$(gofmt -l . | grep -v '^pkg/generated/' | grep -v '^vendor/' || true)"; \
	if [ -n "$$files" ]; then \
		echo '[ERROR] Code is not formatted. Run: gofmt -w .'; \
		echo "$$files"; \
		exit 1; \
	fi
	$(GO) vet ./...

check: test lint ## Run tests and static checks

all: deps wsdl build test lint ## Regenerate, build, test, and lint everything

run start: build ## Build and start the daemon using ENV_FILE (default: .env)
	@if [ ! -f "$(ENV_FILE)" ]; then \
		echo "[ERROR] Env file not found: $(ENV_FILE)"; \
		echo "Usage: make run ENV_FILE=path-to-env-file"; \
		exit 1; \
	fi; \
	if [ ! -x "$(BIN_DIR)/daemon" ]; then \
		echo "[ERROR] daemon binary not found: $(BIN_DIR)/daemon"; \
		exit 1; \
	fi; \
	set -a; source "$(ENV_FILE)"; set +a; \
	if [ -n "$${UAPI_AUTHORIZATION:-}" ]; then \
		echo "[WARN] UAPI_AUTHORIZATION is request-level and is ignored at startup."; \
	fi; \
	echo "[INFO] Starting daemon with env file: $(ENV_FILE)"; \
	exec "$(BIN_DIR)/daemon" -env "$${UAPI_ENV:-test}" -port "$${PORT:-8080}"

dev: ## Run the daemon from source using ENV_FILE (default: .env)
	@if [ ! -f "$(ENV_FILE)" ]; then echo "[ERROR] Env file not found: $(ENV_FILE)"; exit 1; fi
	@set -a; source "$(ENV_FILE)"; set +a; exec $(GO) run ./cmd/daemon

release: all ## Create a release archive (VERSION=... overrides the version)
	@release_dir="release/uapi-go-$(VERSION)"; \
	mkdir -p "$$release_dir"; \
	cp -r $(BIN_DIR) README.md scripts wsdl "$$release_dir/"; \
	tar -czf "release/uapi-go-$(VERSION).tar.gz" -C release "uapi-go-$(VERSION)"; \
	echo "[SUCCESS] release/uapi-go-$(VERSION).tar.gz"

clean: ## Remove generated build and release artifacts
	rm -rf $(BIN_DIR) release
