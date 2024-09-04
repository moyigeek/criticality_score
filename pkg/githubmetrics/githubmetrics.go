package githubmetrics

import (
	"context"
	"database/sql"
	"fmt"
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
}

func Run(ctx context.Context, db *sql.DB, owner string, repo string, config Config) error {
	// 初始化 GitHub API 客户端
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, src)
	client := github.NewClient(tc) // 使用 v3 API 客户端来验证仓库链接

	// 检查仓库是否存在且可访问
	_, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		if githubErr, ok := err.(*github.ErrorResponse); ok && githubErr.Response.StatusCode == 404 {
			return fmt.Errorf("repository %s/%s not found or access denied", owner, repo)
		}
		return fmt.Errorf("error checking repository %s/%s: %v", owner, repo, err)
	}

	// 初始化 GitHub v4 API 客户端
	clientV4 := githubv4.NewClient(tc)

	stats := GitHubStats{}

	// 合并后的查询结构体
	var combinedQuery struct {
		Repository struct {
			// 查询 StarCount 和 ForkCount
			Stargazers struct {
				TotalCount githubv4.Int
			}
			Forks struct {
				TotalCount githubv4.Int
			}
			CreatedAt githubv4.DateTime
			UpdatedAt githubv4.DateTime

			// 查询 Contributor 总数（分页查询）
			MentionableUsers struct {
				TotalCount githubv4.Int
				Edges      []struct {
					Cursor githubv4.String
				}
			} `graphql:"mentionableUsers(first: 100, after: $cursor)"`

			// 查询 CommitFrequency
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
		"owner":  githubv4.String(owner),
		"name":   githubv4.String(repo),
		"cursor": (*githubv4.String)(nil),                                   // 初始化为 nil，用于分页查询
		"since":  githubv4.GitTimestamp{Time: time.Now().AddDate(-1, 0, 0)}, // 查询过去一年的提交历史
	}

	// 执行查询
	err = clientV4.Query(ctx, &combinedQuery, vars)
	if err != nil {
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
	starCount := int(combinedQuery.Repository.Stargazers.TotalCount)
	forkCount := int(combinedQuery.Repository.Forks.TotalCount)
	stats.StarCount = &starCount
	stats.ForkCount = &forkCount
	stats.CreatedSince = &combinedQuery.Repository.CreatedAt.Time
	stats.UpdatedSince = &combinedQuery.Repository.UpdatedAt.Time

	// 处理 Contributors 的总数
	totalContributors := int(combinedQuery.Repository.MentionableUsers.TotalCount)
	stats.ContributorCount = &totalContributors

	// 设置 Commit Frequency
	commitFreq := int(combinedQuery.Repository.DefaultBranchRef.Target.Commit.History.TotalCount / 52)
	stats.CommitFrequency = &commitFreq

	// 更新数据库
	err = updateDatabase(ctx, db, owner, repo, stats)
	if err != nil {
		return fmt.Errorf("error updating database for %s/%s: %v", owner, repo, err)
	}

	return nil
}

func handleRateLimitError(err error) time.Duration {
	if rateLimitError, ok := err.(*github.RateLimitError); ok {
		resetTimestamp := rateLimitError.Rate.Reset.Time
		waitDuration := time.Until(resetTimestamp)
		fmt.Printf("GitHub API rate limit exceeded. Waiting %v before retrying...\n", waitDuration)
		return waitDuration
	}
	return 0
}

func updateDatabase(ctx context.Context, db *sql.DB, owner, repo string, stats GitHubStats) error {
	_, err := db.Exec(`UPDATE git_metrics SET star_count = $1, fork_count = $2, created_since = $3, updated_since = $4, contributor_count = $5, commit_frequency = $6 WHERE git_link = $7`,
		stats.StarCount, stats.ForkCount, stats.CreatedSince, stats.UpdatedSince, stats.ContributorCount, stats.CommitFrequency, fmt.Sprintf("https://github.com/%s/%s", owner, repo))
	return err
}
