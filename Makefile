# Go compiler
GO := go

# Binary name
BINARY_NAME := ./bin/file-to-env
BUILD_STEP := $(GO) build -o $(BINARY_NAME) *.go
.PHONY: all build run clean

all: build

build:
	$(BUILD_STEP)

run: build
	$(BINARY_NAME)

dev:
	air --build.cmd "$(BUILD_STEP)" --build.bin "$(BINARY_NAME)"

clean:
	rm -f $(BINARY_NAME)