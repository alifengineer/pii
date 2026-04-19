GOBIN := $(shell go env GOPATH)/bin
GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)

build:
	go build -o pii .

setup:
	@if [ ! -x "$(GOBIN)/golangci-lint" ]; then \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; \
	fi

lint: setup
	$(GOBIN)/golangci-lint run ./...

test:
	go test -race -v ./...

coverage:
	mkdir -p data
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	grep -v "mock" coverage.out > coverage.out.tmp
	go tool cover -html coverage.out.tmp -o data/coverage.html
	mv coverage.out.tmp coverage.out
	open data/coverage.html

fmt:
	go fmt ./...

check: setup test fmt lint
	@echo "✅ All checks passed!"

.PHONY: build setup lint test coverage fmt check