VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

BINARY_NAME = gstest
DOCKER_IMAGE = gstest

.PHONY: build docker test lint clean help

## build: Build the binary
build:
	go build -ldflags="-w -s -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)" \
		-o bin/$(BINARY_NAME) ./cmd/gstest

## docker: Build the Docker image
docker:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		-t $(DOCKER_IMAGE):$(VERSION) \
		-t $(DOCKER_IMAGE):latest \
		.

## test: Run tests
test:
	go test -v -race -coverprofile=coverage.out ./...

## lint: Run linter
lint:
	golangci-lint run ./...

## clean: Clean build artifacts
clean:
	rm -rf bin/ coverage.out

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
