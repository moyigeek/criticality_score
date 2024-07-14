# Makefile
# Go commands
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOFMT = $(GOCMD) fmt

# Default target
all: build

# Build the binaries
build: build_show_distpkg_deps build_enumerae

build_show_distpkg_deps:
	cd $(CURDIR) && $(GOBUILD) -o ./bin/show_distpkg_deps github.com/HUSTSecLab/criticality_score/cmd/show_distpkg_deps

build_enumerae:
	cd $(CURDIR) && $(GOBUILD) -o ./bin/enumerae github.com/HUSTSecLab/criticality_score/cmd/enumerate_github

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

.PHONY: all build build_show_distpkg_deps build_enumerae fmt clean test #run