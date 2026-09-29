ifneq (,$(wildcard .env))
    include .env
    export
endif

OAPI_CODEGEN_VERSION := v2.8.0

APP_NAME ?= notification
BIN_DIR ?= bin
BIN ?= $(BIN_DIR)/$(APP_NAME)
MAIN_PKG ?= ./cmd/notification
GO ?= /usr/local/go/bin/go
HTTP_ADDR ?= :8080

COMPOSE_FILE ?= deploy/docker-compose.yml

DB_DSN ?= $(DATABASE_DSN)
DATABASE_URL ?= $(DATABASE_DSN)

MIGRATIONS_DIR ?= migrations

.PHONY: deps fmt lint test build run up down migrate-up migrate-down migrate-status e2e check generate

deps:
	$(GO) mod download
	@echo "OK: go modules downloaded"
	$(GO) mod tidy
	@echo "OK: go modules tidied"

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')
	@echo "OK: go files formatted"

lint:
	golangci-lint run
	@echo "OK: lint passed"

test:
	$(GO) test ./...
	@echo "OK: tests passed"

build:
	mkdir -p $(BIN_DIR)
	@echo "OK: build directory ready"
	$(GO) build -o $(BIN) $(MAIN_PKG)
	@echo "OK: binary built at $(BIN)"

run:
	$(GO) run $(MAIN_PKG)
	@echo "OK: application stopped"

up:
	docker compose -f $(COMPOSE_FILE) up -d
	@echo "OK: docker compose services started"

down:
	docker compose -f $(COMPOSE_FILE) down
	@echo "OK: docker compose services stopped"

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres '$(DATABASE_DSN)' up
	@echo "OK: migrations applied"

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres '$(DATABASE_DSN)' down
	@echo "OK: migration rolled back"

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres '$(DATABASE_DSN)' status
	@echo "OK: migration status checked"

generate:
	@mkdir -p internal/api/openapi
	@if [ -f "internal/api/openapi/oapi-codegen.yaml" ] && [ -f "api/openapi.yaml" ]; then \
		$(GO) run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) \
			-config internal/api/openapi/oapi-codegen.yaml api/openapi.yaml; \
		echo "OK: code generated"; \
	else \
		echo "Warning: api/openapi.yaml or config not found. Skipping code generation for now."; \
	fi

e2e:
	@curl -s http://localhost:8080/api/ready
	@echo "OK: e2e ready check passed"


check: fmt lint test build
	@echo "OK: all checks passed"
