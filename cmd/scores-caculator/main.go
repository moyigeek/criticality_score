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
	// linksMap := scores.FetchGitLink(ac)
	linksMap := []string{"git://code.call-cc.org/chicken-core"}
	gitMeticMap := scores.FetchGitMetrics(ac)
	langEcoMetricMap := scores.FetchLangEcoMetadata(ac)
	distMetricMap := scores.FetchDistMetadata(ac)
	var gitMetadataScore = make(map[string]*scores.GitMetadataScore)
	var MaxDistMetricScore, MaxLangEcoMetricScore, MaxGitMetaScore float64

	packageScore := make(map[string]*scores.LinkScore)

	for _, link := range linksMap {
		if _, ok := distMetricMap[link]; !ok {
			distMetricMap[link] = scores.NewDistScore()
		}
		distMetricMap[link].CalculateDistScore()
		if distMetricMap[link].DistScore > MaxDistMetricScore {
			MaxDistMetricScore = distMetricMap[link].DistScore
		}

		if _, ok := langEcoMetricMap[link]; !ok {
			langEcoMetricMap[link] = scores.NewLangEcoScore()
		}
		langEcoMetricMap[link].CalculateLangEcoScore()
		if langEcoMetricMap[link].LangEcoScore > MaxLangEcoMetricScore {
			MaxLangEcoMetricScore = langEcoMetricMap[link].LangEcoScore
		}

		gitMetadataScore[link] = scores.NewGitMetadataScore()
		if _, ok := gitMeticMap[link]; !ok {
			continue
		}
		gitMetadataScore[link].CalculateGitMetadataScore(gitMeticMap[link])
		if gitMetadataScore[link].GitMetadataScore > MaxGitMetaScore {
			MaxGitMetaScore = gitMetadataScore[link].GitMetadataScore
		}
	}

	for _, link := range linksMap {
		langEcoMetricMap[link].NormalizeScore()
		distMetricMap[link].NormalizeScore()
		gitMetadataScore[link].NormalizeScore()
		packageScore[link] = scores.NewLinkScore(gitMetadataScore[link], distMetricMap[link], langEcoMetricMap[link])
		packageScore[link].CalculateScore()
	}
	logger.Println("Updating database...")
	scores.UpdateScore(ac, packageScore)
}
