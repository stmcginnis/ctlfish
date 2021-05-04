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

LD_FLAGS += -X 'github.com/stmcginnis/ctlfish/cmd.BuildVersion=$(BUILD_VERSION)'

ENVS := linux-amd64 linux-arm64 windows-amd64 darwin-amd64 darwin-arm64
CLI_JOBS := $(addprefix build-,${ENVS})

all: lint build test

test:
	go test -v ./...

build: ${CLI_JOBS}
	@echo "Built ctlfish ${BUILD_VERSION}"

build-%:
	$(eval ARCH = $(word 2,$(subst -, ,$*)))
	$(eval OS = $(word 1,$(subst -, ,$*)))

	GOOS=${OS} GOARCH=${ARCH} go build -o bin/ctlfish-${OS}_${ARCH} -ldflags "$(LD_FLAGS)" ./

modules:
	go mod tidy

lint:
	golangci-lint run -v

clean:
	go clean