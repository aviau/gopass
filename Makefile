all: gopass test vet lint

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go') Makefile
VERSION := 0.1.0
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

gopass: $(SOURCES)
	go build ${LDFLAGS} -o gopass cmd/gopass/main.go

.PHONY: test
test:
	go test ${LDFLAGS} -v ./...

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	rm -rf gopass

.PHONY: vet
vet:
	go vet -v ./...

.PHONY: lint
lint:
	golint ./...

.PHONY: fmt
fmt:
	go fmt ./...


.PHONY: get-deps
get-deps:
	go get -t ./...
