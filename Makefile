ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

.PHONY: lint-prepare
lint-prepare:
	@echo "Preparing Linter"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

.PHONY: lint
lint:
	@echo "Applying linter"
	./bin/golangci-lint run --timeout=2m ./...

.PHONY: build
build:
	go build -i -o go-openapi-converter app/main.go