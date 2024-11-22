# Criticality Score

## Motivation

We would like to know the most critical and most downloaded open source software in the real world. With this goal, we select a few metrics different from [ossf/criticality_score](https://github.com/ossf/criticality_score/) and rate open source projects based on them.

Differences:
1. we use Linux downstream distributions as our dataset to evaluate the dependency of open source software, instead of github mention and deps.dev;
2. we only use metrics that can be collected from git repository, instead of only Github API;

## Design

![workflow](figs/workflow.png)

## Overview
This project utilizes a structured workflow to automate and manage various software development tasks, primarily focused on GitHub and Git metrics collection, along with dependency management. Our automated processes are designed to enhance productivity and maintain up-to-date project metrics. Below, we detail the key components and the frequency of each task to give you a better understanding of our development operations.

## Key Tasks:
GitHub Metrics Collection: We periodically collect metrics from GitHub every 3 days to monitor project activity and performance. This helps in keeping track of contributions, issues, and other repository interactions that are crucial for project health.

Git Metrics Acquisition: Every week, we gather detailed metrics from Git to analyze code changes and repository evolution. This data is essential for assessing the progress and trends in our development practices.

Git Link Sharing: Every 6 hours, Git links are manually retrieved and shared among the team. This ensures that all team members have immediate access to the latest versions of repositories.

GitHub Enumeration: Conducted every 6 hours, this task involves a thorough enumeration of GitHub resources. It is crucial for identifying new repositories, forks, and branches that are relevant to ongoing projects.

Dependency Management using deps.dev: To maintain a robust and secure codebase, we use deps.dev every 6 hours to check and manage project dependencies. This tool helps us in identifying outdated libraries and potential security vulnerabilities.

Union Task: A quick, essential operation executed within a few minutes, the Union Task involves merging or consolidating data and resources. This operation is critical for maintaining the integrity and continuity of our codebase.

By automating these tasks, we strive to maintain high standards of efficiency and consistency in our development process. This README aims to provide a comprehensive guide to our project's workflow automation.