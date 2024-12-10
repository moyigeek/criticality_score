# Criticality Score

[English](./README.md) | 简体中文

## 描述

Criticality Score 项目旨在评估和排名开源项目在生态系统中的重要性和维护状态。与仅依赖 GitHub 指标的传统方法不同，该项目结合了来自各种 Linux 发行版的数据，以提供更全面的分析。

通过收集和分析来自多个来源的指标，该项目提供了一个评估开源项目关键性的强大和全面的框架。

## 与 [ossf/criticality_score](https://github.com/ossf/criticality_score) 的区别

- **发行版依赖**：收集来自各种 Linux 发行版（例如 Debian、Nix、Gentoo）的数据，以评估开源软件的依赖性。
- **支持所有 Git 仓库**：分析来自任何 Git 平台的仓库，而不仅仅是 GitHub。
- **全面的指标收集**：从 Git 仓库和包管理器收集更广泛和更精确的指标，例如提交次数、组织数量比 GitHub API 更准确。
- **无需依赖 Google Cloud 或 BigQuery**：`ossf/criticality_score` 依赖于 Google Cloud 服务，难以迁移到其他平台。该项目独立于特定云服务运行，确保易于部署。
- **简易部署**：只需运行一个脚本，系统将自动建立数据库，启动 Docker，完成设置。
- **提供附加信息**：提供额外的信息，例如项目和依赖关系之间的关系。

## 快速开始

如果您想使用 Gentoo 前缀进行设置，请参阅 [如何设置 Gentoo prefix](./docs/setup/gentoo.zh_CN.md)。

然后确保已安装 `docker` 和 `docker-compose-v2`，并运行以下命令。

```sh
export GENTOO_PREFIX_DIR=<你的 Gentoo prefix 路径>
export GITHUB_TOKEN=<你的 GitHub Token>
./setup.sh
```

脚本完成后，尝试连接到数据库（密码存储在 data/DB_PASSWD 中）。

然后手动填充 arch_packages 和 debian_packages 中的 git_link 字段，最后运行以下命令。如果您已经有 git_link 数据，可以使用 scripts/copy-gitlink.py 工具将数据复制到数据库中。

然后，首次运行以下命令以收集和计算分数。这将需要几天时间完成。

```sh
docker compose exec app bash /gitlink.sh
```

## 总体设计、工具和独立组件的文档

有关详细信息，请参阅 [docs/](./docs/)。

## 引用

1. [ossf/criticality_score](https://github.com/ossf/criticality_score)