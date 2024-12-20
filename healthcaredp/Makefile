# Makefile for healthcaredp pipeline tool

# Check for required environment variables
GOROOT ?= $(shell go env GOROOT)
GOPATH ?= $(shell go env GOPATH)

SHELL := /bin/bash

# Define the output binary
BINARY_NAME = healthcaredp
BUILD_DIR = main/build

ifdef INSTALL_DIR
INSTALL_DIR = $(INSTALL_DIR)
else
INSTALL_DIR = .local/bin
endif

COMPLETION_DIR = .local/share/bash-completion/completions
COMPLETION_SCRIPT = healthcaredp

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
	@go mod download
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
	@rm $(HOME)/.local/bin/$(BINARY_NAME)


install: build
	@echo "Copying executable to $(HOME)/$(INSTALL_DIR) ..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(HOME)/$(INSTALL_DIR)
	@echo "Setting up bash completion..."
	@$(BUILD_DIR)/$(BINARY_NAME) completion bash > bash_autocomplete
	@mkdir -p $(HOME)/$(COMPLETION_DIR)
	@mv bash_autocomplete $(HOME)/$(COMPLETION_DIR)/$(COMPLETION_SCRIPT)

help:
	@echo "Makefile for healthcaredp pipeline tool"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all        Check environment and build the program (default)"
	@echo "  check-env  Check if the Go environment is setup correctly"
	@echo "	 install	Install executable into local"
	@echo "  build      Build the program"
	@echo "  run        Build and run the program"
	@echo "  clean      Remove build artifacts"
	@echo "  help       Show this help message"