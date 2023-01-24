BINDIR := bin

$(BINDIR):
	install -d $@


.PHONY: build
build: $(BINDIR)
	@go build -o bin/plan plan.go
.PHONY: test
test:
	@go test -v ./...

.PHONY: install
install:
	@go install .

.DEFAULT_GOAL := build
