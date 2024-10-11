package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	scores "github.com/HUSTSecLab/criticality_score/pkg/gen_scores"
	_ "github.com/lib/pq"
)

type Config struct {
	Database    string `json:"database"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	GitHubToken string `json:"githubToken"`
}

func main() {
	config := loadConfig("config.json")
	db := setupDatabase(config)
	defer db.Close()

	links, err := fetchAllLinks(db)
	if err != nil {
		log.Fatalf("Failed to fetch git links: %v", err)
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

func loadConfig(path string) *Config {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	return &config
}

func setupDatabase(config *Config) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	return db
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
