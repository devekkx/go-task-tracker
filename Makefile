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
