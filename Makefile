all: gopass test

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

gopass: $(SOURCES)
	go build -o gopass cmd/gopass/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: install
install:
	go install ./...

.PHONY: clean
clean:
	rm -rf gopass
