package union_repo

import (
	"fmt"
	"strings"
	"database/sql"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

func fetchMetricsLinks() (map[string]string, error) {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT git_link FROM git_metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make(map[string]string)
	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, err
		}
		links[link] = link
	}
	return links, nil
}

func fetchUnionRepoLinks() (map[string]string, error) {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT git_link FROM git_repositories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make(map[string]string)
	for rows.Next() {
		var link string
		if err := rows.Scan(&link); err != nil {
			return nil, err
		}
		links[link] = link
	}
	return links, nil
}

const batchSize = 1000

func batchInsertLinks(links []string) error {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	for i := 0; i < len(links); i += batchSize {
		end := i + batchSize
		if end > len(links) {
			end = len(links)
		}
		if err := insertBatch(db, links[i:end]); err != nil {
			return err
		}
	}

	return nil
}

func insertBatch(db *sql.DB, links []string) error {
	valueStrings := make([]string, 0, len(links))
	valueArgs := make([]interface{}, 0, len(links))

	for i, link := range links {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i+1))
		valueArgs = append(valueArgs, link)
	}

	query := fmt.Sprintf("INSERT INTO git_repositories (git_link) VALUES %s", strings.Join(valueStrings, ","))
	_, err := db.Exec(query, valueArgs...)
	return err
}

func Run() {
	linkMetrics, err := fetchMetricsLinks()
	if err != nil {
		fmt.Println("Error fetching links from git_metrics:", err)
		return
	}
	linkUnion, err := fetchUnionRepoLinks()
	if err != nil {
		fmt.Println("Error fetching links from git_repositories:", err)
		return
	}

	newLinks := make([]string, 0)
	for link := range linkMetrics {
		if _, ok := linkUnion[link]; !ok {
			newLinks = append(newLinks, link)
		}
	}

	if len(newLinks) > 0 {
		err := batchInsertLinks(newLinks)
		if err != nil {
			fmt.Println("Error batch inserting links:", err)
			return
		}
		fmt.Printf("Successfully inserted %d new links.\n", len(newLinks))
	} else {
		fmt.Println("No new links to insert.")
	}
}
