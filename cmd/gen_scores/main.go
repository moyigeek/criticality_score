package main

import (
	"database/sql"
	"flag"
	"log"

	scores "github.com/HUSTSecLab/criticality_score/pkg/gen_scores"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	_ "github.com/lib/pq"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")

func main() {
	flag.Parse()
	storage.InitializeDatabase(*flagConfigPath)
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	links, err := fetchAllLinks(db)
	if err != nil {
		log.Fatalf("Failed to fetch git links: %v", err)
	}
	for _, link := range links{
		scores.CalculateDepsdistro(db, link)
		data, err := scores.FetchProjectData(db, link)
		if err != nil {
			log.Printf("Failed to fetch project data for %s: %v", link, err)
			continue
		}
		score := scores.CalculateScore(*data)
		if err := scores.UpdateScore(db, link, score * 100); err != nil {
			log.Printf("Failed to update score for %s: %v", link, err)
		}
	}
}

func fetchAllLinks(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT git_link FROM git_metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []string
	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}
