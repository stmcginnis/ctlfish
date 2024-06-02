#
# SPDX-License-Identifier: BSD-3-Clause
#

LAST_TAG = $(shell git describe --tags --abbrev=0 --dirty)
CHANGES_SINCE_TAG = $(shell git rev-list $(shell git describe --tags --abbrev=0).. --count)
ifneq ($(CHANGES_SINCE_TAG),0)
BUILD_VERSION = $(LAST_TAG)-$(CHANGES_SINCE_TAG)
else
BUILD_VERSION = $(LAST_TAG)
endif

ROOT_DIR := $(shell git rev-parse --show-toplevel)
GOLANGCI_VERSION := "v1.57"

LD_FLAGS += -X 'github.com/stmcginnis/ctlfish/cmd.BuildVersion=$(BUILD_VERSION)'

all: lint build test

test:
	go test -v ./...

build:
	go build -o bin/ctlfish -ldflags "$(LD_FLAGS)" ./

modules:
	go mod tidy

lint:
	docker run --rm \
		-v "$(ROOT_DIR)":/src \
		-w /src \
		"golangci/golangci-lint:$(GOLANGCI_VERSION)" \
		golangci-lint run -v

clean:
	go clean
