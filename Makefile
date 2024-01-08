BINDIR := bin
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)
LDFLAGS := -X 'main.version=$(VERSION)' \
           -X 'main.goversion=$(GOVERSION)'
BUILD_GOOS ?= $(shell go env GOOS)
BUILD_GOARCH ?= $(shell go env GOARCH)

RELEASE_ARTIFACTS_DIR := .release_artifacts
CHECKSUM_FILE := $(RELEASE_ARTIFACTS_DIR)/checksums.txt

$(RELEASE_ARTIFACTS_DIR):
	install -d $@

$(BINDIR):
	install -d $@


.PHONY: build
build: $(BINDIR)
	GOOS=$(BUILD_GOOS) GOARCH=$(BUILD_GOARCH) go build -ldflags "$(LDFLAGS)" -o bin/plan plan.go

.PHONY: build-standalone
build-standalone: build $(RELEASE_ARTIFACTS_DIR)
	mv bin/plan $(RELEASE_ARTIFACTS_DIR)/plan-$(VERSION).$(BUILD_GOOS).$(BUILD_GOARCH)
	shasum -a 256 $(RELEASE_ARTIFACTS_DIR)/plan-$(VERSION).$(BUILD_GOOS).$(BUILD_GOARCH) >> $(CHECKSUM_FILE)

.PHONY: test
test:
	go test -v ./...

.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)" .

.DEFAULT_GOAL := build

.PHONY: github-release
github-release:
	gh release create $(VERSION) --title 'Release $(VERSION)' \
	 	--notes-file docs/releases/$(VERSION).md $(RELEASE_ARTIFACTS_DIR)/*
