BINDIR := bin

$(BINDIR):
	install -d $@

.DEFAULT_GOAL: build

.PHONY: build
build: $(BINDIR)
	@go build -o bin/plan plan.go
.PHONY: test
test:
	@go test -v ./...

.PHONY: install
install:
	@go install .
