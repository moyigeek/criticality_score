# Criticality Score

English | [简体中文](./README.zh_CN.md) 

## Description

This project aims to evaluate and rank open-source projects based on their importance within the ecosystem and maintainance status. Unlike traditional methods that rely solely on GitHub metrics, this project incorporates data from various Linux distributions to provide a more comprehensive analysis. 

By collecting and analyzing metrics from multiple sources, this project offers a robust and comprehensive framework for assessing the criticality of open-source projects.

## Difference from [ossf/criticality_score](https://github.com/ossf/criticality_score)

- **Distribution Dependents**: Collects data from various Linux distributions (e.g. Debian, Nix, Gentoo) to evaluate the dependency of open-source software.
- **Support for All Git Repositories**: Analyzes repositories from any Git platform, not just GitHub.
- **Comprehensive Metrics Collection**: Gathers a wider and more precise metrics from Git repositories and package managers, for example, the number of commits, organization count is more accurate than GitHub API.
- **No Dependency on Google Cloud or BigQuery**: `ossf/criticality_score` depends on Google Cloud service, make it hard to migrate to other platforms. This project runs independently of specific cloud services, ensuring ease of deployment.
- **Easy Deployment**: Just running a script, the system will be setup with simple setup process using Docker.
- **Provides Additional Information**: Provides extra insights, such as relationships between projects and dependencies.

## Quick Start

If you want to setup with Gentoo prefix, please refer to [How to setup Gentoo prefix](./docs/setup/gentoo.md).

Then make sure `docker` and `docker-compose-v2` is installed, and run the following.

```sh
export GENTOO_PREFIX_DIR=<your Gentoo prefix location>
export GITHUB_TOKEN=<your GitHub token>
./setup.sh
```

After the script finish, try to connect to database (the 
password is stored in `data/DB_PASSWD`.

Then populate git_link fields in arch_packages and debian_packages manually and finally run following command). If you have git_link data already, you can use `scripts/copy-gitlink.py` tool to copy the data to the database.

Then, running the following command for the first time collecting and calculating the score. This will take days to finish.

```sh
docker compose exec app bash /gitlink.sh
```

## Documentation of general design, tools and separate components

See [docs/](./docs/) for details

## Reference

1. [ossf/criticality_score](https://github.com/ossf/criticality_score)
