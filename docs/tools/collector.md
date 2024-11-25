# Collector Module Documentation

## Overview

The collector module is responsible for collecting metrics from various Linux distributions and package management ecosystems. This module aims to understand the dependencies and relationships between packages across popular distributions, providing a unified set of metrics to support criticality analysis.

The distributions and package management systems covered in this module include:

**Nix**

**Gentoo**

**Homebrew**

**Debian**

**Arch Linux**

Each distribution requires a slightly different approach to gather metrics, but the core process is similar: fetching package information, extracting dependencies, and storing the data for further analysis.

## Supported Distributions

### Homebrew

The Homebrew collector is responsible for collecting package information from the Homebrew package repository. Below is an overview of the process:

**Repository Cloning:** The Homebrew repository (homebrew-core) is cloned to the local system.

**Formula Parsing**: Each formula file (.rb files) in the Formula directory is parsed to extract package details, including name, description, homepage, and dependencies.

**Dependency Storage**: Dependencies for each package are stored in a relational database to maintain relationships between packages.

**Generate Dependency Graph**: The dependency graph is generated for analysis, which helps visualize how different packages depend on each other.

### Gentoo

For Gentoo, the collector operates in the Gentoo environment to access the package database (Portage). The process involves:

**Repository Cloning**: Cloning the Gentoo repository (gentoo) to the local system for analysis.

**Environment Synchronization**: Synchronizing the environment with the latest package information using emerge --sync.

**Ebuild Parsing**: Parsing each ebuild file to extract package details such as name, version, description, homepage, and dependencies.

**Dependency Extraction**: Extracting package dependencies by executing commands like equery depgraph to analyze package relationships.

**Storage**: Storing the package information and dependencies in the central metrics database.

**Generate Dependency Graph**: Generating a graph representation of dependencies for better visualization and understanding of package relationships.

### Nix

The Nix collector uses the Nix package manager within a Nix distribution environment. The steps are:

**Environment Setup**: Accessing the full list of available packages using nix-env commands.

**Package Information Retrieval**: Extracting details such as name, version, homepage, description, and repository URLs using nix eval and attribute path conversion.

**Dependency Analysis**: Identifying package dependencies by evaluating build inputs through custom Nix expressions.

**Database Update**: Adding collected package information and dependencies into the shared database for further analysis.

**Generate Dependency Graph**: Generating a graph representation of dependencies to understand package relationships.

### Debian

For Debian, the collector accesses package information from the official Debian mirrors. The process involves:

**Accessing Package Repositories**: Downloading metadata from Debian mirrors, such as the Packages.gz file from the stable distribution.

**Package Parsing**: Decompressing and parsing package metadata to extract information such as name, version, description, homepage, and dependencies.

**Dependency Analysis**: Extracting and organizing the dependency relationships for each package.

**Database Integration**: Storing the package information and dependency relationships in the shared database for further analysis.

**Generate Dependency Graph**: Generating a visual graph representation of package dependencies to understand the overall structure.

### Arch Linux

The Arch Linux collector works with Pacman to access official repositories for package information. The process includes:

**Repository Access**: Downloading .tar.gz packages from the official Arch repository.

**Package Extraction**: Extracting metadata such as desc files from downloaded packages.

**Dependency Analysis**: Extracting and parsing package information, including dependencies, from the metadata.

**Database Integration**: Storing parsed package information and dependency relationships in the shared database.

**Generate Dependency Graph**: Generating a graph representation of package dependencies to provide insights into the ecosystem structure.

## Metrics Collection Workflow

Repository Access: Each distribution has its specific repository or environment access mechanism. For Homebrew, it's a GitHub repository, whereas others like Gentoo, Nix, Debian, and Arch Linux use their native environments or official mirrors.

**Package Parsing**: Package details are extracted from distribution-specific files, such as .rb files for Homebrew, ebuild files for Gentoo, and metadata files for Debian and Arch Linux.

**Dependency Analysis**: Dependencies are identified and processed recursively to build the entire dependency graph, using commands like equery for Gentoo, custom Nix expressions, or metadata parsing.

**Database Storage**: All extracted data, including package details and dependencies, are stored in a centralized database.

**Graph Generation**: The data can then be used to generate dependency graphs, providing insights into package relationships and critical points within the ecosystem.

## Database Integration

The collected metrics from each package manager are stored in a shared relational database. The data structure includes tables for:

**Package Information**: Details like name, description, homepage, etc.

**Dependency Relationships**: Relationships between packages, stored to facilitate efficient querying and graph generation.

This allows for easy analysis of package criticality, dependency counts, and overall metrics across multiple Linux distributions.

## Summary

The collector module plays a crucial role in gathering metrics from different Linux distributions and package ecosystems. The key goal is to create a comprehensive dataset of package dependencies and relationships, enabling analysis of critical open-source software projects. Each distribution has a slightly different implementation, but they all follow the general process of fetching package data, parsing dependencies, and storing the information for analysis.