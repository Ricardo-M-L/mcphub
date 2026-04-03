.PHONY: build test fmt lint clean

BINARY_NAME=mcphub
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mcphub

test:
	go test ./... -v

fmt:
	gofmt -w .
	goimports -w .

lint:
	golangci-lint run ./...

clean:
	rm -rf $(BUILD_DIR) dist

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

.DEFAULT_GOAL := build
