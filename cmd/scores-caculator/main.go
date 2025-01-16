package main

import (
	"log"

	"github.com/HUSTSecLab/criticality_score/pkg/config"
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
	gitMeticMap := scores.FetchGitMetrics(ac)
	langEcoMetricMap := scores.FetchLangEcoMetadata(ac)
	distMetricMap := scores.FetchDistMetadata(ac)

	packageScore := make(map[string]*scores.LinkScore)

	for _, link := range linksMap {
		distMetricMap[link].CalculateDistScore()
		langEcoMetricMap[link].CalculateLangEcoScore()

		gitMetadataScore := scores.NewGitMetadataScore()
		gitMetadataScore.CalculateGitMetadataScore(gitMeticMap[link])

		packageScore[link] = scores.NewLinkScore(gitMetadataScore, distMetricMap[link], langEcoMetricMap[link])
		packageScore[link].CalculateScore()
	}
	log.Println("Updating database...")
	scores.UpdateScore(ac, packageScore)
}
