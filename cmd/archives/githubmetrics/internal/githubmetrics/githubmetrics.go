package githubmetrics

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Config struct {
	Database    string `json:"database"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	GitHubToken string `json:"githubToken"`
}

type GitHubStats struct {
	StarCount        *int
	ForkCount        *int
	CreatedSince     *time.Time
	UpdatedSince     *time.Time
	ContributorCount *int
	CommitFrequency  *int
	OrgCount         *int
}

type UpdateOptions struct {
	UpdateStarCount        bool
	UpdateForkCount        bool
	UpdateCreatedSince     bool
	UpdateUpdatedSince     bool
	UpdateContributorCount bool
	UpdateCommitFrequency  bool
	UpdateOrgCount         bool
	ForceUpdate            bool // 新增选项，决定是否强制更新
}

func Run(ctx context.Context, db *sql.DB, owner string, repo string, config Config, opts UpdateOptions) error {
	// 首先从数据库中查询现有的值，决定是否需要更新
	var currentStats GitHubStats
	err := db.QueryRowContext(ctx, `
		SELECT star_count, fork_count, created_since, updated_since, contributor_count, commit_frequency, org_count
		FROM git_metrics 
		WHERE git_link = $1
	`, fmt.Sprintf("https://github.com/%s/%s", owner, repo)).Scan(
		&currentStats.StarCount,
		&currentStats.ForkCount,
		&currentStats.CreatedSince,
		&currentStats.UpdatedSince,
		&currentStats.ContributorCount,
		&currentStats.CommitFrequency,
		&currentStats.OrgCount,
	)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error querying current values for %s/%s: %v", owner, repo, err)
	}

	// 如果不是强制更新且数据库中已有值，则将对应的更新选项设置为 false
	if !opts.ForceUpdate {
		if currentStats.StarCount != nil {
			opts.UpdateStarCount = false
		}
		if currentStats.ForkCount != nil {
			opts.UpdateForkCount = false
		}
		if currentStats.CreatedSince != nil {
			opts.UpdateCreatedSince = false
		}
		if currentStats.UpdatedSince != nil {
			opts.UpdateUpdatedSince = false
		}
		if currentStats.ContributorCount != nil {
			opts.UpdateContributorCount = false
		}
		if currentStats.CommitFrequency != nil {
			opts.UpdateCommitFrequency = false
		}
		if currentStats.OrgCount != nil {
			opts.UpdateOrgCount = false
		}
	}

	// 如果所有更新选项都为 false，则直接返回，不进行查询和更新
	if !opts.UpdateStarCount && !opts.UpdateForkCount && !opts.UpdateCreatedSince &&
		!opts.UpdateUpdatedSince && !opts.UpdateContributorCount &&
		!opts.UpdateCommitFrequency && !opts.UpdateOrgCount {
		fmt.Println("No updates required, skipping GitHub API queries and database update.")
		return nil
	}

	// 初始化 GitHub API 客户端
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, src)
	client := github.NewClient(tc) // 使用 v3 API 客户端来验证仓库链接

	// 检查仓库是否存在且可访问
	_, _, err = client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		if githubErr, ok := err.(*github.ErrorResponse); ok && githubErr.Response.StatusCode == 404 {
			return fmt.Errorf("repository %s/%s not found or access denied", owner, repo)
		}
		return fmt.Errorf("error checking repository %s/%s: %v", owner, repo, err)
	}

	// 初始化 GitHub v4 API 客户端
	clientV4 := githubv4.NewClient(tc)

	stats := GitHubStats{}
	if opts.UpdateStarCount || opts.UpdateForkCount || opts.UpdateCreatedSince || opts.UpdateUpdatedSince || opts.UpdateContributorCount || opts.UpdateCommitFrequency {
		// 合并后的查询结构体
		var combinedQuery struct {
			Repository struct {
				Stargazers struct {
					TotalCount githubv4.Int
				}
				Forks struct {
					TotalCount githubv4.Int
				}
				CreatedAt githubv4.DateTime
				UpdatedAt githubv4.DateTime

				MentionableUsers struct {
					TotalCount githubv4.Int
				}

				DefaultBranchRef struct {
					Target struct {
						Commit struct {
							History struct {
								TotalCount githubv4.Int
							} `graphql:"history(since: $since)"`
						} `graphql:"... on Commit"`
					}
				} `graphql:"defaultBranchRef"`
			} `graphql:"repository(owner: $owner, name: $name)"`
		}

		// 设置查询变量
		vars := map[string]interface{}{
			"owner": githubv4.String(owner),
			"name":  githubv4.String(repo),
			"since": githubv4.GitTimestamp{Time: time.Now().AddDate(-1, 0, 0)}, // 查询过去一年的提交历史
		}

		// 执行查询
		err = clientV4.Query(ctx, &combinedQuery, vars)
		if err != nil {
			fmt.Println(err)
			wait := handleRateLimitError(err)
			if wait > 0 {
				time.Sleep(wait)
				err = clientV4.Query(ctx, &combinedQuery, vars)
			}
			if err != nil {
				return err
			}
		}

		// 设置查询结果
		if opts.UpdateStarCount {
			starCount := int(combinedQuery.Repository.Stargazers.TotalCount)
			stats.StarCount = &starCount
		}
		if opts.UpdateForkCount {
			forkCount := int(combinedQuery.Repository.Forks.TotalCount)
			stats.ForkCount = &forkCount
		}
		if opts.UpdateCreatedSince {
			stats.CreatedSince = &combinedQuery.Repository.CreatedAt.Time
		}
		if opts.UpdateUpdatedSince {
			stats.UpdatedSince = &combinedQuery.Repository.UpdatedAt.Time
		}
		if opts.UpdateContributorCount {
			totalContributors := int(combinedQuery.Repository.MentionableUsers.TotalCount)
			stats.ContributorCount = &totalContributors
		}
		if opts.UpdateCommitFrequency {
			commitFreq := int(combinedQuery.Repository.DefaultBranchRef.Target.Commit.History.TotalCount / 52)
			stats.CommitFrequency = &commitFreq
		}
	}

	if opts.UpdateOrgCount {
		orgCount, err := FetchOrgCount(ctx, client, owner, repo, config.GitHubToken)
		if err != nil {
			return fmt.Errorf("error fetching organization count for %s/%s: %v", owner, repo, err)
		}
		stats.OrgCount = &orgCount
	}

	// 更新数据库
	err = updateDatabase(ctx, db, owner, repo, stats, opts)
	if err != nil {
		return fmt.Errorf("error updating database for %s/%s: %v", owner, repo, err)
	}

	return nil
}

func updateDatabase(ctx context.Context, db *sql.DB, owner, repo string, stats GitHubStats, opts UpdateOptions) error {
	query := "UPDATE git_metrics SET "
	args := []interface{}{}
	argIndex := 1

	if opts.UpdateStarCount {
		query += fmt.Sprintf("star_count = $%d, ", argIndex)
		args = append(args, stats.StarCount)
		argIndex++
	}
	if opts.UpdateForkCount {
		query += fmt.Sprintf("fork_count = $%d, ", argIndex)
		args = append(args, stats.ForkCount)
		argIndex++
	}
	if opts.UpdateCreatedSince {
		query += fmt.Sprintf("created_since = $%d, ", argIndex)
		args = append(args, stats.CreatedSince)
		argIndex++
	}
	if opts.UpdateUpdatedSince {
		query += fmt.Sprintf("updated_since = $%d, ", argIndex)
		args = append(args, stats.UpdatedSince)
		argIndex++
	}
	if opts.UpdateContributorCount {
		query += fmt.Sprintf("contributor_count = $%d, ", argIndex)
		args = append(args, stats.ContributorCount)
		argIndex++
	}
	if opts.UpdateCommitFrequency {
		query += fmt.Sprintf("commit_frequency = $%d, ", argIndex)
		args = append(args, stats.CommitFrequency)
		argIndex++
	}
	if opts.UpdateOrgCount {
		query += fmt.Sprintf("org_count = $%d, ", argIndex)
		args = append(args, stats.OrgCount)
		argIndex++
	}

	// 去掉最后一个逗号和空格
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE git_link = $%d", argIndex)
	args = append(args, fmt.Sprintf("https://github.com/%s/%s", owner, repo))

	_, err := db.Exec(query, args...)
	return err
}

func handleRateLimitError(err error) time.Duration {
	if rateLimitError, ok := err.(*github.RateLimitError); ok {
		// 检查 Reset 时间是否有效
		if !rateLimitError.Rate.Reset.IsZero() {
			resetTimestamp := rateLimitError.Rate.Reset.Time
			waitDuration := time.Until(resetTimestamp)
			if waitDuration > 0 {
				fmt.Printf("GitHub API rate limit exceeded. Waiting %v before retrying...\n", waitDuration)
				return waitDuration
			}
		}
	}
	// 如果没有找到重置时间，则休眠一小时
	fmt.Println("Unable to determine rate limit reset time. Waiting 1 hour before retrying...")
	return time.Hour
}

func FetchOrgCount(ctx context.Context, client *github.Client, owner, repo string, Token string) (int, error) {
	// 初始化组织名称过滤器
	orgFilter := strings.NewReplacer(
		"inc.", "",
		"llc", "",
		"@", "",
		" ", "",
	)

	// 设置获取贡献者的选项
	opts := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100, // 每页最多100个贡献者
		},
	}

	// 获取贡献者列表
	contributors, _, err := client.Repositories.ListContributors(ctx, owner, repo, opts)
	if err != nil {
		if wait := handleRateLimitError(err); wait > 0 {
			fmt.Printf("Waiting for %v before retrying...\n", wait)
			time.Sleep(wait)

			// 重试 ListContributors
			contributors, _, err = client.Repositories.ListContributors(ctx, owner, repo, opts)
			if err != nil {
				return 0, err // 重试后仍然出错，返回错误
			}
		} else {
			return 0, err // 其他类型的错误，直接返回
		}
	}

	if len(contributors) == 0 {
		return 0, nil // 没有有效的贡献者
	}

	// 初始化 GitHub GraphQL 客户端
	clientV4 := githubv4.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: Token})))

	// 提取和去重组织名称
	orgSet := make(map[string]struct{})
	var mu sync.Mutex // 用于保护 orgSet 的并发访问

	var wg sync.WaitGroup // 用于等待所有协程完成

	// 对每个贡献者逐个查询公司信息
	for _, contributor := range contributors {
		login := contributor.GetLogin()
		if login == "" || strings.HasSuffix(login, "[bot]") {
			continue
		}

		// 增加 WaitGroup 计数器
		wg.Add(1)

		// 启动协程并行查询每个贡献者的公司信息
		go func(login string) {
			defer wg.Done() // 协程结束时减少计数器

			for {
				// 构建 GraphQL 查询
				var query struct {
					User struct {
						Company *string
					} `graphql:"user(login: $login)"`
				}

				variables := map[string]interface{}{
					"login": githubv4.String(login),
				}

				// 执行查询
				err := clientV4.Query(ctx, &query, variables)
				if err != nil {
					waitDuration := handleRateLimitError(err)
					if waitDuration > 0 {
						// 如果触发速率限制，等待指定时间后重试
						time.Sleep(waitDuration)
						continue
					} else {
						fmt.Printf("Error querying user %s: %v\n", login, err)
						return
					}
				}

				// 处理查询结果
				if query.User.Company != nil {
					org := strings.ToLower(*query.User.Company)
					org = strings.TrimRight(orgFilter.Replace(org), ",")

					if org != "" {
						// 使用 mutex 锁保护共享资源 orgSet
						mu.Lock()
						orgSet[org] = struct{}{}
						mu.Unlock()
					}
				}

				// 如果没有错误，退出循环
				break
			}
		}(login) // 将 login 传递给协程
	}

	// 等待所有协程完成
	wg.Wait()

	// 返回唯一组织数量
	return len(orgSet), nil
}
