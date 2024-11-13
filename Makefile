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
build: build_show_distpkg_deps build_enumerate_github build_show_depsdev_deps \
	build_home2git build_gitmetricsync build_githubmetrics \
	build_gen_scores build_package_calculator update_git_metrics \
	apiserver

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

build_githubmetrics:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/githubmetrics github.com/HUSTSecLab/criticality_score/cmd/githubmetrics

build_gen_scores:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/gen_scores github.com/HUSTSecLab/criticality_score/cmd/gen_scores

build_package_calculator:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/package_calculator github.com/HUSTSecLab/criticality_score/cmd/package_calculator

update_git_metrics:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/update_git_metrics github.com/HUSTSecLab/criticality_score/pkg/collector_git/cmd/integrate


apiserver:
	cd $(CURDIR) && $(GOBUILD) -o $(BIN_DIR)/apiserver github.com/HUSTSecLab/criticality_score/cmd/apiserver

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

.PHONY: all build build_show_distpkg_deps build_enumerate_github \
	build_show_depsdev_deps build_home2git build_gitmetricsync \
	build_githubmetrics build_gen_scores build_package_calculator \
	update_git_metrics apiserver \
	fmt clean test #run
