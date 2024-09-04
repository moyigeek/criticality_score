package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/githubmetrics"
	_ "github.com/lib/pq"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	config := readConfig(*configPath)
	ctx := context.Background()

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 获取所有git_links
	links, err := fetchGitLinks(db)
	if err != nil {
		log.Fatalf("Failed to fetch git links: %v", err)
	}

	// 遍历git_links并更新它们的统计信息
	for _, link := range links {
		parts := strings.Split(link, "/")
		if len(parts) < 5 {
			log.Printf("Invalid git link format: %s", link)
			continue
		}
		owner := parts[3]
		repo := parts[4]

		if err := githubmetrics.Run(ctx, db, owner, repo, config); err != nil {
			log.Println("Failed to update metrics for %s/%s: %v", owner, repo, err)
		}
	}
}

func readConfig(path string) githubmetrics.Config {
	var config githubmetrics.Config
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	return config
}

func fetchGitLinks(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT git_link FROM git_metrics WHERE git_link LIKE 'https://github.com/%'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []string
	var gitLink string
	for rows.Next() {
		if err := rows.Scan(&gitLink); err != nil {
			return nil, err
		}
		links = append(links, gitLink)
	}
	return links, nil
}
