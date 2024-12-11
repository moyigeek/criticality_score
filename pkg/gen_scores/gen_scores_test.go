package scores

import (
	"flag"
	"log"
	"testing"
	"fmt"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	_ "github.com/lib/pq"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")

func TestCalculateScore(t *testing.T) {
	fmt.Println("Testing CalculateScore")
	flag.Parse()
	storage.InitializeDatabase(*flagConfigPath)
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	link := "https://github.com/WangYves/Tinyblog.git"
	linkCount := make(map[string]map[string]int)
	for repo := range PackageList{
		linkCount[repo] = FetchdLinkCount(repo, db)
	}
	CalculateDepsdistro(link, linkCount)
	data, err := FetchProjectData(db, link)
	if err != nil {
		log.Printf("Failed to fetch project data for %s: %v", link, err)
		return
	}
	score := CalculateScore(*data, 0) * 100
	fmt.Println(score)
}