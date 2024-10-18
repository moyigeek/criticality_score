package scores

import (
	"database/sql"
	"log"
	"math"
	"fmt"
	"time"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

type ProjectData struct {
	StarCount        *int
	ForkCount        *int
	CreatedSince     *time.Time
	UpdatedSince     *time.Time
	ContributorCount *int
	CommitFrequency  *float64
	DepsdevCount     *int
	deps_distro      *float64
	Pkg_Manager      *string
}

// Define weights (αi) and max thresholds (Ti)
var weights = map[string]float64{
	// "star_count":        1,
	// "fork_count":        1,
	"created_since":     1,
	"updated_since":     -1,
	"contributor_count": 2,
	"commit_frequency":  1,
	"depsdev_ratios":    2,
	"deps_distro":       1,
}

var thresholds = map[string]float64{
	// "star_count":        10000,
	// "fork_count":        5000,
	"created_since":     120, // in months
	"updated_since":     120, // in months
	"contributor_count": 5000,
	"commit_frequency":  1000,
	"depsdev_ratios":    40,
	"deps_distro":       50,
}

var PackageManagerData = map[string]int{
	"npm":   3400000,
	"go":    1230000,
	"maven": 636000,
	"pypi":  538000,
	"nuget": 406000,
	"cargo": 155000,
}

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

func GetProjectTypeFromDB(link string) string {
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

// CalculateScore calculates the criticality score for a project.
func CalculateScore(data ProjectData) float64 {
	score := 0.0

	// Calculate each parameter's contribution to the score.
	// if data.StarCount != nil {
	// 	normalized := math.Min(float64(*data.StarCount)/thresholds["star_count"], 1)
	// 	score += weights["star_count"] * normalized
	// }

	// if data.ForkCount != nil {
	// 	normalized := math.Min(float64(*data.ForkCount)/thresholds["fork_count"], 1)
	// 	score += weights["fork_count"] * normalized
	// }

	if data.CreatedSince != nil {
		monthsSinceCreation := time.Since(*data.CreatedSince).Hours() / (24 * 30)
		normalized := math.Min(monthsSinceCreation/thresholds["created_since"], 1)
		score += weights["created_since"] * normalized
	}

	if data.UpdatedSince != nil {
		monthsSinceUpdate := time.Since(*data.UpdatedSince).Hours() / (24 * 30)
		normalized := math.Min(monthsSinceUpdate/thresholds["updated_since"], 1)
		score += weights["updated_since"] * normalized
	}

	if data.ContributorCount != nil {
		normalized := math.Min(float64(*data.ContributorCount)/thresholds["contributor_count"], 1)
		score += weights["contributor_count"] * normalized
	}

	if data.CommitFrequency != nil {
		normalized := math.Min(*data.CommitFrequency/thresholds["commit_frequency"], 1)
		score += weights["commit_frequency"] * normalized
	}
	if data.Pkg_Manager != nil {
		// 确保包管理器的值是有效的
		pkgManager, ok := PackageManagerData[*data.Pkg_Manager]
		if ok && data.DepsdevCount != nil {
			normalized := math.Min(float64(*data.DepsdevCount)/float64(pkgManager)/thresholds["depsdev_ratios"], 1)
			score += weights["depsdev_ratios"] * normalized
		}
	}
	if data.deps_distro != nil {
		normalized := math.Min((*data.deps_distro*100)/thresholds["deps_distro"], 1)
		score += weights["deps_distro"] * normalized
	}

	return score / 6
}

func UpdateDepsdistro(db *sql.DB, link, packageManager string, totalRatio float64) error {
	_, err := db.Exec("UPDATE git_metrics SET deps_distro = $1 WHERE git_link = $2", totalRatio, link)
	return err
}

// UpdateScore updates the criticality score in the database for a given project.
func UpdateScore(db *sql.DB, gitLink string, score float64) error {
	_, err := db.Exec("UPDATE git_metrics SET scores = $1 WHERE git_link = $2", score, gitLink)
	return err
}

// FetchProjectData retrieves the project data from the database.
func FetchProjectData(db *sql.DB, gitLink string) (*ProjectData, error) {
	row := db.QueryRow("SELECT created_since, updated_since, contributor_count, commit_frequency, depsdev_count, deps_distro, ecosystem FROM git_metrics WHERE git_link = $1", gitLink)
	var data ProjectData
	err := row.Scan(&data.CreatedSince, &data.UpdatedSince, &data.ContributorCount, &data.CommitFrequency, &data.DepsdevCount, &data.deps_distro, &data.Pkg_Manager)
	if err != nil {
		log.Printf("Failed to fetch data for git link %s: %v", gitLink, err)
		return nil, err
	}
	return &data, nil
}
