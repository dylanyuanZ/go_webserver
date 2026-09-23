BINARY := server
BIN_DIR := bin
PKG := ./...

.PHONY: all build run test check fmt vet clean

all: check build

build:
	go build -o $(BIN_DIR)/$(BINARY) ./cmd/server

run:
	go run ./cmd/server

test:
	go test $(PKG)

vet:
	go vet $(PKG)

fmt:
	gofmt -l -w .

check: vet test

clean:
	rm -rf $(BIN_DIR)
