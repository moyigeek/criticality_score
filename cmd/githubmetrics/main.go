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
	// 读取配置文件路径和更新选项的命令行参数
	configPath := flag.String("config", "config.json", "Path to config file")
	updateAll := flag.Bool("updateAll", false, "Update all fields")
	updateStarCount := flag.Bool("star", false, "Update star count")
	updateForkCount := flag.Bool("fork", false, "Update fork count")
	updateCreatedSince := flag.Bool("created", false, "Update created since date")
	updateUpdatedSince := flag.Bool("updated", false, "Update updated since date")
	updateContributorCount := flag.Bool("contributors", false, "Update contributor count")
	updateCommitFrequency := flag.Bool("commitfreq", false, "Update commit frequency")
	updateOrgCount := flag.Bool("orgcount", false, "Update unique organization count")

	flag.Parse()

	config := readConfig(*configPath)
	ctx := context.Background()

	// 打开数据库连接
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

		// 设置更新选项
		opts := githubmetrics.UpdateOptions{
			UpdateStarCount:        *updateAll || *updateStarCount,
			UpdateForkCount:        *updateAll || *updateForkCount,
			UpdateCreatedSince:     *updateAll || *updateCreatedSince,
			UpdateUpdatedSince:     *updateAll || *updateUpdatedSince,
			UpdateContributorCount: *updateAll || *updateContributorCount,
			UpdateCommitFrequency:  *updateAll || *updateCommitFrequency,
			UpdateOrgCount:         *updateAll || *updateOrgCount,
		}

		// 执行更新
		if err := githubmetrics.Run(ctx, db, owner, repo, config, opts); err != nil {
			log.Printf("Failed to update metrics for %s/%s: %v", owner, repo, err)
		}
	}
}

// 读取配置文件
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

// 获取GitHub链接列表
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
