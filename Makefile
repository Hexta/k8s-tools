ENV := CGO_ENABLED=1
LDFLAGS := "-w -s"

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

.PHONY: lint-fix
lint-fix:
	@golangci-lint run -v --fix ./...

.PHONY: docs
docs: build
	@./dist/k8s-tools docs generate
