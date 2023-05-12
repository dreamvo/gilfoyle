# Use bash syntax
SHELL=/bin/bash
# Go parameters
GOCMD=go
GOBINPATH=$(shell $(GOCMD) env GOPATH)/bin
GOMOD=$(GOCMD) mod
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=gotestsum
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
GOTOOL=$(GOCMD) tool
GOFMT=$(GOCMD) fmt

.PHONY: FORCE

.PHONY: all
all: deps gen build test lint fmt

.PHONY: build
build:
	$(GOBUILD) -v -o ./bin/gilfoyle .

.PHONY: test
test: deps
	$(GOTEST) --format testname -- -mod=readonly -coverprofile=cover.out -coverpkg=./... ./...

.PHONY: coverage
coverage: test
	$(GOTOOL) cover -func=cover.out

.PHONY: fmt
fmt:
	$(GOFMT) ./...
	if [ -n "$(git status --porcelain)" ]; then
	  exit 1;
	fi

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f bin/* logs/* cover.out

.PHONY: lint
lint:
	@which golangci-lint > /dev/null 2>&1 || (curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b $(GOBINPATH) v1.31.0)
	golangci-lint run -v --timeout=4m

.PHONY: deps
deps:
	$(GOMOD) download

.PHONY: gen
gen:
	$(GOCMD) generate ./...

.PHONY: install-tools
install-tools:
	$(GOGET) github.com/facebook/ent/entc/gen@v0.5.3
	$(GOGET) github.com/facebook/ent/cmd/internal/printer@v0.5.3
	$(GOGET) github.com/swaggo/swag/gen@v1.7.0
	$(GOGET) github.com/swaggo/swag/cmd/swag@v1.7.0

go.mod: FORCE
	$(GOMOD) tidy
	$(GOMOD) verify
go.sum: go.mod
