# GitMetricsSync

`GitMetricsSync` is a Go module that synchronizes Git repository metrics with a PostgreSQL database. It connects to various package sources and fetches GitHub repository links to maintain a centralized record of Git metrics for analysis and monitoring.

## Features

- **Database Synchronization**: Extracts GitHub links from multiple package sources stored in a database.
- **GitHub Link Normalization**: Ensures GitHub links are stored consistently, with normalization of link cases and optional `.git` suffix.
- **Link Deletion and Addition**: Keeps the Git metrics database up-to-date by adding new links and removing outdated ones.

## Underlying Principle

The synchronization process involves collecting GitHub links from several Linux distribution package sources—specifically `Arch`, `Debian`, `Gentoo`, `Homebrew`, and `Nix`. These links are compared against the entries in the `git_metrics` table. The system follows these steps to ensure the data is always up-to-date:

- **Addition**: If a GitHub link exists in the package sources but is missing from the `git_metrics` database, it will be added.
- **Deletion**: If a GitHub link is present in the `git_metrics` table but no longer found in the package sources, it will be removed.

This ensures that the `git_metrics` table always reflects the current state of available GitHub links in the specified distributions.

## Project Structure

- **`gitmetricsync`**: Main package responsible for Git metrics synchronization.
- **`storage`**: Handles database connection and setup.
- **`main`**: Initializes the database connection and runs the synchronization process.
