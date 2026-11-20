BIN_DIR:=bin
BIN:=$(BIN_DIR)/ambicertd
MAIN_PKG:=.

VERSION?=$(shell git describe --tags --always --dirty || echo "dev")
GIT_SHA?=$(shell git rev-parse HEAD)
BUILD_TIME?=$(shell date --iso-8601=minutes)
GOBUILD_FLAGS?=-ldflags "-X github.com/okunix/ambicertd/version.version=$(VERSION) -X github.com/okunix/ambicertd/version.gitSHA=$(GIT_SHA) -X github.com/okunix/ambicertd/version.buildTime=$(BUILD_TIME)"

.PHONY: all clean build deps run sqlc

all: clean deps build

clean: 
	rm -r $(BIN_DIR) || true

build: $(BIN)

deps: go.mod
	go mod tidy
	go mod download

$(BIN): $(MAIN_PKG)
	CGO_ENABLED=0 go build $(GOBUILD_FLAGS) -o $@ $^

run: clean deps build
	$(BIN)

sqlc:
	sqlc generate
