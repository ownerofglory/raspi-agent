# ============================================
# Raspberry Pi Agent - Build Makefile
# ============================================

# --- Version Info ---
VERSION       := $(shell git rev-parse --short HEAD)
GIT_TAG       := $(shell git rev-list --tags --max-count=1)
VERSION_TAG   := $(if $(GIT_TAG),$(shell git describe --tags $(GIT_TAG)),v0)
#LDFLAGS       := -X 'main.Version=$(VERSION_TAG)-$(VERSION)'

# --- Directories ---
CMD_DIR_BE    := $(shell pwd)/cmd/raspi-agent-backend
CMD_DIR_ONB   := $(shell pwd)/cmd/raspi-agent-onboard
BIN_DIR       := $(shell pwd)/bin
RELEASE_DIR   := $(shell pwd)/release

# --- Targets ---
TARGET_NAME_BE  := raspi-agent-backend
TARGET_NAME_ONB := raspi-agent-onboard

# --- Build Config ---
TARGET_OS_BE    := $(shell go env GOOS)
TARGET_ARCH_BE  := $(shell go env GOARCH)
TARGET_OS_ONB   := linux          # Raspbian
TARGET_ARCH_ONB := arm64

# ============================================
# Build Targets
# ============================================

all: build-backend build-onboard

build-backend: fmt
	@echo "Building backend for OS=$(TARGET_OS_BE) Arch=$(TARGET_ARCH_BE)"
	@echo "Version: $(VERSION_TAG)-$(VERSION)"
	@mkdir -p $(BIN_DIR)
	go mod tidy
	GOOS=$(TARGET_OS_BE) GOARCH=$(TARGET_ARCH_BE) \
	go build -o $(BIN_DIR)/$(TARGET_NAME_BE) $(CMD_DIR_BE)/main.go

build-onboard: fmt
	@echo "Building onboard for OS=$(TARGET_OS_ONB) Arch=$(TARGET_ARCH_ONB)"
	@echo "Version: $(VERSION_TAG)-$(VERSION)"
	@mkdir -p $(BIN_DIR)
	go mod tidy
	GOOS=$(TARGET_OS_ONB) GOARCH=$(TARGET_ARCH_ONB) \
	go build -o $(BIN_DIR)/$(TARGET_NAME_ONB) $(CMD_DIR_ONB)/main.go

# ============================================
# Dev / Test / Utility
# ============================================

test: fmt
	@echo "Running tests..."
	go test -coverprofile=coverage.out ./...

fmt:
	@echo "Formatting code..."
	go fmt ./...

clean:
	@echo "Cleaning..."
	rm -rf $(BIN_DIR) $(RELEASE_DIR) coverage.out


.PHONY: all build-backend build-onboard test fmt clean release