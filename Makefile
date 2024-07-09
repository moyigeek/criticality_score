# Makefile

# Go commands
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOFMT = $(GOCMD) fmt
BINARY_NAME = main

# Default target
all: build

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) .

# Format the code
fmt:
	$(GOFMT) ./...

# Clean the build
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Run the binary
run: build
	./$(BINARY_NAME)

.PHONY: all build fmt clean test run
