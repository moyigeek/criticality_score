package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/HUSTSecLab/criticality_score/pkg/ghdepratios"
	"github.com/google/go-github/v32/github"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
)

func main() {
	config := loadConfig("config.json")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.GitHubToken})
	tc := oauth2.NewClient(ctx, ts)
	gitClient := github.NewClient(tc)

	links, err := ghdepratios.FetchGitLinks(db)
	if err != nil {
		log.Fatalf("Failed to fetch git links: %v", err)
	}

	for _, link := range links {
		var totalRatio float64

		depRatio, err := ghdepratios.CalculateDependencyRatio(db, link, "debian_packages")
		if err != nil {
			log.Printf("Failed to calculate dependency ratio for Debian packages link %s: %v", link, err)
			continue
		}
		totalRatio += depRatio

		depRatio, err = ghdepratios.CalculateDependencyRatio(db, link, "arch_packages")
		if err != nil {
			log.Printf("Failed to calculate dependency ratio for Arch packages link %s: %v", link, err)
			continue
		}
		totalRatio += depRatio

		var depsdevCount int
		err = db.QueryRow("SELECT depsdev_count FROM git_metrics WHERE git_link = $1", link).Scan(&depsdevCount)
		if err != nil {
			log.Printf("Failed to fetch depsdev_count for %s: %v", link, err)
			continue
		}

		if depsdevCount > 0 {
			pm, err := ghdepratios.DetectPackageManager(gitClient, link)
			if err != nil {
				log.Printf("Failed to detect package manager for %s: %v", link, err)
				continue
			}

			if pm != "" {
				totalPackages, ok := ghdepratios.PackageManagerData[pm]
				if !ok || totalPackages == 0 {
					log.Printf("Package manager %s not recognized or has zero total packages", pm)
					continue
				}

				ratio := float64(depsdevCount) / float64(totalPackages)
				totalRatio += ratio

				ghdepratios.UpdateDatabase(db, link, pm, totalRatio)
			}
		}
	}
}

func loadConfig(path string) ghdepratios.Config {
	var config ghdepratios.Config
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	return config
}
