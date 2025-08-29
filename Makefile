.DEFAULT_GOAL := build

ENV := CGO_ENABLED=1
LDFLAGS ?= -w -s

DIST_DIR := dist
BIN_NAME := k8s-tools
CMD_PATH := ./cmd/k8s-tools

BIN_NAME_SUFFIX :=

ifdef GOOS
	BIN_NAME_SUFFIX := $(BIN_NAME_SUFFIX)-$(GOOS)
endif

ifdef GOARCH
	BIN_NAME_SUFFIX := $(BIN_NAME_SUFFIX)-$(GOARCH)
endif

BINARY_PATH := $(DIST_DIR)/$(BIN_NAME)$(BIN_NAME_SUFFIX)

.PHONY: build
build:
	mkdir -p dist
	$(ENV) go build -ldflags '$(LDFLAGS)' -o $(BINARY_PATH) -v $(CMD_PATH)

.PHONY: test
test:
	$(ENV) go test -v ./...

.PHONY: lint
lint:
	@golangci-lint run -v ./...

.PHONY: lint-fix
lint-fix:
	@golangci-lint run -v --fix ./...

.PHONY: docs
docs: build
	@$(BINARY_PATH) docs generate
