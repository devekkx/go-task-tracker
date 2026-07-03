BINARY_NAME=tracker
BUILD_DIR=./bin
MAIN_PKG=.


.PHONY: build clean

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -v -race ./...

lint:
	go vet ./...

help:
	@echo 'Available targets: build clean test lint'

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(shell go env GOPATH)/bin/

test-cover:
	go test -cover ./...

.PHONY: coverage
coverage:
	go test -cover ./...
