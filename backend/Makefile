include ../.env

GOBIN := $(shell go env GOPATH)/bin

BINARY_NAME := server
BUILD_DIR   := bin

.PHONY: help tools fmt vet lint test tidy build run check hooks

help: 
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-8s %s\n", $$1, $$2}'

tools: 
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

fmt:
	$(GOBIN)/golangci-lint fmt ./...

vet: 
	go vet ./...

lint: 
	$(GOBIN)/golangci-lint run ./...

test:
	go test -race -cover ./...

tidy: 
	go mod tidy

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

run: build 
	./$(BUILD_DIR)/$(BINARY_NAME)

check: fmt vet lint test 

hooks:
	git config core.hooksPath ../.githooks

create_migration:
	$(GOBIN)/migrate create -ext=sql -dir=migrations -seq init

migrate_up:
	$(GOBIN)/migrate -path=migrations -database="${DB_URL}" up

migrate_down:
	$(GOBIN)/migrate -path=migrations -database="${DB_URL}" down
