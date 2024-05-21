# Define the output binary name and directory
BINARY_NAME=tokenclone
BIN_DIR=bin

.PHONY: compile-linux compile-windows compile-local clean

# The default target that runs when no target is specified
all: compile-osx

# The compile target for cross-compiling to Linux ARM architecture
compile-linux:
	# Create the bin directory if it doesn't exist
	mkdir -p $(BIN_DIR)
	# Build the binary for Linux ARM architecture
	env GOOS=linux GOARCH=arm go build -v -o $(BIN_DIR)/$(BINARY_NAME)

# The compile target for cross-compiling to Windows
compile-windows:
	# Create the bin directory if it doesn't exist
	mkdir -p $(BIN_DIR)
	# Build the binary for Windows
	env GOOS=windows GOARCH=amd64 go build -v -o $(BIN_DIR)/$(BINARY_NAME).exe

# The compile-local target for building on local architecture (macOS)
compile-osx:
	# Create the bin directory if it doesn't exist
	mkdir -p $(BIN_DIR)
	# Build the binary for the local architecture
	env GO111MODULE=on go build -v -o $(BIN_DIR)/$(BINARY_NAME)

# Clean up build artifacts
clean:
	rm -rf $(BIN_DIR)
