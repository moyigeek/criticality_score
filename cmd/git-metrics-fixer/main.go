// This file can manual fix the git metrics of a repository
package main

import (
	"fmt"
	"os"

	"github.com/HUSTSecLab/criticality_score/pkg/config"
	collector "github.com/HUSTSecLab/criticality_score/pkg/gitfile/collector"
	git "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	scores "github.com/HUSTSecLab/criticality_score/pkg/score"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/repository"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/sqlutil"
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
	gitMetric := &repository.GitMetric{
		GitLink:          &link,
		CommitFrequency:  sqlutil.ToNullable(repo.CommitFrequency),
		ContributorCount: sqlutil.ToNullable(repo.ContributorCount),
		CreatedSince:     sqlutil.ToNullable(repo.CreatedSince),
		UpdatedSince:     sqlutil.ToNullable(repo.UpdatedSince),
		OrgCount:         sqlutil.ToNullable(repo.OrgCount),
	}
	gitMetadata := InsertGitMeticAndFetch(ac, gitMetric)

	gitMetadataScore := scores.NewGitMetadataScore()
	gitMetadataScore.CalculateGitMetadataScore(gitMetadata[link])

	distScore := scores.FetchDistMetadataSingle(ac, link)
	distScore[link].CalculateDistScore()
	langEcoScore := scores.FetchLangEcoMetadataSingle(ac, link)
	langEcoScore[link].CalculateLangEcoScore()

	if updateDB {
		scores.UpdateScore(ac, map[string]*scores.LinkScore{
			link: scores.NewLinkScore(gitMetadataScore, distScore[link], langEcoScore[link]),
		})
	}
}

func InsertGitMeticAndFetch(ac storage.AppDatabaseContext, gitMetadata *repository.GitMetric) map[string]*scores.GitMetadata {
	repo := repository.NewGitMetricsRepository(ac)
	repo.InsertOrUpdate(gitMetadata)
	gitMetric := scores.FetchGitMetricsSingle(ac, *gitMetadata.GitLink)
	return gitMetric
}
