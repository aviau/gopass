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

GO_IMPORT_PATH := github.com/aviau/gopass

gopass: $(SOURCES)
	go build -v -o gopass ${GO_IMPORT_PATH}/cmd/gopass

.PHONY: test
test:
	go test -v ${GO_IMPORT_PATH}/...

.PHONY: install
install:
	go install ${GO_IMPORT_PATH}/...

.PHONY: clean
clean:
	rm -rf gopass

.PHONY: vet
vet: install
	go get github.com/stretchr/testify/assert
	go vet -v ${GO_IMPORT_PATH}/...

.PHONY: lint
lint:
	golint ${GO_IMPORT_PATH}/...

.PHONY: fmt
fmt:
	go fmt ${GO_IMPORT_PATH}/...


.PHONY: get-deps
get-deps:
	go get -t ${GO_IMPORT_PATH}/...
