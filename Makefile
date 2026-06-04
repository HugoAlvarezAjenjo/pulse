.PHONY: build test test-short test-integration lint clean install run fmt cover

BINARY := pulse
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/pulse

run: build
	./$(BUILD_DIR)/$(BINARY)

test:
	go test ./... -v -race

test-short:
	go test ./... -short

test-integration:
	go test ./tests/integration/... -v -race

cover:
	go test ./... -coverprofile=coverage.out -race
	go tool cover -func=coverage.out

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR) coverage.out

install:
	go install ./cmd/pulse

fmt:
	go fmt ./...
	goimports -w .

.DEFAULT_GOAL := build
