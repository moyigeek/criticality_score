package collector_depsdev

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	_ "github.com/lib/pq" // Assuming PostgreSQL, adjust as needed
)

type DependentInfo struct {
	DependentCount         int `json:"dependentCount"`
	DirectDependentCount   int `json:"directDependentCount"`
	IndirectDependentCount int `json:"indirectDependentCount"`
}

func updateDatabase(link, projectName string, dependentCount int) error {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE git_metrics SET depsdev_count = $1 WHERE git_link = $2", dependentCount, link)
	return err
}

func Run(configPath string) {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		fmt.Errorf("Error initializing database: %v\n", err)
		return
	}
	defer db.Close()

	// Query git_links from git_metrics table
	rows, err := db.Query("SELECT git_link FROM git_metrics")
	if err != nil {
		fmt.Println("Error querying git_metrics:", err)
		return
	}
	defer rows.Close()

	var gitLinks []string
	for rows.Next() {
		var gitLink string
		if err := rows.Scan(&gitLink); err != nil {
			fmt.Println("Error scanning git_link:", err)
			return
		}
		gitLinks = append(gitLinks, gitLink)
	}
	// Process each git link
	for _, link := range gitLinks {
		if strings.HasPrefix(link, "https://github.com/") || strings.HasPrefix(link, "https://gitlab.com/") || strings.HasPrefix(link, "https://bitbucket.org/") {
			parts := strings.Split(link, "/")
			if len(parts) >= 5 {
				owner := parts[3]
				repo := parts[4]
		
				// Remove the .git suffix if it exists
				if strings.HasSuffix(repo, ".git") {
					repo = strings.TrimSuffix(repo, ".git")
				}
		
				projectType := getProjectTypeFromDB(link)
				if projectType != "" {
					latestVersion := getLatestVersion(owner, repo, projectType)
					if latestVersion != "" {
						queryDepsDev(link, projectType, repo, latestVersion)
					}
				}
			}
		}		
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error reading rows:", err)
	}
}

func getProjectTypeFromDB(link string) string {
	var projectType string
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return ""
	}
	defer db.Close()
	err = db.QueryRow("SELECT ecosystem FROM git_metrics WHERE git_link = $1", link).Scan(&projectType)
	if err != nil {
		fmt.Println("Error querying project type:", err)
		return ""
	}

	return projectType
}

type VersionInfo struct {
	VersionKey struct {
		Version string `json:"version"`
	} `json:"versionKey"`
	PublishedAt time.Time `json:"publishedAt"`
}

type PackageInfo struct {
	Versions []VersionInfo `json:"versions"`
}

func getLatestVersion(owner, repo, projectType string) string {
	ctx := context.Background()

	// 构造请求URL
	url := fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s", projectType, repo)

	// 发起HTTP GET请求
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Println("Error fetching package information:", err)
		return ""
	}
	defer resp.Body.Close()

	// 解析响应体
	body, _ := io.ReadAll(resp.Body)
	var result PackageInfo
	json.Unmarshal(body, &result)

	// 寻找最新版本
	var latestVersion string
	var latestDate time.Time
	for _, version := range result.Versions {
		if version.PublishedAt.After(latestDate) {
			latestDate = version.PublishedAt
			latestVersion = version.VersionKey.Version
		}
	}

	return latestVersion
}

func queryDepsDev(link, projectType, projectName, version string) {
	url := fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s/versions/%s:dependents", projectType, projectName, version)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error querying deps.dev:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code")
		return
	}

	var info DependentInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	err = updateDatabase(link, projectName, info.DependentCount)
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
}
