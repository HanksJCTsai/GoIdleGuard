# Makefile – Build, test and package for multiple platforms

BIN_DIR := bin

.PHONY: all build linux windows macos ui daemon test clean

all: test build

# Build for host OS
build: ui daemon

# Build UI
ui:
	@echo "→ Building UI (host OS)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/app-ui ./cmd/main
	@echo "   → $(BIN_DIR)/app-ui"

# Build Daemon
daemon:
	@echo "→ Building Daemon (host OS)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/app-daemon ./cmd/daemon
	@echo "   → $(BIN_DIR)/app-daemon"

# Run tests
test:
	@echo "→ Running tests..."
	@go test ./...
	@echo "   ✔ All tests passed."

# --------------------
# Cross‑platform targets
# --------------------

# Linux amd64
linux: linux-ui linux-daemon

linux-ui:
	@echo "→ Building UI for Linux/amd64..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/app-ui-linux ./cmd/main
	@echo "   → $(BIN_DIR)/app-ui-linux"

linux-daemon:
	@echo "→ Building Daemon for Linux/amd64..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/app-daemon-linux ./cmd/daemon
	@echo "   → $(BIN_DIR)/app-daemon-linux"

# Windows amd64
windows: windows-ui windows-daemon

windows-ui:
	@echo "→ Building UI for Windows/amd64..."
	@mkdir -p $(BIN_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/app-ui.exe ./cmd/main
	@echo "   → $(BIN_DIR)/app-ui.exe"

windows-daemon:
	@echo "→ Building Daemon for Windows/amd64..."
	@mkdir -p $(BIN_DIR)
	@GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/app-daemon.exe ./cmd/daemon
	@echo "   → $(BIN_DIR)/app-daemon.exe"

# macOS amd64
macos: macos-ui macos-daemon

macos-ui:
	@echo "→ Building UI for macOS/amd64..."
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/app-ui-darwin ./cmd/main
	@echo "   → $(BIN_DIR)/app-ui-darwin"

macos-daemon:
	@echo "→ Building Daemon for macOS/amd64..."
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/app-daemon-darwin ./cmd/daemon
	@echo "   → $(BIN_DIR)/app-daemon-darwin"

# --------------------
# Clean
# --------------------

clean:
	@echo "→ Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@echo "   ✔ Removed $(BIN_DIR)/"
