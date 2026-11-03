BIN_DIR:=bin
BIN:=$(BIN_DIR)/ambicertd
MAIN_PKG:=.

GOBUILD_FLAGS:=

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
