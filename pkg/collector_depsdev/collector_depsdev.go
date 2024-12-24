package collector_depsdev

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"strconv"
	"time"
	"log"
	"database/sql"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	_ "github.com/lib/pq" // Assuming PostgreSQL, adjust as needed
)

type DependentInfo struct {
	DependentCount         int `json:"dependentCount"`
	DirectDependentCount   int `json:"directDependentCount"`
	IndirectDependentCount int `json:"indirectDependentCount"`
}


func updateDatabase(db *sql.DB, linkDepCountList map[string]int, batchSize int) error {
	var linkDepList [][]string
	for link, count := range linkDepCountList {
		linkDepList = append(linkDepList, []string{link, fmt.Sprintf("%d", count)})
	}
	for i := 0; i < len(linkDepList); i += batchSize {
		end := i + batchSize
		if end > len(linkDepList) {
			end = len(linkDepList)
		}
		batch := linkDepList[i:end]
		query := "UPDATE git_metrics SET depsdev_count = CASE git_link"
		args := []interface{}{}
		for j, link := range batch {
			count, _ := strconv.Atoi(link[1]) 
			query += fmt.Sprintf(" WHEN $%d THEN $%d::Integer", j*2+1, j*2+2)
			args = append(args, link[0], count)
		}
		query += " END WHERE git_link IN ("
		for j := 0; j < len(batch); j++ {
			if j > 0 {
				query += ", "
			}
			query += fmt.Sprintf("$%d", j*2+1)
		}
		query += ")"
		
		result, err := db.Exec(query, args...)
		if err != nil {
			log.Printf("Error executing query: %v", err)
			return fmt.Errorf("failed to update batch: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Error retrieving rows affected: %v", err)
			return fmt.Errorf("failed to retrieve affected rows: %w", err)
		}
		log.Printf("Batch [%d - %d]: %d rows updated", i, end, rowsAffected)
	}
	return nil
}

func Run(configPath string, batchSize int) {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		fmt.Errorf("Error initializing database: %v\n", err)
		return
	}
	defer db.Close()

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
	linkDepCountList := make(map[string]int)
	for _, link := range gitLinks {
		var repo string
		if strings.HasSuffix(link, ".git") {
			result := strings.Split(strings.TrimSuffix(link, ".git"), "/")
			repo = result[len(result)-1]
		} else {
			result := strings.Split(link, "/")
			repo = result[len(result)-1]
		}
		projectType := getProjectTypeFromDB(link)
		if projectType != "" {
			var projectTypeList []string
			var count int
			if strings.Contains(projectType, " ") {
				projectTypeList = strings.Fields(projectType)
			}else {
				projectTypeList = []string{projectType}
			}
			for _, types := range projectTypeList {
				latestVersion := getLatestVersion(repo, types)
				if latestVersion != "" {
					count += queryDepsDev(link, types, repo, latestVersion)
				}
			}
			linkDepCountList[link] = count
		}
	}
	err = updateDatabase(db, linkDepCountList, batchSize)
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error reading rows:", err)
	}
}
func getProjectTypeFromDB(link string) string {
	var projectType sql.NullString
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

	if projectType.Valid {
		return projectType.String
	}
	return ""
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

func getLatestVersion(repo, projectType string) string {
	ctx := context.Background()

	url := fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s", projectType, repo)

	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Println("Error fetching package information:", err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result PackageInfo
	json.Unmarshal(body, &result)

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

func queryDepsDev(link, projectType, projectName, version string) int{
	url := fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s/versions/%s:dependents", projectType, projectName, version)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error querying deps.dev:", err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code")
		return 0
	}

	var info DependentInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		fmt.Println("Error decoding response:", err)
		return 0
	}
	return info.DependentCount
}
