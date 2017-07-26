#    Copyright (C) 2017 Alexandre Viau <alexandre@alexandreviau.net>
#
#    This file is part of gopass.
#
#    gopass is free software: you can redistribute it and/or modify
#    it under the terms of the GNU General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.
#
#    gopass is distributed in the hope that it will be useful,
#    but WITHOUT ANY WARRANTY; without even the implied warranty of
#    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#    GNU General Public License for more details.
#
#    You should have received a copy of the GNU General Public License
#    along with gopass.  If not, see <http://www.gnu.org/licenses/>.

all: fmt gopass test vet lint

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go') Makefile

gopass: $(SOURCES)
	go build -v -o gopass cmd/gopass/main.go

.PHONY: test
test:
	go test -v ./...

.PHONY: install
install:
	go install ./...

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
