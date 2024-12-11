package gitmetricsync

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	_ "github.com/lib/pq"
)

type DependentInfo struct {
	DependentCount         int `json:"dependentCount"`
	DirectDependentCount   int `json:"directDependentCount"`
	IndirectDependentCount int `json:"indirectDependentCount"`
}

var unionTables = [][]string{
	[]string{"debian_packages", "arch_packages", "gentoo_packages", "homebrew_packages", "nix_packages"},
	[]string{"github_links"},
}

func Run() {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for i := 0; i < len(unionTables); i++ {
		gitLinks := fetchGitLinks(db, i)
		syncGitMetrics(db, gitLinks, i)
	}

}

func fetchGitLinks(db *sql.DB, from int) map[string]string {
	gitLinks := make(map[string]string)
	for _, table := range unionTables[from] {
		rows, err := db.Query(fmt.Sprintf("SELECT git_link FROM %s", table))
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var gitLink sql.NullString
		for rows.Next() {
			if err := rows.Scan(&gitLink); err != nil {
				log.Fatal(err)
			}
			if gitLink.Valid{
				link := strings.TrimSpace(gitLink.String)
				if link == "" || link == "NA" || link == "NaN" {
					continue
				}
				if !strings.HasPrefix(link, "git://") && !strings.HasPrefix(link, "https://") && !strings.HasPrefix(link, "http://") {
					continue
				}
				if !strings.HasSuffix(link, ".git") {
					link += ".git"
				}
				gitLinks[strings.ToLower(link)] = link
			}
		}
	}
	return gitLinks
}

func syncGitMetrics(db *sql.DB, gitLinks map[string]string, from int) {
	normalizedLinks := make(map[string]string)
	for link := range gitLinks {
			lowercaseLink := strings.ToLower(gitLinks[link])
			normalizedLinks[lowercaseLink] = gitLinks[link]
	}
	dbLinks := make(map[string]string)
	query := `SELECT git_link FROM git_metrics WHERE "from" = $1`
	rows, err := db.Query(query, from)
	if err != nil {
		log.Fatalf("Failed to fetch git_links from git_metrics: %v", err)
	}
	defer rows.Close()

	var gitLink string
	for rows.Next() {
		if err := rows.Scan(&gitLink); err != nil {
			log.Fatalf("Failed to scan git_link from git_metrics: %v", err)
		}
		dbLinks[strings.ToLower(gitLink)] = gitLink
	}

	for dbLinkLower, dbLinkOriginal := range normalizedLinks {
		if _, exists := dbLinks[dbLinkLower]; !exists {
			if from == 0 {
				_, err := db.Exec(`
					INSERT INTO git_metrics (git_link, "from", need_update)
					VALUES ($1, $2, $3)
					ON CONFLICT (git_link) 
					DO UPDATE SET "from" = EXCLUDED."from"`,
					dbLinkOriginal, from, true)
				if err != nil {
					log.Printf("Failed to insert or update git_link %s: %v", dbLinkOriginal, err)
				}
			} else {
				_, err := db.Exec(`
					INSERT INTO git_metrics (git_link, "from", need_update)
					VALUES ($1, $2, $3)
					ON CONFLICT (git_link) DO NOTHING`,
					dbLinkOriginal, from, true)
				if err != nil {
					log.Printf("Failed to insert git_link %s: %v", dbLinkOriginal, err)
				}
			}
		}
	}
	

	for normLinkLower, normLinkOriginal := range dbLinks {
		if _, exists := normalizedLinks[normLinkLower]; !exists {
			_, err := db.Exec(`DELETE FROM git_metrics WHERE LOWER(git_link) = $1 AND "from" = $2`, normLinkLower, from)
			if err != nil {
				log.Printf("Failed to delete git_link %s: %v", normLinkOriginal, err)
			}
		}
	}
}

func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		if k != "" {
			keys = append(keys, k)
		}
	}
	return keys
}
