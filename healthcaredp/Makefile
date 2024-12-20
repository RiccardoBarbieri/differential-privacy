# Makefile for healthcaredp pipeline tool

# Check for required environment variables
GOROOT ?= $(shell go env GOROOT)
GOPATH ?= $(shell go env GOPATH)

SHELL := /bin/bash

# Define the output binary
BINARY_NAME = healthcaredp
BUILD_DIR = main/build

all: check-env build



check-env:
	@echo "Checking Go environment..."
	@if [ -z "$(GOROOT)" ]; then \
        echo "GOROOT is not set. Please set GOROOT."; \
        exit 1; \
    fi
	@if [ -z "$(GOPATH)" ]; then \
        echo "GOPATH is not set. Please set GOPATH."; \
        exit 1; \
    fi
	@if [ -z $(shell command -v go) ]; then \
  		echo "go executable not found. Please add go executable to path."; \
		exit 1; \
	fi
	@echo "Go command: $(shell command -v go)"
	@echo "GOROOT: $(GOROOT)"
	@echo "GOPATH: $(GOPATH)"

build: check-env
	@echo "Downloading packages"
	@go mod tidy
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main/main.go

run: build
	@echo "Run $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

help:
	@echo "Makefile for healthcaredp pipeline tool"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all        Check environment and build the program (default)"
	@echo "  check-env  Check if the Go environment is setup correctly"
	@echo "  build      Build the program"
	@echo "  run        Build and run the program"
	@echo "  clean      Remove build artifacts"
	@echo "  help       Show this help message"