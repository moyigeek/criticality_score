package checkvalid
import (
	"database/sql"
	"encoding/csv"
	"os"
	"time"
	"strings"
)

type Metrics struct {
	CreatedSince     time.Time
	UpdatedSince     time.Time
	Score        	 float64
}

var repoList = []string{"debain_packages", "arch_packages", "gentoo_packages", "nix_packages", "homebrew_packages"}

func fetchDistroGitlink(gitlink *sql.DB, repo string)[]string{
	query := "SELECT git_link FROM " + repo
	rows, err := gitlink.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var gitLinks []string
	for rows.Next() {
		var gitLink sql.NullString
		err := rows.Scan(&gitLink)
		if err != nil {
			panic(err)
		}
		if gitLink.Valid {
			gitLinks = append(gitLinks, gitLink.String)
		}
	}
	return gitLinks
}

func checkDistroValid(gitlink *sql.DB, repo string)[]string{
	gitLinks := fetchDistroGitlink(gitlink, repo)
	var invalidLinks []string
	for _, link := range gitLinks {
		if link == "" || link == "NA" || link == "NAN" {
			continue
		}
		if !strings.HasSuffix(link, ".git") {
			invalidLinks = append(invalidLinks, link)
		}
	}
	return invalidLinks
}

func fetchMetrics(db *sql.DB)map[string]Metrics{
	query := "SELECT git_link, created_since, updated_since, score FROM git_metrics"
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	MetricsList := make(map[string]Metrics)
	for rows.Next() {
		var gitLink string
		var createdSince time.Time
		var updatedSince time.Time
		var score float64
		err := rows.Scan(&gitLink, &createdSince, &updatedSince, &score)
		if err != nil {
			panic(err)
		}
		MetricsList[gitLink] = Metrics{CreatedSince: createdSince, UpdatedSince: updatedSince, Score: score}
	}
	return MetricsList
}

func checkMetricsValid(db *sql.DB)[]string{
	MetricsList := fetchMetrics(db)
	var invalidLinks []string
	for link, metrics := range MetricsList {
		duration := metrics.CreatedSince.Sub(metrics.UpdatedSince)
		if duration > 0 {
            invalidLinks = append(invalidLinks, link)
        }
		if metrics.Score < 0 {
			invalidLinks = append(invalidLinks, link)
		}
	}
	return invalidLinks
}

func CheckVaild(db *sql.DB)[]string{
	var invalidLinks []string
	for _, repo := range repoList {
		invalidLinks = append(invalidLinks, checkDistroValid(db, repo)...)
	}
	invalidLinks = append(invalidLinks, checkMetricsValid(db)...)
	return invalidLinks
}

func WriteCsv(invalidLinks []string, outputFile string){
	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	for _, link := range invalidLinks {
		writer.Write([]string{link})
		writer.Flush()
	}
}