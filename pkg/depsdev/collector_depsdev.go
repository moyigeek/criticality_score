package depsdev

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/repository"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/samber/lo"
)

var PackageCounts = map[repository.LangEcosystemType]int{
	repository.Npm:   3.37e6,
	repository.Go:    1.29e6,
	repository.Maven: 668e3,
	repository.Pypi:  574e3,
	repository.NuGet: 430e3,
	repository.Cargo: 168e3,
}

type DependentInfo struct {
	DependentCount         int `json:"dependentCount"`
	DirectDependentCount   int `json:"directDependentCount"`
	IndirectDependentCount int `json:"indirectDependentCount"`
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

type Node struct {
	VersionKey Version  `json:"versionKey"`
	Bundled    bool     `json:"bundled"`
	Relation   string   `json:"relation"`
	Errors     []string `json:"errors"`
}

type Edge struct {
	FromNode    int    `json:"fromNode"`
	ToNode      int    `json:"toNode"`
	Requirement string `json:"requirement"`
}

type Dependencies struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
	Error string `json:"error"`
}

type Version struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type PkgInfo struct {
	VersionKey         Version
	RelationType       string   `json:"relationType"`
	RelationProvenance string   `json:"relationProvenance"`
	SlsaProvenances    []string `json:"slsaProvenances"`
	Attestations       []string `json:"attestations"`
}

type DepsDevInfo struct {
	Versions []PkgInfo `json:"versions"`
}

type EcoSystemRatio struct {
	NpmRatio   float64
	GoRatio    float64
	MavenRatio float64
	PyPiRatio  float64
	NuGetRatio float64
	CargoRatio float64
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

func queryDepsDev(projectType, projectName, version string) int {
	url := fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s/versions/%s:dependents", projectType, projectName, version)
	resp, err := http.Get(url)
	if err != nil {
		version = getLatestVersion(projectName, projectType)
		url = fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s/versions/%s:dependents", projectType, projectName, version)
		resp, err = http.Get(url)
		if err != nil {
			fmt.Println("Error fetching package information:", err)
			return 0
		}
		defer resp.Body.Close()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	var info DependentInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		fmt.Println("Error decoding response:", err)
		return 0
	}
	return info.DependentCount
}

func getGitlink(db *sql.DB) []string {
	rows, err := db.Query("SELECT git_link FROM git_metrics")
	if err != nil {
		fmt.Println("Error querying git_metrics:", err)
		return nil
	}
	defer rows.Close()
	var gitLinks []string
	for rows.Next() {
		var gitLink string
		if err := rows.Scan(&gitLink); err != nil {
			fmt.Println("Error scanning git_link:", err)
			return nil
		}
		gitLinks = append(gitLinks, gitLink)
	}
	return gitLinks
}

func queryDepsName(gitlink string, rdb *redis.Client) map[string]Version {
	depMap := make(map[string]Version)
	if strings.Contains(gitlink, ".git") {
		gitlink = strings.TrimSuffix(gitlink, ".git")
	}
	var repo, name string
	if len(strings.Split(gitlink, "/")) == 5 {
		repo = strings.Split(gitlink, "/")[3]
		name = strings.Split(gitlink, "/")[4]
	}
	url := fmt.Sprintf("https://api.deps.dev/v3alpha/projects/github.com%%2f%s%%2f%s:packageversions", repo, name)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error querying deps.dev:", err)
		return depMap
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return depMap
	}
	var result DepsDevInfo

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return depMap
	}

	latestVersions := make(map[string]map[string]string)
	for _, item := range result.Versions {
		name := item.VersionKey.Name
		version := item.VersionKey.Version
		system := item.VersionKey.System
		if strings.Contains(name, "\u003E") {
			split := strings.Split(name, "\u003E")
			name = split[len(split)-1]
		}
		if strings.Contains(name, "/") {
			split := strings.Split(name, "/")
			name = split[len(split)-1]
		}
		if _, exists := latestVersions[name]; !exists {
			latestVersions[name] = make(map[string]string)
		}
		if currentVersion, exists := latestVersions[name][system]; !exists || version > currentVersion {
			latestVersions[name][system] = version
			depMap[name] = Version{Name: name, System: system, Version: version}
			storage.SetKeyValue(rdb, name, gitlink)
		}
	}
	return depMap
}

type GitMetrics struct {
	LangEcoImpact   float64
	LangEcoPageRank float64
}

func Depsdev(batchSize int, workerPoolSize int, calculatePageRankFlag bool) {
	ac := storage.GetDefaultAppDatabaseContext()
	repo := repository.NewLangEcoLinkRepository(ac)
	rdb, _ := storage.InitRedis()
	gitLinks := fetchGitLink(ac)
	// gitLinks := []string{"https://github.com/facebook/react.git"}
	pkgMap := make(map[string][]Version)
	pkgDepMap := make(map[string]map[string]int)
	for _, gitlink := range gitLinks {
		depMap := queryDepsName(gitlink, rdb)
		for pkgName, pkgInfo := range depMap {
			if _, exists := pkgDepMap[pkgInfo.System]; !exists {
				pkgDepMap[pkgInfo.System] = make(map[string]int)
			}
			pkgDepMap[pkgInfo.System][pkgName] = queryDepsDev(pkgInfo.System, pkgInfo.Name, pkgInfo.Version)
		}
		if calculatePageRankFlag {
			pkgdepMap := fetchDep(depMap, workerPoolSize)
			for pkgName, pkgInfo := range pkgdepMap {
				pkgMap[pkgName] = pkgInfo
			}
			storage.PersistData(rdb)
		}
	}
	var pageRank map[string]float64
	if calculatePageRankFlag {
		pageRank = calculatePageRank(pkgMap, 100, 0.85)
	} else {
		pageRank = make(map[string]float64)
		for _, pkgMap := range pkgDepMap {
			for pkgName := range pkgMap {
				pageRank[pkgName] = 0.0
			}
		}
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, workerPoolSize)
	type langEcoKey struct {
		gitLink string
		ltype   repository.LangEcosystemType
	}

	langEco := make(map[langEcoKey]int)

	for system, pkgMap := range pkgDepMap {
		for pkgName := range pkgMap {
			wg.Add(1)
			semaphore <- struct{}{}
			go func(system, pkgName string) {
				defer wg.Done()
				defer func() { <-semaphore }()
				gitlink, err := storage.GetKeyValue(rdb, pkgName)
				if err != nil {
					fmt.Println("Error getting key:", err)
					return
				}

				var ltype repository.LangEcosystemType
				switch strings.ToLower(system) {
				case "cargo":
					ltype = repository.Cargo
				case "go":
					ltype = repository.Go
				case "maven":
					ltype = repository.Maven
				case "npm":
					ltype = repository.Npm
				case "nuget":
					ltype = repository.NuGet
				case "pypi":
					ltype = repository.Pypi
				}

				key := langEcoKey{
					gitLink: gitlink,
					ltype:   ltype,
				}

				mu.Lock()

				if _, exists := langEco[key]; !exists {
					langEco[key] = pkgDepMap[system][pkgName]
				} else {
					langEco[key] += pkgDepMap[system][pkgName]
				}

				mu.Unlock()
			}(system, pkgName)
		}
	}
	wg.Wait()
	var toUpdateList []*repository.LangEcosystem
	for key, info := range langEco {
		toUpdateList = append(toUpdateList, lo.ToPtr(repository.LangEcosystem{
			GitLink:       lo.ToPtr(key.gitLink),
			Type:          lo.ToPtr(key.ltype),
			DepCount:      lo.ToPtr(info),
			LangEcoImpact: lo.ToPtr(float64(info) / float64(PackageCounts[key.ltype])),
		}))
	}
	err := repo.BatchInsertOrUpdate(toUpdateList)
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
	}
}

func fetchDep(depMap map[string]Version, threadnum int) map[string][]Version {
	depMapNew := make(map[string][]Version)
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, threadnum)
	var mu sync.Mutex
	for depName, depInfo := range depMap {
		wg.Add(1)
		semaphore <- struct{}{}
		go func(depName string, depInfo Version) {
			defer wg.Done()
			defer func() { <-semaphore }()
			var system, name, version string
			system = depInfo.System
			name = depName
			version = depInfo.Version
			result := getAndProcessDependencies(system, name, version)
			mu.Lock()
			depMapNew[name] = []Version{}
			for _, node := range result.Nodes {
				if node.Relation == "DIRECT" {
					depMapNew[name] = append(depMapNew[name], node.VersionKey)
				}
			}
			mu.Unlock()
		}(depName, depInfo)
	}
	wg.Wait()
	return depMapNew
}

func calculatePageRank(pkgInfoMap map[string][]Version, iterations int, dampingFactor float64) map[string]float64 {
	pageRank := make(map[string]float64)
	numPackages := len(pkgInfoMap)

	for pkgName := range pkgInfoMap {
		pageRank[pkgName] = 1.0 / float64(numPackages)
	}

	for i := 0; i < iterations; i++ {
		newPageRank := make(map[string]float64)

		for pkgName := range pkgInfoMap {
			newPageRank[pkgName] = (1 - dampingFactor) / float64(numPackages)
		}

		for pkgName, deps := range pkgInfoMap {
			depNum := len(deps)
			for _, depName := range deps {
				if _, exists := pkgInfoMap[depName.Name]; exists {
					newPageRank[depName.Name] += dampingFactor * (pageRank[pkgName] / float64(depNum))
				}
			}
		}
		pageRank = newPageRank
	}
	return pageRank
}

func getAndProcessDependencies(system, name, version string) Dependencies {
	var result Dependencies
	url := fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s/versions/%s:dependencies", system, name, version)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error querying deps.dev:", err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		version = getLatestVersion(name, system)
		url = fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s/versions/%s:dependencies", system, name, version)
		resp, err = http.Get(url)
		if err != nil {
			fmt.Println("Error querying deps.dev:", err)
			return result
		}
		defer resp.Body.Close()
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return result
	}
	cleanedBody := removeInvisibleChars(string(body))
	err = json.Unmarshal([]byte(cleanedBody), &result)
	if err != nil {
		return result
	}

	return result
}

func removeInvisibleChars(input string) string {
	re := regexp.MustCompile(`[[:cntrl:]]+`)
	return re.ReplaceAllString(input, "")
}

func fetchGitLink(ac storage.AppDatabaseContext) []string {
	repo := repository.NewAllGitLinkRepository(ac)
	linksIter, err := repo.Query()
	if err != nil {
		log.Fatalf("Failed to fetch git links: %v", err)
	}
	links := []string{}
	for link := range linksIter {
		links = append(links, link)
	}
	return links
}
