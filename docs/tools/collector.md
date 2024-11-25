# Collector Module README

## Overview

The **Collector Module** gathers metrics from various Linux distributions and package management systems to analyze dependencies and relationships between packages. This data supports the criticality analysis of open-source projects.

### Supported Distributions

- **Nix**
- **Gentoo**
- **Homebrew**
- **Debian**
- **Arch Linux**

Each distribution requires a slightly different approach to data collection, but the core process remains the same: accessing package repositories, extracting dependency information, and storing data for analysis.

## Metrics Collection Process

### Homebrew

- **Repository Cloning**: The Homebrew repository (`homebrew-core`) is cloned locally.
- **Formula Parsing**: Parses `.rb` files to collect package information and dependencies.
- **Dependency Storage**: Stores data in a relational database.
- **Generate Dependency Graph**: Visualizes package relationships.

### Gentoo

- **Repository Cloning**: The Gentoo repository is cloned.
- **Environment Sync**: Updates with `emerge --sync`.
- **Ebuild Parsing**: Extracts package data from `.ebuild` files.
- **Dependency Analysis**: Uses `equery depgraph`.
- **Storage & Visualization**: Stores data and generates a dependency graph.

### Nix

- **Environment Setup**: Uses `nix-env` to list available packages.
- **Package Retrieval**: Extracts package info using `nix eval`.
- **Dependency Analysis**: Analyzes dependencies via custom Nix expressions.
- **Database Update**: Saves data in a central database.
- **Graph Generation**: Creates a dependency graph.

### Debian

- **Repository Access**: Downloads metadata from Debian mirrors.
- **Package Parsing**: Decompresses `Packages.gz` to extract package data.
- **Dependency Analysis**: Parses dependencies.
- **Database Integration**: Stores data.
- **Generate Dependency Graph**: Visualizes dependencies.

### Arch Linux

- **Repository Access**: Downloads `.tar.gz` packages.
- **Package Extraction**: Extracts metadata files.
- **Dependency Analysis**: Parses package information.
- **Database Integration**: Stores data.
- **Graph Generation**: Creates dependency graph.

## Database Integration

Collected data from each distribution is stored in a relational database. This includes:

- **Package Information**: Basic package details like name, description, and homepage.
- **Dependency Relationships**: Data on how packages depend on each other, useful for visualizing and querying package ecosystems.

## Usage Guide

### Build and Installation

To build the project, use the `make` command from the root directory:

```
make
```

### Execution Command

After building, run the Collector module with the following command:

```
./bin/show_dispkg_deps -config=config.json -type=<distribution> [-gendot=output.dot]
```

### Parameters Explanation

- `-config`: Specifies the path to the configuration file, containing database connection details. Default is `config.json`.
- `-type`: Specifies the distribution type to collect metrics from. Options include `archlinux`, `debian`, `nix`, `homebrew`, and `gentoo`.
- `-gendot`: (Optional) Specifies the output file for a `.dot` dependency graph. Note: This option is not supported for `nix`.

### Example Commands

- **Arch Linux**:

  ```
  ./bin/show_dispkg_deps -config=config.json -type=archlinux -gendot=arch_deps.dot
  ```

- **Debian**:

  ```
  ./bin/show_dispkg_deps -config=config.json -type=debian -gendot=debian_deps.dot
  ```

- **Homebrew**:

  ```
  ./bin/show_dispkg_deps -config=config.json -type=homebrew -gendot=brew_deps.dot
  ```

- **Gentoo**:

  ```
  ./bin/show_dispkg_deps -config=config.json -type=gentoo -gendot=gentoo_deps.dot
  ```

- **Nix** (Graph generation is not supported):

  ```
  ./bin/show_dispkg_deps -config=config.json -type=nix
  ```

## Summary

The Collector Module centralizes the collection of dependency data from multiple Linux distributions, supporting criticality analysis. This unified dataset facilitates the evaluation of open-source projects, enabling better insights into their dependencies and relationships. Each distribution is handled with a tailored approach, but follows a common workflow for accessing repositories, parsing data, and storing it in a structured format for analysis.