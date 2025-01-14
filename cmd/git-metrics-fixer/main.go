// This file can manual fix the git metrics of a repository
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/HUSTSecLab/criticality_score/pkg/config"
	collector "github.com/HUSTSecLab/criticality_score/pkg/gitfile/collector"
	git "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	scores "github.com/HUSTSecLab/criticality_score/pkg/score"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	gogit "github.com/go-git/go-git/v5"
	"github.com/spf13/pflag"
)

var (
	flagUpdateDB   = pflag.Bool("update-db", false, "Whether to update the database")
	flagUpdateLink = pflag.String("update-link", "", "Which link to update")
)

func main() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "This program collects metrics for git repositories.\n")
		fmt.Fprintf(os.Stderr, "This tool can be used to fix the git metrics of a repository manually.\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		pflag.PrintDefaults()
	}

	config.RegistCommonFlags(pflag.CommandLine)
	config.ParseFlags(pflag.CommandLine)
	ac := storage.GetDefaultAppDatabaseContext()

	updateDB := *flagUpdateDB
	link := *flagUpdateLink

	logger.Infof("Collecting %s", link)
	r := &gogit.Repository{}
	var err error
	u := url.ParseURL(link)
	r, err = collector.EzCollect(&u)
	if err != nil {
		logger.Panicf("Collecting %s Failed", u.URL)
	}

	repo, err := git.ParseRepo(r)
	if err != nil {
		logger.Panicf("Parsing %s Failed", link)
	}
	logger.Infof("%s Collected", repo.Name)

	repo.Show()
	gitMetadata := &scores.GitMetadata{
		CommitFrequency:  repo.CommitFrequency,
		ContributorCount: repo.ContributorCount,
		CreatedSince:     repo.CreatedSince,
		UpdatedSince:     repo.UpdatedSince,
		Org_Count:        repo.OrgCount,
	}
	gitMetadataScore := scores.NewGitMetadataScore()
	gitMetadataScore.CalculateGitMetadataScore(gitMetadata)

	distScore := scores.NewDistScore()
	distMetadata := scores.FetchDistMetadataSingle(ac, link)
	distScore.CalculateDistMerics(distMetadata[link], scores.PackageList[distMetadata[link].Type])
	distScore.CalculateDistScore()

	langEcoScore := scores.NewLangEcoScore()
	langEcoMetadata := scores.FetchLangEcoMetadataSingle(ac, link)
	langEcoScore.CalulateLangEcoMeritcs(langEcoMetadata[link], scores.PackageCounts[langEcoMetadata[link].Type])
	langEcoScore.CalculateLangEcoScore()

	if updateDB {
		scores.UpdateScore(ac, map[string]*scores.LinkScore{
			link: scores.NewLinkScore(gitMetadataScore, distScore, langEcoScore),
		})
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
