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
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, src)
	client := githubv4.NewClient(tc)

	stats := GitHubStats{}

	var repoQuery struct {
		Repository struct {
			Stargazers struct {
				TotalCount githubv4.Int
			}
			Forks struct {
				TotalCount githubv4.Int
			}
			CreatedAt    githubv4.DateTime
			UpdatedAt    githubv4.DateTime
			Contributors struct {
				TotalCount githubv4.Int
			} `graphql:"mentionableUsers"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	repoVars := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
	}

	err := client.Query(ctx, &repoQuery, repoVars)
	if err != nil {
		wait := handleRateLimitError(err)
		if wait > 0 {
			time.Sleep(wait)
			err = client.Query(ctx, &repoQuery, repoVars)
		}
		if err != nil {
			return nil
		}
	}

	starCount := int(repoQuery.Repository.Stargazers.TotalCount)
	forkCount := int(repoQuery.Repository.Forks.TotalCount)
	contributorCount := int(repoQuery.Repository.Contributors.TotalCount)

	stats.StarCount = &starCount
	stats.ForkCount = &forkCount
	stats.CreatedSince = &repoQuery.Repository.CreatedAt.Time
	stats.UpdatedSince = &repoQuery.Repository.UpdatedAt.Time
	stats.ContributorCount = &contributorCount

	var commitQuery struct {
		Repository struct {
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

	commitVars := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
		"since": githubv4.GitTimestamp{Time: time.Now().AddDate(-1, 0, 0)},
	}

	err = client.Query(ctx, &commitQuery, commitVars)
	if err != nil {
		wait := handleRateLimitError(err)
		if wait > 0 {
			time.Sleep(wait)
			err = client.Query(ctx, &commitQuery, commitVars)
		}
		if err != nil {
			return nil
		}
	}

	commitFreq := int(commitQuery.Repository.DefaultBranchRef.Target.Commit.History.TotalCount / 52)
	stats.CommitFrequency = &commitFreq

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
