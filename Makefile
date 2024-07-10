# Makefile

# Go commands
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOFMT = $(GOCMD) fmt

# Default target
all: build

# Build the binary
build:
	$(GOBUILD) -o ./bin/show_distpkg_deps github.com/HUSTSeclab/criticality_score/cmd/show_distpkg_deps

# Format the code
fmt:
	$(GOFMT) ./...

# Clean the build
clean:
	$(GOCLEAN)
	rm -f bin/*

# Run tests
test:
	$(GOTEST) -v ./...

# Run the binary
#run: build
#	./$(BINARY_NAME)

.PHONY: all build fmt clean test #run
