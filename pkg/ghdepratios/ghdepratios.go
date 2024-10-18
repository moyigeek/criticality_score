package ghdepratios

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v32/github"
)

type Config struct {
	Database    string
	User        string
	Password    string
	Host        string
	Port        string
	GitHubToken string
}

var PackageManagerData = map[string]int{
	"npm":   3400000,
	"go":    1230000,
	"maven": 636000,
	"pypi":  538000,
	"nuget": 406000,
	"cargo": 155000,
}

// FetchGitLinks retrieves GitHub links with a non-null depsdev_count.
func FetchGitLinks(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT git_link FROM git_metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []string
	var link string
	for rows.Next() {
		if err := rows.Scan(&link); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

// CalculateDependencyRatio calculates and returns the dependency ratio.
func CalculateDependencyRatio(db *sql.DB, link, packageType string) (float64, error) {
	var packageDependencies, totalPackages int
	err := db.QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(depends_count), 0) FROM %s WHERE git_link = $1", packageType), link).Scan(&packageDependencies)
	if err != nil {
		return 0.0, err
	}

	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", packageType)).Scan(&totalPackages)
	if err != nil {
		return 0.0, err
	}

	if totalPackages == 0 {
		return 0.0, nil
	}

	return float64(packageDependencies) / float64(totalPackages), nil
}

func DetectPackageManager(client *github.Client, link string) (string, error) {
	parts := strings.Split(link, "/")
	if len(parts) < 5 {
		return "", fmt.Errorf("invalid GitHub link format: %s", link)
	}
	owner, repo := parts[3], parts[4]

	return GetProjectType(client, owner, repo)
}

func GetProjectType(client *github.Client, owner, repo string) (string, error) {
attempt:
	_, dirContent, resp, err := client.Repositories.GetContents(context.Background(), owner, repo, "", nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			body := err.Error()
			re := regexp.MustCompile(`rate reset in (\d+)m(\d+)s`)
			matches := re.FindStringSubmatch(body)
			if len(matches) == 3 {
				minutes, _ := strconv.Atoi(matches[1])
				seconds, _ := strconv.Atoi(matches[2])
				waitTime := time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
				fmt.Printf("Rate limit exceeded. Waiting %v to retry...\n", waitTime)
				time.Sleep(waitTime)
				goto attempt
			}
			return "", err
		}
		return "", err
	}

	for _, file := range dirContent {
		if file.Name != nil {
			switch *file.Name {
			case "package.json":
				return "npm", nil
			case "setup.py":
				return "pypi", nil
			case "Cargo.toml":
				return "cargo", nil
			case "pom.xml":
				return "maven", nil
			case "build.gradle":
				return "gradle", nil
			case "go.mod":
				return "go", nil
			}
		}
	}
	return "", nil
}

func UpdateDatabase(db *sql.DB, link, packageManager string, totalRatio float64) error {
	_, err := db.Exec("UPDATE git_metrics SET pkg_manager = $1, deps_distro = $2 WHERE git_link = $3", packageManager, totalRatio, link)
	return err
}
