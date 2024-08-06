# Makefile
# Go commands
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOFMT = $(GOCMD) fmt

# Binary directory
BIN_DIR=./bin

# Default target
all: build

# Build the binaries
build: build_show_distpkg_deps build_enumerate_github build_show_depsdev_deps build_home2git build_gitmetricsync

build_home2git:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/home2git github.com/HUSTSecLab/criticality_score/cmd/home2git

build_show_distpkg_deps:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/show_distpkg_deps github.com/HUSTSecLab/criticality_score/cmd/show_distpkg_deps

build_enumerate_github:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/enumerate_github github.com/HUSTSecLab/criticality_score/cmd/enumerate_github

build_show_depsdev_deps:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/show_depsdev_deps github.com/HUSTSecLab/criticality_score/cmd/show_depsdev_deps

build_gitmetricsync:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/gitmetricsync github.com/HUSTSecLab/criticality_score/cmd/gitmetricsync

# Format the code
fmt:
	$(GOFMT) ./...

# Clean the build
clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/*

# Run tests
test:
	$(GOTEST) -v ./...

.PHONY: all build build_show_distpkg_deps build_enumerate_github build_show_depsdev_deps build_home2git build_gitmetricsync fmt clean test #run