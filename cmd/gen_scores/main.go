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

	for _, link := range links {
		totalRatio := 0.0

		depRatio, err := scores.CalculateDependencyRatio(db, link, "debian_packages")
		if err == nil {
			totalRatio += depRatio
		}

		depRatio, err = scores.CalculateDependencyRatio(db, link, "arch_packages")
		if err == nil {
			totalRatio += depRatio
		}
		var depsdevCount int
		var pm string
		err = db.QueryRow("SELECT COALESCE(depsdev_count, 0) FROM git_metrics WHERE git_link = $1", link).Scan(&depsdevCount)
		if err == nil && depsdevCount > 0 {
			pm = scores.GetProjectTypeFromDB(link)
			if err == nil && pm != "" {
				totalPackages, ok := scores.PackageManagerData[pm]
				if ok && totalPackages > 0 {
					ratio := float64(depsdevCount) / float64(totalPackages)
					totalRatio += 10 * ratio
				}
			}
		}

		// Update the database with the computed total ratio and the detected package manager
		err = scores.UpdateDepsdistro(db, link, pm, totalRatio)
		if err != nil {
			log.Printf("Failed to update database for %s: %v", link, err)
		}
	}

	for _, link := range links {
		data, err := scores.FetchProjectData(db, link)
		if err != nil {
			log.Printf("Failed to fetch project data for %s: %v", link, err)
			continue
		}

		score := scores.CalculateScore(*data)
		if err := scores.UpdateScore(db, link, score); err != nil {
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
