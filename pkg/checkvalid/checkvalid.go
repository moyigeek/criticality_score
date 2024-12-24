package checkvalid
import (
	"database/sql"
	"encoding/csv"
	"os"
	"time"
	"strings"
	"os/exec"
)

type Metrics struct {
	CreatedSince     time.Time
	UpdatedSince     time.Time
	Score        	 float64
}

var repoList = []string{"debian_packages", "arch_packages", "gentoo_packages", "nix_packages", "homebrew_packages"}

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

func checkDistroValid(gitlink *sql.DB, repo string)[][]string{
	gitLinks := fetchDistroGitlink(gitlink, repo)
	var invalidLinks [][]string
	for _, link := range gitLinks {
		if link == "" || link == "NA" || link == "NaN" {
			continue
		}
		if !strings.HasSuffix(link, ".git") {
			invalidLinks = append(invalidLinks, []string{link, "missing .git suffix"})
		}
	}
	return invalidLinks
}

func fetchMetrics(db *sql.DB)map[string]Metrics{
	query := "SELECT git_link, created_since, updated_since, scores FROM git_metrics"
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	MetricsList := make(map[string]Metrics)
	for rows.Next() {
		var gitLink string
		var createdSince sql.NullTime
		var updatedSince sql.NullTime
		var createdSinceres time.Time
		var updatedSinceres time.Time
		var score float64
		err := rows.Scan(&gitLink, &createdSince, &updatedSince, &score)
		if err != nil {
			panic(err)
		}
		if createdSince.Valid {
			createdSinceres = createdSince.Time
		}
		if updatedSince.Valid {
			updatedSinceres = updatedSince.Time
		}
		if !createdSince.Valid {
			createdSinceres = time.Time{}
		}
		if !updatedSince.Valid {
			updatedSinceres = time.Time{}
		}
		MetricsList[gitLink] = Metrics{CreatedSince: createdSinceres, UpdatedSince: updatedSinceres, Score: score}
	}
	return MetricsList
}

func checkMetricsValid(db *sql.DB)[][]string{
	MetricsList := fetchMetrics(db)
	var invalidLinks [][]string
	for link, metrics := range MetricsList {
		duration := metrics.CreatedSince.Sub(metrics.UpdatedSince)
		if duration > 0 {
            invalidLinks = append(invalidLinks, []string{link, "created_since is after updated_since"})
        }
		if metrics.Score < 0 {
			invalidLinks = append(invalidLinks, []string{link, "score is less than 0"})
		}
	}
	return invalidLinks
}

func checkCloneValid(db *sql.DB)[][]string{
	query := "SELECT git_link FROM git_metrics WHERE clone_valid = false"
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var gitLinks []string
	for rows.Next() {
		var gitLink string
		err := rows.Scan(&gitLink)
		if err != nil {
			panic(err)
		}
		gitLinks = append(gitLinks, gitLink)
	}
	var invalidLinks [][]string
	for _, link := range gitLinks {
		cmd := exec.Command("git", "clone", "--depth=1", link, "/tmp/test_repo")
		err := cmd.Run()
		if err != nil {
			invalidLinks = append(invalidLinks, []string{link, "failed to clone"})
		}
		os.RemoveAll("/tmp/test_repo")
	}
	return invalidLinks
}

func CheckVaild(db *sql.DB, checkCloneValidflag bool)[][]string{
	var invalidLinks [][]string
	for _, repo := range repoList {
		invalidLinks = append(invalidLinks, checkDistroValid(db, repo)...)
	}
	invalidLinks = append(invalidLinks, checkMetricsValid(db)...)
	if checkCloneValidflag {
		invalidLinks = append(invalidLinks, checkCloneValid(db)...)
	}
	return invalidLinks
}

func WriteCsv(invalidLinks [][]string, outputFile string){
	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	for _, value := range invalidLinks {
		link := value[0]
		reason := value[1]
		writer.Write([]string{link, reason})
		writer.Flush()
	}
}