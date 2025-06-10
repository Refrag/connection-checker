# Refrag ConnectionChecker - Cross-platform build system
# Go application for network diagnostics

# Application details
APP_NAME := RefragConnectionChecker
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Build directories
BUILD_DIR := build
DIST_DIR := dist

# Color output
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RED := \033[31m
RESET := \033[0m

# Operating systems and architectures
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	linux/386 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/386 \
	freebsd/amd64 \
	openbsd/amd64

.PHONY: help all clean build-all release local test deps tidy fmt vet

# Default target
all: clean build-all

help: ## Show this help message
	@echo "$(BLUE)Refrag ConnectionChecker Build System$(RESET)"
	@echo ""
	@echo "$(YELLOW)Available targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)Supported platforms:$(RESET)"
	@for platform in $(PLATFORMS); do \
		echo "  $(BLUE)$$platform$(RESET)"; \
	done

deps: ## Download and verify dependencies
	@echo "$(YELLOW)Downloading dependencies...$(RESET)"
	go mod download
	go mod verify
	@echo "$(GREEN)Dependencies ready$(RESET)"

tidy: ## Clean up go.mod and go.sum
	@echo "$(YELLOW)Tidying Go modules...$(RESET)"
	go mod tidy
	@echo "$(GREEN)Modules tidied$(RESET)"

fmt: ## Format Go code
	@echo "$(YELLOW)Formatting code...$(RESET)"
	go fmt ./...
	@echo "$(GREEN)Code formatted$(RESET)"

vet: ## Run go vet
	@echo "$(YELLOW)Running go vet...$(RESET)"
	go vet ./...
	@echo "$(GREEN)No issues found$(RESET)"

test: ## Run tests
	@echo "$(YELLOW)Running tests...$(RESET)"
	go test -v ./...
	@echo "$(GREEN)Tests completed$(RESET)"

local: deps ## Build for current platform only
	@echo "$(YELLOW)Building for current platform...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "$(GREEN)Local build complete: $(BUILD_DIR)/$(APP_NAME)$(RESET)"

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(RESET)"
	rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "$(GREEN)Clean complete$(RESET)"

build-all: deps ## Build for all platforms
	@echo "$(BLUE)Building for all platforms...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@$(foreach platform,$(PLATFORMS), \
		$(call build_platform,$(platform)))
	@echo "$(GREEN)All builds complete!$(RESET)"
	@echo "$(YELLOW)Build artifacts in: $(BUILD_DIR)/$(RESET)"

release: clean build-all ## Create release packages
	@echo "$(BLUE)Creating release packages...$(RESET)"
	@mkdir -p $(DIST_DIR)
	@$(foreach platform,$(PLATFORMS), \
		$(call package_platform,$(platform)))
	@echo "$(GREEN)Release packages created in: $(DIST_DIR)/$(RESET)"
	@ls -la $(DIST_DIR)/

# Individual platform targets
linux-amd64: deps ## Build for Linux AMD64
	$(call build_single,linux,amd64)

linux-arm64: deps ## Build for Linux ARM64
	$(call build_single,linux,arm64)

darwin-amd64: deps ## Build for macOS Intel
	$(call build_single,darwin,amd64)

darwin-arm64: deps ## Build for macOS Apple Silicon
	$(call build_single,darwin,arm64)

windows-amd64: deps ## Build for Windows 64-bit
	$(call build_single,windows,amd64)

windows-386: deps ## Build for Windows 32-bit
	$(call build_single,windows,386)

# Build function for a single platform
define build_single
	@echo "$(YELLOW)Building for $1/$2...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	GOOS=$1 GOARCH=$2 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-$1-$2$(if $(filter windows,$1),.exe) .
	@echo "$(GREEN)Built: $(BUILD_DIR)/$(APP_NAME)-$1-$2$(if $(filter windows,$1),.exe)$(RESET)"
endef

# Build function called by build-all
define build_platform
	$(eval OS := $(word 1,$(subst /, ,$(1))))
	$(eval ARCH := $(word 2,$(subst /, ,$(1))))
	@echo "$(YELLOW)Building $(OS)/$(ARCH)...$(RESET)"
	@GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-$(OS)-$(ARCH)$(if $(filter windows,$(OS)),.exe) . || \
		(echo "$(RED)Failed to build $(OS)/$(ARCH)$(RESET)" && exit 1)
endef

# Package function for releases
define package_platform
	$(eval OS := $(word 1,$(subst /, ,$(1))))
	$(eval ARCH := $(word 2,$(subst /, ,$(1))))
	$(eval BINARY := $(APP_NAME)-$(OS)-$(ARCH)$(if $(filter windows,$(OS)),.exe))
	$(eval PACKAGE := $(APP_NAME)-$(VERSION)-$(OS)-$(ARCH))
	@echo "$(YELLOW)Packaging $(OS)/$(ARCH)...$(RESET)"
	@if [ -f "$(BUILD_DIR)/$(BINARY)" ]; then \
		mkdir -p $(DIST_DIR)/$(PACKAGE); \
		cp $(BUILD_DIR)/$(BINARY) $(DIST_DIR)/$(PACKAGE)/; \
		cp README.md $(DIST_DIR)/$(PACKAGE)/; \
		if [ "$(OS)" = "windows" ]; then \
			cd $(DIST_DIR) && zip -r $(PACKAGE).zip $(PACKAGE)/; \
		else \
			cd $(DIST_DIR) && tar -czf $(PACKAGE).tar.gz $(PACKAGE)/; \
		fi; \
		rm -rf $(DIST_DIR)/$(PACKAGE); \
		echo "$(GREEN)Packaged: $(DIST_DIR)/$(PACKAGE)$(if $(filter windows,$(OS)),.zip,.tar.gz)$(RESET)"; \
	else \
		echo "$(RED)Binary not found: $(BUILD_DIR)/$(BINARY)$(RESET)"; \
	fi
endef

# Development helpers
dev: local ## Quick development build and run
	@echo "$(BLUE)Running development build...$(RESET)"
	@./$(BUILD_DIR)/$(APP_NAME)

install: local ## Install to system (requires sudo on Unix)
	@echo "$(YELLOW)Installing $(APP_NAME)...$(RESET)"
	@if [ "$$(uname)" = "Darwin" ] || [ "$$(uname)" = "Linux" ]; then \
		sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/; \
		echo "$(GREEN)Installed to /usr/local/bin/$(APP_NAME)$(RESET)"; \
	else \
		echo "$(RED)Manual installation required on this platform$(RESET)"; \
	fi

uninstall: ## Uninstall from system
	@echo "$(YELLOW)Uninstalling $(APP_NAME)...$(RESET)"
	@if [ "$$(uname)" = "Darwin" ] || [ "$$(uname)" = "Linux" ]; then \
		sudo rm -f /usr/local/bin/$(APP_NAME); \
		echo "$(GREEN)Uninstalled from /usr/local/bin/$(RESET)"; \
	else \
		echo "$(RED)Manual uninstallation required on this platform$(RESET)"; \
	fi

info: ## Show build information
	@echo "$(BLUE)Build Information$(RESET)"
	@echo "App Name:    $(APP_NAME)"
	@echo "Version:     $(VERSION)"
	@echo "Build Time:  $(BUILD_TIME)"
	@echo "Go Version:  $$(go version)"
	@echo "Platforms:   $(words $(PLATFORMS)) supported"

# Check for required tools
check-tools: ## Check for required build tools
	@echo "$(YELLOW)Checking build tools...$(RESET)"
	@command -v go >/dev/null 2>&1 || (echo "$(RED)Go is not installed$(RESET)" && exit 1)
	@command -v git >/dev/null 2>&1 || echo "$(YELLOW)Git not found - version will be 'dev'$(RESET)"
	@echo "$(GREEN)Build tools ready$(RESET)"
