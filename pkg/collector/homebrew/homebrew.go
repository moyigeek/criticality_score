package homebrew

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/lib/pq"
)

type PackageInfo struct {
	Name         string
	Description  string
	Homepage     string
	Depends      []string
	DependsCount int
	URL          string
	GitRepo      string
	PageRank     float64
}

type HomebrewCollector struct {
	RepoDir    string
	PkgInfoMap map[string]PackageInfo
}

func NewHomebrewCollector() *HomebrewCollector {
	return &HomebrewCollector{
		PkgInfoMap: make(map[string]PackageInfo),
	}
}

func (hc *HomebrewCollector) CloneHomebrewRepo() error {
	repoURL := "https://github.com/Homebrew/homebrew-core.git"
	dir := "homebrew-core"

	if _, err := os.Stat(dir); err == nil {
		cmd := exec.Command("git", "-C", dir, "pull")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to pull repository: %v", err)
		}
		hc.RepoDir = dir
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check directory: %v", err)
	}

	cmd := exec.Command("git", "clone", repoURL, dir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	hc.RepoDir = dir
	return nil
}

func (hc *HomebrewCollector) FetchAndParseFormulaFiles() error {
	if err := hc.CloneHomebrewRepo(); err != nil {
		return err
	}

	formulaDir := filepath.Join(hc.RepoDir, "Formula")

	err := filepath.Walk(formulaDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".rb") {
			formulaContent, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read formula file %s: %v", path, err)
			}

			pkgInfo := hc.parseFormulaContent(string(formulaContent), info.Name())
			if pkgInfo.Name != "" {
				pkgInfo.DependsCount = len(pkgInfo.Depends)
				hc.PkgInfoMap[pkgInfo.Name] = pkgInfo
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk through formula directory: %v", err)
	}

	return nil
}

func (hc *HomebrewCollector) parseFormulaContent(content string, fileName string) PackageInfo {
	var pkgInfo PackageInfo

	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	pkgInfo.Name = baseName

	descRe := regexp.MustCompile(`desc\s+"([^"]+)"`)
	if match := descRe.FindStringSubmatch(content); len(match) > 1 {
		pkgInfo.Description = match[1]
	}

	homepageRe := regexp.MustCompile(`homepage\s+"([^"]+)"`)
	if match := homepageRe.FindStringSubmatch(content); len(match) > 1 {
		pkgInfo.Homepage = match[1]
	}

	dependsRe := regexp.MustCompile(`depends_on\s+"([^"]+)"`)
	dependsMatches := dependsRe.FindAllStringSubmatch(content, -1)
	for _, match := range dependsMatches {
		if len(match) > 1 {
			pkgInfo.Depends = append(pkgInfo.Depends, match[1])
		}
	}

	urlRe := regexp.MustCompile(`url\s+"([^"]+)"`)
	if match := urlRe.FindStringSubmatch(content); len(match) > 1 {
		url := match[1]
		pkgInfo.URL = url
		if strings.HasPrefix(url, "https://github.com/") {
			parts := strings.Split(url, "/")
			if len(parts) >= 5 {
				orgName := parts[3]
				repoName := parts[4]
				pkgInfo.GitRepo = fmt.Sprintf("https://github.com/%s/%s.git", orgName, repoName)
			}
		} else if strings.HasPrefix(url, "https://gitlab.com/") {
			parts := strings.Split(url, "/")
			if len(parts) >= 5 {
				orgName := parts[3]
				repoName := parts[4]
				pkgInfo.GitRepo = fmt.Sprintf("https://gitlab.com/%s/%s.git", orgName, repoName)
			}
		} else if strings.HasPrefix(url, "https://gitee.com/") {
			parts := strings.Split(url, "/")
			if len(parts) >= 5 {
				orgName := parts[3]
				repoName := parts[4]
				pkgInfo.GitRepo = fmt.Sprintf("https://gitee.com/%s/%s.git", orgName, repoName)
			}
		} else if strings.HasPrefix(url, "https://bitbucket.org/") {
			parts := strings.Split(url, "/")
			if len(parts) >= 5 {
				orgName := parts[3]
				repoName := parts[4]
				pkgInfo.GitRepo = fmt.Sprintf("https://bitbucket.org/%s/%s.git", orgName, repoName)
			}
		}
	}

	return pkgInfo
}

func (hc *HomebrewCollector) getAllDep(pkgName string, visited map[string]bool, deps []string) []string {
	if visited[pkgName] {
		return deps
	}

	visited[pkgName] = true
	deps = append(deps, pkgName)

	if pkg, ok := hc.PkgInfoMap[pkgName]; ok {
		for _, depName := range pkg.Depends {
			deps = hc.getAllDep(depName, visited, deps)
		}
	}
	return deps
}

func (hc *HomebrewCollector) calculatePageRank(iterations int, dampingFactor float64) map[string]float64 {
	pageRank := make(map[string]float64)
	numPackages := len(hc.PkgInfoMap)

	for pkgName := range hc.PkgInfoMap {
		pageRank[pkgName] = 1.0 / float64(numPackages)
	}

	for i := 0; i < iterations; i++ {
		newPageRank := make(map[string]float64)

		for pkgName := range hc.PkgInfoMap {
			newPageRank[pkgName] = (1 - dampingFactor) / float64(numPackages)
		}

		for pkgName, pkgInfo := range hc.PkgInfoMap {
			var depNum int
			for _, depName := range pkgInfo.Depends {
				if _, exists := hc.PkgInfoMap[depName]; exists {
					depNum++
				}
			}
			for _, depName := range pkgInfo.Depends {
				if _, exists := hc.PkgInfoMap[depName]; exists {
					newPageRank[depName] += dampingFactor * (pageRank[pkgName] / float64(depNum))
				}
			}
		}
		pageRank = newPageRank
	}
	return pageRank
}

func (hc *HomebrewCollector) generateDependencyGraph(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("digraph {\n")

	packageIndices := make(map[string]int)
	index := 0

	for pkgName, pkgInfo := range hc.PkgInfoMap {
		packageIndices[pkgName] = index
		label := fmt.Sprintf("%s@%s", pkgName, pkgInfo.Description)
		writer.WriteString(fmt.Sprintf("  %d [label=\"%s\"];\n", index, label))
		index++
	}

	for pkgName, pkgInfo := range hc.PkgInfoMap {
		pkgIndex := packageIndices[pkgName]
		for _, depName := range pkgInfo.Depends {
			if depIndex, ok := packageIndices[depName]; ok {
				writer.WriteString(fmt.Sprintf("  %d -> %d;\n", pkgIndex, depIndex))
			}
		}
	}

	writer.WriteString("}\n")
	writer.Flush()
	return nil
}

func (hc *HomebrewCollector) storeDependenciesInDatabase(pkgName string, dependencies []string) error {
	db, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, dep := range dependencies {
		_, err := db.Exec("INSERT INTO homebrew_relationships (frompackage, topackage) VALUES ($1, $2)", pkgName, dep)
		if err != nil {
			return err
		}
	}
	return nil
}

func (hc *HomebrewCollector) updateOrInsertDatabase() error {
	db, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	for pkgName, pkgInfo := range hc.PkgInfoMap {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM homebrew_packages WHERE package = $1)", pkgName).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			_, err := db.Exec("INSERT INTO homebrew_packages (package, depends_count, description, homepage, page_rank) VALUES ($1, $2, $3, $4, $5)",
				pkgName, pkgInfo.DependsCount, pkgInfo.Description, pkgInfo.Homepage, pkgInfo.PageRank)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec("UPDATE homebrew_packages SET depends_count = $1, description = $2, homepage = $3, page_rank = $4 WHERE package = $5",
				pkgInfo.DependsCount, pkgInfo.Description, pkgInfo.Homepage, pkgInfo.PageRank, pkgName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (hc *HomebrewCollector) isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}

func (hc *HomebrewCollector) Collect(outputPath string) {
	if err := hc.FetchAndParseFormulaFiles(); err != nil {
		fmt.Printf("Error fetching package info: %v\n", err)
		return
	}

	depMap := make(map[string][]string)
	for pkgName := range hc.PkgInfoMap {
		visited := make(map[string]bool)
		deps := hc.getAllDep(pkgName, visited, []string{})
		depMap[pkgName] = deps
	}

	countMap := make(map[string]int)
	for _, deps := range depMap {
		for _, dep := range deps {
			countMap[dep]++
		}
	}

	pagerank := hc.calculatePageRank(20, 0.85)

	for pkgName, pkgInfo := range hc.PkgInfoMap {
		pagerankVal := pagerank[pkgName]
		depCount := countMap[pkgName]
		pkgInfo.PageRank = pagerankVal
		pkgInfo.DependsCount = depCount
		hc.PkgInfoMap[pkgName] = pkgInfo
	}
	err := hc.updateOrInsertDatabase()
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
	for pkgName, pkgInfo := range hc.PkgInfoMap {
		if err := hc.storeDependenciesInDatabase(pkgName, pkgInfo.Depends); err != nil {
			if hc.isUniqueViolation(err) {
				continue
			}
			fmt.Printf("Error storing dependencies for package %s: %v\n", pkgName, err)
		}
	}
	fmt.Println("Database updated successfully.")

	if outputPath != "" {
		err := hc.generateDependencyGraph(outputPath)
		if err != nil {
			fmt.Printf("Error generating dependency graph: %v\n", err)
			return
		}
		fmt.Println("Dependency graph generated successfully.")
	}
}
