package scores

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"
	"strings"

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
	Org_Count		 *int
}

// Define weights (αi) and max thresholds (Ti)
var weights = map[string]float64{
	// "star_count":        1,
	// "fork_count":        1,
	"created_since":     1,
	"updated_since":     -1,
	"contributor_count": 2,
	"commit_frequency":  1,
	"depsdev_ratios":    3,
	"deps_distro":       3,
	"org_count":		 1,
}

var thresholds = map[string]float64{
	// "star_count":        10000,
	// "fork_count":        5000,
	"created_since":     120, // in months
	"updated_since":     120, // in months
	"contributor_count": 40000,
	"commit_frequency":  1000,
	"depsdev_ratios":    50,
	"deps_distro":       50,
	"org_count":		 8400,
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
	err := db.QueryRow(fmt.Sprintf("SELECT COALESCE(SUM(depends_count), 0) FROM %s WHERE git_link like $1", packageType), strings.TrimSuffix(link, ".git")).Scan(&packageDependencies)
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

func CalculateScore(data ProjectData) float64 {
	score := 0.0

	var createdSinceScore, updatedSinceScore, contributorCountScore, commitFrequencyScore, Org_CountScore float64
	if data.CreatedSince != nil {
		monthsSinceCreation := time.Since(*data.CreatedSince).Hours() / (24 * 30)
		normalized := math.Log(monthsSinceCreation + 1) / math.Log(math.Max(monthsSinceCreation, thresholds["created_since"]) + 1)
		createdSinceScore = weights["created_since"] * normalized
		score += createdSinceScore
	}

	if data.UpdatedSince != nil {
		monthsSinceUpdate := time.Since(*data.UpdatedSince).Hours() / (24 * 30)
		normalized := math.Log(monthsSinceUpdate + 1) / math.Log(math.Max(monthsSinceUpdate, thresholds["updated_since"]) + 1)
		updatedSinceScore = weights["updated_since"] * normalized
		score += updatedSinceScore
	}

	if data.ContributorCount != nil {
		normalized := math.Log(float64(*data.ContributorCount) + 1) / math.Log(math.Max(float64(*data.ContributorCount), thresholds["contributor_count"]) + 1)
		contributorCountScore = weights["contributor_count"] * normalized
		score += contributorCountScore
	}

	if data.CommitFrequency != nil {
		normalized := math.Log(float64(*data.CommitFrequency) + 1) / math.Log(math.Max(float64(*data.CommitFrequency), thresholds["commit_frequency"]) + 1)
		commitFrequencyScore = weights["commit_frequency"] * normalized
		score += commitFrequencyScore
	}
	if data.Org_Count != nil {
		normalized := math.Log(float64(*data.Org_Count) + 1) / math.Log(math.Max(float64(*data.Org_Count), thresholds["org_count"]) + 1)
		Org_CountScore = weights["org_count"] * normalized
		score += Org_CountScore
	}
	if data.Pkg_Manager != nil {
		pkgManager, ok := PackageManagerData[*data.Pkg_Manager]
		if ok && data.DepsdevCount != nil {
			ratios := float64(*data.DepsdevCount)/float64(pkgManager)
			normalized := math.Log(ratios + 1) / math.Log(math.Max(ratios, thresholds["depsdev_ratios"]) + 1)
			score += weights["depsdev_ratios"] * normalized
		}
	}
	if data.deps_distro != nil {
		normalized := math.Log(float64(*data.deps_distro) + 1) / math.Log(math.Max(float64(*data.deps_distro), thresholds["deps_distro"]) + 1)
		score += weights["deps_distro"] * normalized
	}

	return score / 10
}

func UpdateDepsdistro(db *sql.DB, link string, totalRatio float64) error {
	_, err := db.Exec("UPDATE git_metrics SET deps_distro = $1 WHERE git_link = $2", totalRatio, link)
	return err
}

func UpdateScore(db *sql.DB, gitLink string, score float64) error {
	_, err := db.Exec("UPDATE git_metrics SET scores = $1 WHERE git_link = $2", score, gitLink)
	return err
}

func FetchProjectData(db *sql.DB, gitLink string) (*ProjectData, error) {
	row := db.QueryRow("SELECT created_since, updated_since, contributor_count, commit_frequency, depsdev_count, deps_distro, ecosystem, org_count FROM git_metrics WHERE git_link = $1", gitLink)
	var data ProjectData
	err := row.Scan(&data.CreatedSince, &data.UpdatedSince, &data.ContributorCount, &data.CommitFrequency, &data.DepsdevCount, &data.deps_distro, &data.Pkg_Manager, &data.Org_Count)
	if err != nil {
		log.Printf("Failed to fetch data for git link %s: %v", gitLink, err)
		return nil, err
	}
	return &data, nil
}

func CalculateDepsdistro(db *sql.DB, link string) float64{
	totalRatio := 0.0
	depRatio, err := CalculateDependencyRatio(db, link, "debian_packages")
	if err == nil {
		totalRatio += depRatio
	}
	
	depRatio, err = CalculateDependencyRatio(db, link, "arch_packages")
	if err == nil {
			totalRatio += depRatio
	}

	depRatio, err = CalculateDependencyRatio(db, link, "nix_packages")
	if err == nil {
			totalRatio += depRatio
	}

	depRatio, err = CalculateDependencyRatio(db, link, "homebrew_packages")
	if err == nil {
			totalRatio += depRatio
	}

	depRatio, err = CalculateDependencyRatio(db, link, "gentoo_packages")
	if err == nil {
			totalRatio += depRatio
	}

	err = UpdateDepsdistro(db, link, totalRatio)
	if err != nil {
		log.Printf("Failed to update database for %s: %v", link, err)
	}
	return totalRatio
}