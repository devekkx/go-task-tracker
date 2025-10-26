BINARY_NAME=tracker
BUILD_DIR=./bin

.PHONY: build clean

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) .

clean:
	rm -rf $(BUILD_DIR)
