# DepsDev Synchronization

## Overview

This project is designed to fetch dependency information from the [deps.dev API](https://deps.dev/) and update a PostgreSQL database with the number of dependents for each project in the system. It retrieves the latest version of each project from the Deps.dev API and queries for the number of direct and indirect dependents for that version.

The project consists of the following main components:

- **Database Operations**: Queries and updates the PostgreSQL database to store the dependent counts for different GitHub, GitLab, and Bitbucket projects.
- **API Integration**: Connects to the [deps.dev](https://deps.dev/) API to fetch the latest versions of projects and their dependency counts.
- **Project Synchronization**: Processes each project in the database, retrieves the latest version from deps.dev, and queries for its dependent information.

## Project Structure

- `collector_depsdev`: This package contains the logic for querying the Deps.dev API and updating the PostgreSQL database.
- `storage`: Handles database connection and initialization.
- `gitmetricsync`: The main package that orchestrates the synchronization of project data.

### Key Functions

- `Run(configPath string)`: Initializes the database connection, queries the `git_metrics` table for project links, and processes each project.
- `updateDatabase(link, projectName string, dependentCount int)`: Updates the `depsdev_count` in the database for a given project.
- `getLatestVersion(owner, repo, projectType string)`: Retrieves the latest version of a project from the Deps.dev API.
- `queryDepsDev(link, projectType, projectName, version string)`: Queries the Deps.dev API to get the dependent count for a specific version of a project.
- `getProjectTypeFromDB(link string)`: Retrieves the ecosystem (project type) of a project from the database.

## How It Works

1. **Database Connection:**
   - The script connects to a PostgreSQL database using the credentials provided in `config.json`.
2. **Git Project Links:**
   - The script retrieves a list of Git project links (e.g., GitHub, GitLab, Bitbucket) from the `git_metrics` table in the database.
3. **Fetch Latest Version:**
   - For each project, the script queries the [deps.dev API](https://api.deps.dev/) to get the latest version available for the project.
4. **Dependent Information:**
   - After obtaining the latest version, the script queries the Deps.dev API again to get the number of dependents for that version (both direct and indirect).
5. **Database Update:**
   - Finally, the dependent count is updated in the `git_metrics` table in the database, ensuring that the `depsdev_count` column reflects the number of dependents for the latest version of the project.

## Troubleshooting

- **Database Connection Issues**: Ensure your PostgreSQL instance is running and that the credentials in `config.json` are correct.
- **API Rate Limiting**: If the script exceeds the API rate limits of Deps.dev, consider implementing retries or rate-limiting logic.
- **Error Logs**: Check the logs for any errors in fetching data from Deps.dev or database queries. The script logs any issues encountered during execution.
