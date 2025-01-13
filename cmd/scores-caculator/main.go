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
	db, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	links, err := scores.FetchAllLinks(db)
	if err != nil {
		log.Fatalf("Failed to fetch git links: %v", err)
	}
	scores.CalculaterepoCount(db)
	packageScore := make(map[string]*scores.LinkScore)
	linkCount := make(map[string]map[string]scores.PackageData)
	for repo := range scores.PackageList {
		linkCount[repo] = scores.FetchdLinkCount(repo, db)
	}
	for _, link := range links {
		gitMetadata := scores.NewGitMetadata()
		distScore := scores.NewDistScore()
		langEcoScore := scores.NewLangEcoScore()
		gitMetadata.FetchGitMetadata(db, link)

		if *calcType == "distro" || *calcType == "all" {
			distScore.CalculateDistSubScore(link, linkCount)
			distScore.CalculateDistScore()
		}
		if *calcType == "git" || *calcType == "all" {
			gitMetadata.CalculateGitMetadataScore()
		}
		if *calcType == "langeco" || *calcType == "all" {
			langEcoScore.CalculateLangEcoScore()
		}

		packageScore[link] = scores.NewLinkScore(gitMetadata, distScore, langEcoScore)
		packageScore[link].CalculateScore()
	}
	log.Println("Updating database...")
	scores.UpdateScore(db, packageScore, *batchSize, *calcType)
}
