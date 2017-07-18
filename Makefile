all: fmt gopass test vet lint

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go') Makefile
VERSION := 1.0.4
LDFLAGS=-ldflags "-X github.com/aviau/gopass/version.Version=$(VERSION)"

gopass: $(SOURCES)
	go build -v ${LDFLAGS} -o gopass cmd/gopass/main.go

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
