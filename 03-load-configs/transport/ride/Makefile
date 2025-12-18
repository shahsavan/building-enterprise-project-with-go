.PHONY: fmt test tidy build run clean

BIN_DIR ?= ../bin
BIN_NAME ?= ride

fmt:
	go fmt ./...

test:
	go test ./...

integration_test:
	go test ./... -timeout 90s -tags=integration_test

tidy:
	go mod tidy

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) ./cmd

run:
	go run ./cmd

clean:
	rm -rf $(BIN_DIR)/$(BIN_NAME)
