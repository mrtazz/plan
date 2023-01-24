BINDIR := bin
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.goversion=$(GOVERSION)'


$(BINDIR):
	install -d $@


.PHONY: build
build: $(BINDIR)
	@go build -ldflags "$(LDFLAGS)" -o bin/plan plan.go
.PHONY: test
test:
	@go test -v ./...

.PHONY: install
install:
	@go install -ldflags "$(LDFLAGS)" .

.DEFAULT_GOAL := build
