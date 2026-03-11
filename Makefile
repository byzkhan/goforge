.PHONY: build test lint clean install

VERSION ?= dev
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -ldflags "-X github.com/byzkhan/goforge/internal/version.Version=$(VERSION) -X github.com/byzkhan/goforge/internal/version.Commit=$(COMMIT) -X github.com/byzkhan/goforge/internal/version.Date=$(DATE)"

## build: Compile goforge
build:
	go build $(LDFLAGS) -o bin/goforge ./cmd/goforge

## test: Run all tests
test:
	go test ./... -v -race -count=1

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## install: Install goforge to GOPATH/bin
install:
	go install $(LDFLAGS) ./cmd/goforge

## clean: Remove build artifacts
clean:
	rm -rf bin/

## help: Show this help
help:
	@grep -E '^##' Makefile | sed 's/## //'
