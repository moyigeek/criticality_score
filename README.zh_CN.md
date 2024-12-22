# Criticality Score

[English](./README.md) | 简体中文

## 简介

Criticality Score 项目旨在基于开源项目在开源生态系统中的关键性和下载量开展关键开源项目评估和排名。与仅依赖 GitHub 指标的传统方法不同，该项目结合了来自各种 Linux 发行版的数据，以提供更全面的分析。

This project aims to evaluate and rank open-source projects based on their criticality within the open source ecosystem and download volume. Unlike traditional methods that rely solely on GitHub metrics, this project incorporates data from various Linux distributions and more code hosting platforms to provide a more comprehensive analysis.

通过收集和分析来自多个来源的指标，该项目提供了一个评估开源项目关键性的强大和全面的框架。

## 与 [ossf/criticality_score](https://github.com/ossf/criticality_score) 的区别

- **发行版依赖**：收集来自各种 Linux 发行版（例如 Debian、Arch、Nix、Gentoo）的数据，以评估开源软件的依赖性。
- **支持所有 Git 仓库**：分析来自任何 Git 平台的仓库，而不仅仅是 GitHub。
- **全面的指标收集**：从 Git 仓库和包管理器收集更广泛和更精确的指标，例如提交次数、组织数量比 GitHub API 更准确。
- **指标定制友好**: 可定制任意用于关键性评估算法的指标，而不是仅使用 Github API 中收集到的 metrics。
- **无需依赖 Google Cloud 或 BigQuery**：`ossf/criticality_score` 依赖于 Google Cloud 服务，难以迁移到其他平台。该项目独立于特定云服务运行，确保易于部署。
- **简易部署**：只需运行一个脚本，系统将通过 Docker 自动建立数据库，完成设置。
- **提供附加信息**：提供额外的信息，例如各个仓库（依赖项）之间的关系。

## 快速开始

如果您想收集 gentoo 的数据，请参阅 [如何设置 Gentoo prefix](./docs/setup/gentoo.zh_CN.md) 设置 Gentoo prefix。

然后确保已安装 `docker` 和 `docker-compose-v2`，并运行以下命令。

```sh
export GENTOO_PREFIX_DIR=<你的 Gentoo prefix 路径> # 如果没有设置 Gentoo prefix，请忽略
export GITHUB_TOKEN=<你的 GitHub Token>
./setup.sh
```

1. 启动脚本完成后，尝试连接到数据库（密码存储在 data/DB_PASSWD 中）。

2. 手动填充 arch_packages 和 debian_packages 中的 git_link 字段，最后运行以下命令。如果 git_link 数据已存在，可以使用 scripts/copy-gitlink.py 工具将数据复制到数据库中。

3. 首次运行以下命令以收集和计算分数，这将需要几天时间完成。

```sh
docker compose exec app bash /gitlink.sh
```

## 总体设计、工具和组件的文档

有关详细信息，请参阅 [docs/](./docs/)。

## 指标描述

| **指标名称**              | **描述**                                                                                 | **理由**                                                                                  |
|---------------------------|-----------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------|
| created_since（项目创建时间） | 项目创建时间越久，越有可能被广泛使用或依赖，代表成熟度和稳定性                                    | 项目创建时间久说明成熟度高，可能拥有更广泛的用户基础和依赖关系。                                        |
| updated_since（最后更新时间） | 长期未更新的项目可能不再维护，可靠性较低，依赖性下降                                               | 未维护的项目可能存在安全风险或功能不再适用，对下游用户的影响较小。                                      |
| contributor_count（贡献者数量） | 更多贡献者表明项目更受关注，社区活跃度高，重要性更强                                              | 不同贡献者的参与表明项目社区活跃，且得到广泛支持，意味着项目更可靠。                                    |
| org_count（组织数量）       | 来自不同组织的贡献者表明项目被跨组织依赖，影响范围广                                                 | 跨组织贡献者表示项目具有更高的通用性和重要性，被更多企业或机构依赖。                                    |
| commit_frequency（提交频率） | 更高的代码更改频率表明项目活跃度高，也可能存在漏洞风险                                               | 提交频繁说明项目活跃，且易于快速响应问题，但可能也需关注代码质量和安全性。                                |
| depsdev_count（依赖开发项目数量） | 表示deps.dev平台中展示各种包的依赖数目，体现其在开发生态中的重要性和影响力                                   | 更多依赖开发项目说明该项目在开发者生态中不可或缺，是多个项目的基础依赖组件。                                |
| deps_distro（发行版依赖比例）   | 反映该项目在不同软件发行版中的依赖程度，说明其在生产环境中的广泛应用和稳定性                              | 更多发行版依赖该项目表明其在软件打包和生产环境中广泛应用，具有较高的实用性和稳定性。                         |


## 引用

[1] <https://github.com/ossf/criticality_score>
