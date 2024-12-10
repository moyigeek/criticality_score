package union_repo
import (
	"fmt"
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

func fetchMetricsLinks() (map[string]string, error) {
	db, err := storage.GetDatabaseConnection()
	defer db.Close()
	rows, err := db.Query("SELECT git_link from git_metrics")
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
	return links, err
}
func Updatelink(gitlink string)error {
	db, err := storage.GetDatabaseConnection()
	defer db.Close()
	_, err = db.Exec("INSERT INTO git_repositories (git_link) VALUES ($1)", gitlink)
	if err != nil {
		return err
	}
	return nil
}

func isDuplicateKeyError(err error) bool {
    if err == nil {
        return false
    }
    return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

func Run(){
	linkMetrics, err := fetchMetricsLinks()
	if err != nil {
		fmt.Println("Error fetching links:", err)
		return
	}
	for _, link := range linkMetrics {
		err := Updatelink(link)
		if err != nil{
			if isDuplicateKeyError(err) {
				continue
			}
			fmt.Println("Error updating link:", err)
		}
	}
}