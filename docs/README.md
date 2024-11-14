# Criticality Score

## Motivation

We would like to know the most critical and most downloaded open source software in the real world. With this goal, we select a few metrics different from [ossf/criticality_score](https://github.com/ossf/criticality_score/) and rate open source projects based on them.

Differences:
1. we use Linux downstream distributions as our dataset to evaluate the dependency of open source software, instead of github mention and deps.dev;
2. we only use metrics that can be collected from git repository, instead of only Github API;

## Design

![workflow](figs/workflow.png)

TODO: modify the Chinese font in the figure into English
