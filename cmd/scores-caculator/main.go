package main

import (
	"github.com/HUSTSecLab/criticality_score/pkg/config"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	scores "github.com/HUSTSecLab/criticality_score/pkg/score"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
)

var (
	batchSize = pflag.Int("batch", 1000, "batch size")
	calcType  = pflag.String("calc", "all", "calculation type: distro, git, langeco, all")
)

func main() {
	config.RegistCommonFlags(pflag.CommandLine)
	config.ParseFlags(pflag.CommandLine)
	ac := storage.GetDefaultAppDatabaseContext()
	scores.UpdatePackageList(ac)
	linksMap := scores.FetchGitLink(ac)
	// linksMap := []string{"https://gcc.gnu.org/git/gcc.git"}
	gitMeticMap := scores.FetchGitMetrics(ac)
	langEcoMetricMap := scores.FetchLangEcoMetadata(ac)
	distMetricMap := scores.FetchDistMetadata(ac)
	var gitMetadataScore = make(map[string]*scores.GitMetadataScore)

	packageScore := make(map[string]*scores.LinkScore)

	for _, link := range linksMap {
		if _, ok := distMetricMap[link]; !ok {
			distMetricMap[link] = scores.NewDistScore()
		}
		distMetricMap[link].CalculateDistScore()

		if _, ok := langEcoMetricMap[link]; !ok {
			langEcoMetricMap[link] = scores.NewLangEcoScore()
		}
		langEcoMetricMap[link].CalculateLangEcoScore()

		gitMetadataScore[link] = scores.NewGitMetadataScore()
		if _, ok := gitMeticMap[link]; !ok {
			continue
		}
		gitMetadataScore[link].CalculateGitMetadataScore(gitMeticMap[link])
		packageScore[link] = scores.NewLinkScore(gitMetadataScore[link], distMetricMap[link], langEcoMetricMap[link])
		packageScore[link].CalculateScore()
	}
	logger.Println("Updating database...")
	scores.UpdateScore(ac, packageScore)
}
