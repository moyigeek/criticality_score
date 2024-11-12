package gitmetricsync

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
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
	[]string{"debian_packages", "arch_packages"},
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

func fetchGitLinks(db *sql.DB, from int) map[string]bool {
	gitLinks := make(map[string]bool)
	githubRegex := regexp.MustCompile(`https?://github\.com/[^\s/]+/[^\s/]+`)
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
			if gitLink.Valid && githubRegex.MatchString(gitLink.String) {
				gitLinks[gitLink.String] = true
			}
		}
	}
	return gitLinks
}

func syncGitMetrics(db *sql.DB, gitLinks map[string]bool, from int) {
	githubRegex := regexp.MustCompile(`https?://github\.com/([^/\s]+/[^/\s]+)`)

	// 规范化输入链接：保留原始大小写但根据情况添加或不添加 .git，存储用于插入；用小写化链接用于比较
	normalizedLinks := make(map[string]string) // key 为小写化的链接，用于比较；值为原始大小写的链接
	for link := range gitLinks {
		if matches := githubRegex.FindStringSubmatch(link); len(matches) > 1 {
			originalLink := matches[0] // 使用捕获到的完整链接
			if !strings.HasSuffix(originalLink, ".git") {
				originalLink += ".git"
			}
			lowercaseLink := strings.ToLower(originalLink)
			normalizedLinks[lowercaseLink] = originalLink
		}
	}

	// 获取数据库中所有 git_links 的当前状态，用于比较（小写化比较）
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

	// 检查哪些需要删除
	for dbLinkLower, dbLinkOriginal := range dbLinks {
		if _, exists := normalizedLinks[dbLinkLower]; !exists {
			_, err := db.Exec(`DELETE FROM git_metrics WHERE LOWER(git_link) = $1 AND "from" = $2`, dbLinkLower, from)
			if err != nil {
				log.Printf("Failed to delete git_link %s: %v", dbLinkOriginal, err)
			}
		}
	}

	// 检查哪些需要添加
	for normLinkLower, normLinkOriginal := range normalizedLinks {
		if _, exists := dbLinks[normLinkLower]; !exists {
			parts := strings.Split(normLinkOriginal, "/")
			if len(parts) >= 5 {
				_, err := db.Exec(`INSERT INTO git_metrics (git_link, "from", need_update) VALUES ($1, $2, $3)`, normLinkOriginal, from, true)
				if err != nil {
					log.Printf("Failed to insert git_link %s: %v", normLinkOriginal, err)
				}
			}
		}
	}
}

func getKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		if k != "" { // 确保不处理空字符串
			keys = append(keys, k)
		}
	}
	return keys
}
