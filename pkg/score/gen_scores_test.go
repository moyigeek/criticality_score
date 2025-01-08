package score

import (
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	_ "github.com/lib/pq"
	"github.com/spf13/pflag"
)

var flagConfigPath = pflag.String("config", "config.json", "path to the config file")

func TestCalculateScore(t *testing.T) {
	fmt.Println("Testing CalculateScore")
	pflag.Parse()
	storage.BindDefaultConfigPath("config")
	db, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	CalculaterepoCount(db)
	var packageScore = make(map[string]LinkScore)
	links := []string{
		"https://github.com/gcc-mirror/gcc.git",
	}
	linkCount := make(map[string]map[string]PackageData)
	for repo := range PackageList {
		linkCount[repo] = FetchdLinkCount(repo, db)
	}
	var LinkScoreMap = make(map[string][]float64)
	for _, link := range links {
		data, err := FetchProjectData(db, link)
		if err != nil {
			log.Printf("Failed to fetch project data for %s: %v", link, err)
			return
		}
		distro_scores, page_rank := CalculateDepsdistro(link, linkCount)
		LinkScoreMap[link] = append(LinkScoreMap[link], distro_scores, page_rank)
	}
	var maxDistroScore, maxPageRank, maxDepsDevCount float64
	for _, score := range LinkScoreMap {
		if len(score) < 2 {
			continue
		}
		if score[0] > maxDistroScore {
			maxDistroScore = score[0]
		}
		if score[1] > maxPageRank {
			maxPageRank = score[1]
		}
		if float64(*data.DepsdevCount) > maxDepsDevCount {
			maxDepsDevCount = float64(*data.DepsdevCount)
		}
	}
	for link, score := range LinkScoreMap {
		if len(score) < 2 {
			log.Printf("Insufficient scores for link %s", link)
			continue
		}
		distro_scores := score[0]
		page_rank := score[1]

		normalized_distro_scores := math.Log(distro_scores+1) / math.Log(maxDistroScore+1)
		normalized_page_rank := math.Log(page_rank+1) / math.Log(maxPageRank+1)
		normalized_lang_eco_impact := math.Log(float64(*data.DepsdevCount)+1) / math.Log(maxDepsDevCount+1)
		packageScore[link] = LinkScore{
			DistroScores:        normalized_distro_scores,
			PageRank:            normalized_page_rank,
			DepsdevDistroScores: normalized_lang_eco_impact,
		}
	}
	for _, link := range links {
		data, err := FetchProjectData(db, link)
		if err != nil {
			log.Printf("Failed to fetch project data for %s: %v", link, err)
			continue
		}
		score := CalculateScore(*data, packageScore[link])
		packageScore[link] = LinkScore{
			DepsdevDistroScores: packageScore[link].DepsdevDistroScores,
			PageRank:            packageScore[link].PageRank,
			DistroScores:        packageScore[link].DistroScores,
			Score:               score * 100,
		}
	}
	fmt.Println(packageScore)
}
