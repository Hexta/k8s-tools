CI_VERSION = $(shell git describe --tags --abbrev=8 --dirty --always --long)

LDFLAGS := -w -s
LDFLAGS := "$(LDFLAGS) -X 'github.com/Hexta/k8s-tools/pkg/version.version=${CI_VERSION}'"

ENV := CGO_ENABLED=1

BIN_NAME_SUFFIX :=

ifdef GOOS
		BIN_NAME_SUFFIX := $(BIN_NAME_SUFFIX)-$(GOOS)
endif

ifdef GOARCH
		BIN_NAME_SUFFIX := $(BIN_NAME_SUFFIX)-$(GOARCH)
endif

.PHONY: build
build:
	mkdir -p dist
	$(ENV) go build -ldflags=$(LDFLAGS) -o dist/k8s-tools$(BIN_NAME_SUFFIX) -v ./cmd/k8s-tools

.PHONY: test
test:
	$(ENV) go test -v ./...

.PHONY: lint
lint:
	@golangci-lint run -v ./...

.PHONY: docs
docs: build
	@./dist/k8s-tools docs generate
