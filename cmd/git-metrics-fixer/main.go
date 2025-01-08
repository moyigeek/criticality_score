// This file can manual fix the git metrics of a repository
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	collector "github.com/HUSTSecLab/criticality_score/pkg/gitfile/collector"
	git "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	scores "github.com/HUSTSecLab/criticality_score/pkg/score"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/bytedance/gopkg/util/gopool"
	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/pflag"
)

var (
	flagConfigPath = pflag.String("config", "config.json", "path to the config file")
	flagUpdateDB   = pflag.Bool("update-db", false, "Whether to update the database")
)

func main() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "This program collects metrics for git repositories.\n")
		fmt.Fprintf(os.Stderr, "This tool can be used to fix the git metrics of a repository manually.\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		pflag.PrintDefaults()
	}
	pflag.Parse()
	storage.BindDefaultConfigPath("config")

	updateDB := *flagUpdateDB
	paths := flag.Args()

	var wg sync.WaitGroup
	wg.Add(len(paths))

	repos := make([]*git.Repo, 0)

	for _, path := range paths {
		gopool.Go(func() {
			defer wg.Done()
			logger.Infof("Collecting %s", path)

			r := &gogit.Repository{}
			var err error

			if strings.Contains(path, "://") {
				u := url.ParseURL(path)
				r, err = collector.EzCollect(&u)
				if err != nil {
					logger.Panicf("Collecting %s Failed", u.URL)
				}
			} else {
				r, err = collector.Open(path)
				if err != nil {
					logger.Panicf("Opening %s Failed", path)
				}
			}

			repo, err := git.ParseRepo(r)
			if err != nil {
				logger.Panicf("Parsing %s Failed", path)
			}

			repos = append(repos, repo)
			logger.Infof("%s Collected", repo.Name)
		})
	}

	wg.Wait()
	storage.BindDefaultConfigPath("config")
	db, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	for _, repo := range repos {
		repo.Show()
		depsdev_count := FetchDepsdev(db, repo.URL)
		projectData := scores.ProjectData{
			CommitFrequency:  &repo.CommitFrequency,
			ContributorCount: &repo.ContributorCount,
			CreatedSince:     &repo.CreatedSince,
			UpdatedSince:     &repo.UpdatedSince,
			Org_Count:        &repo.OrgCount,
			Pkg_Manager:      &repo.Ecosystems,
			DepsdevCount:     &depsdev_count,
		}
		linkCount := make(map[string]map[string]scores.PackageData)
		scores.CalculaterepoCount(db)
		for pkg := range scores.PackageList {
			linkCount[pkg] = make(map[string]scores.PackageData)
			count := scores.FetchdLinkCountSingle(pkg, repo.URL, db)
			linkCount[pkg][strings.ToLower(repo.URL)] = count
		}
		dist_impact, pagerank := scores.CalculateDepsdistro(repo.URL, linkCount)
		score := scores.CalculateScore(projectData, scores.LinkScore{DistroScores: dist_impact, PageRank: pagerank})
		if updateDB {
			err = updateGitMetrics(db, repo, score, dist_impact)
			if err != nil {
				logger.Fatal(err)
			}
		}
	}
}
func updateGitMetrics(db *sql.DB, repo *git.Repo, score float64, depsDistro float64) error {
	query := `
		 UPDATE git_metrics 
		 SET created_since = $1, updated_since = $2, contributor_count = $3, commit_frequency = $4, org_count = $5, scores = $6, dist_impact = $7
		 WHERE git_link = $8
	 `
	_, err := db.Exec(query, repo.CreatedSince, repo.UpdatedSince, repo.ContributorCount, repo.CommitFrequency, repo.OrgCount, score, depsDistro, repo.URL)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}

func FetchDepsdev(db *sql.DB, git_link string) int {
	query := fmt.Sprintf("SELECT depsdev_count FROM git_metrics WHERE git_link = '%s'", git_link)
	var depsdev_count sql.NullInt64
	err := db.QueryRow(query).Scan(&depsdev_count)
	if err != nil {
		log.Fatal(err)
	}
	if !depsdev_count.Valid {
		return 0
	}
	return int(depsdev_count.Int64)
}
