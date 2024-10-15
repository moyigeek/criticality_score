# Criticality Score Tool

## Quick Start

make sure `docker` is installed, and run the following.

```
./setup.sh
```

After the script finish, try to connect to database (the 
password is stored in `data/DB_PASSWD` and populate 
git_link fields in arch_packages and debian_packages 
manually and finally run following command.

```
docker compose exec app bash /gitlink.sh
```

## Overview

This tool is designed to calculate the criticality score of packages in Arch Linux and Debian distributions. It can generate dependency graphs in DOT format and CSV files with reference counts.

## File Structure

- `cmd/show_distpkg_deps/main.go`: The main entry point of the tool.
- `pkg/collector/archlinux/archlinux.go`: Contains functions to process Arch Linux packages.
- `pkg/collector/debian/debian.go`: Contains functions to process Debian packages.
- `go.mod`: Go module dependencies.

## Usage

To use this tool, you can simply run the provided Makefile. The tool supports generating DOT graphs and CSV files for both Arch Linux and Debian packages.

### Generating DOT Graphs and CSV Files

To generate a DOT graph and CSV file for Arch Linux or Debian packages, use the following commands:

First, compile the tool:

```sh
make
```

Second, generate the DOT graph and CSV file for Debian or Arch Linux packages, in the following commands, `gendot` and `<output_file.dot>` are optional arguments:

For Debian:
```sh
    ./bin/show_distpkg_deps debian gendot output_file.dot
```
For Arch Linux:
```sh
    ./bin/show_distpkg_deps archlinux gendot output_file.dot
```

By using the `gendot` argument, the tool will generate a DOT graph and save it to the specified file, we can use Graphviz to visualize the graph.

```sh
dot -Tpng output_file.dot -o output_file.png
```

### Using the enumerate_github tools

The `enumerate_github` tools are designed to enumerate GitHub repositories and gather relevant data. Detailed usage instructions for these tools can be found in the `cmd/enumerate_github/README.md` file.

To use the `enumerate_github` tools, navigate to the `cmd/enumerate_github` directory and follow the instructions provided in the README file there.

This will provide you with detailed steps on how to use the tools effectively.

### Using the show_depsdev_deps tools

```sh
./bin/show_depsdev_deps input.txt output.txt
```

By using the `show_depsdev_deps` tools, you can use deps.dev api to collect the dependencies of the projects in the input file which enumerate_github tools has generated.
