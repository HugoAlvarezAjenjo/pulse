.PHONY: build test lint clean install run

BINARY := pulse
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/pulse

run: build
	./$(BUILD_DIR)/$(BINARY)

test:
	go test ./... -v

test-short:
	go test ./... -short

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR)

install:
	go install ./cmd/pulse

fmt:
	go fmt ./...
	goimports -w .

.DEFAULT_GOAL := build
