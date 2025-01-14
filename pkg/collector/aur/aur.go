package aur

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lib/pq"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var URL = "https://aur.archlinux.org/packages-meta-ext-v1.json.gz"

type PackageInfo struct {
	Depends      []string
	DependsCount int
	Description  string
	GitRepo      string
	Homepage     string
	Name         string
	PageRank     float64
	URL          string
	Version      string
}

type AurCollector struct {
	PkgInfoMap map[string]PackageInfo
	DepMap     map[string][]string
	CountMap   map[string]int
	PageRank   map[string]float64
}

func NewAurCollector() *AurCollector {
	return &AurCollector{
		PkgInfoMap: make(map[string]PackageInfo),
		DepMap:     make(map[string][]string),
		CountMap:   make(map[string]int),
		PageRank:   make(map[string]float64),
	}
}

func (ac *AurCollector) Collect(outputPath string) {
	err := ac.getDependencies()
	if err != nil {
		log.Fatal(err)
	}

	ac.calculateDependencies()
	ac.calculatePageRank(20, 0.85)
	err = ac.updateOrInsertDatabase()
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
	ac.storeDependenciesInDatabase()
	fmt.Println("Database updated successfully.")

	if outputPath != "" {
		err := ac.generateDependencyGraph(outputPath)
		if err != nil {
			fmt.Printf("Error generating dependency graph: %v\n", err)
			return
		}
		fmt.Println("Dependency graph generated successfully.")
	}
}

func (ac *AurCollector) getDependencies() error {
	resp, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var packages []PackageInfo
	err = json.Unmarshal([]byte(body), &packages)
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		ac.PkgInfoMap[pkg.Name] = pkg
	}
	return nil
}

func (ac *AurCollector) calculateDependencies() {
	for pkgName := range ac.PkgInfoMap {
		visited := make(map[string]bool)
		deps := ac.getAllDep(pkgName, visited, []string{})
		ac.DepMap[pkgName] = deps
	}

	for _, deps := range ac.DepMap {
		for _, dep := range deps {
			ac.CountMap[dep]++
		}
	}
}

func (ac *AurCollector) getAllDep(pkgName string, visited map[string]bool, deps []string) []string {
	if visited[pkgName] {
		return deps
	}

	visited[pkgName] = true
	deps = append(deps, pkgName)

	if pkg, ok := ac.PkgInfoMap[pkgName]; ok {
		for _, depName := range pkg.Depends {
			deps = ac.getAllDep(depName, visited, deps)
		}
	}
	return deps
}

func (ac *AurCollector) calculatePageRank(iterations int, dampingFactor float64) {
	numPackages := len(ac.PkgInfoMap)

	for pkgName := range ac.PkgInfoMap {
		ac.PageRank[pkgName] = 1.0 / float64(numPackages)
	}

	for i := 0; i < iterations; i++ {
		newPageRank := make(map[string]float64)

		for pkgName := range ac.PkgInfoMap {
			newPageRank[pkgName] = (1 - dampingFactor) / float64(numPackages)
		}

		for pkgName, pkgInfo := range ac.PkgInfoMap {
			var depNum int
			for _, depName := range pkgInfo.Depends {
				if _, exists := ac.PkgInfoMap[depName]; exists {
					depNum++
				}
			}
			for _, depName := range pkgInfo.Depends {
				if _, exists := ac.PkgInfoMap[depName]; exists {
					newPageRank[depName] += dampingFactor * (ac.PageRank[pkgName] / float64(depNum))
				}
			}
		}
		ac.PageRank = newPageRank
	}

	for pkgName, pkgInfo := range ac.PkgInfoMap {
		pkgInfo.PageRank = ac.PageRank[pkgName]
		pkgInfo.DependsCount = ac.CountMap[pkgName]
		ac.PkgInfoMap[pkgName] = pkgInfo
	}
}

func (ac *AurCollector) updateOrInsertDatabase() error {
	db, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	for pkgName, pkgInfo := range ac.PkgInfoMap {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM aur_packages WHERE package = $1)", pkgName).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			_, err := db.Exec("INSERT INTO aur_packages (package, depends_count, description, homepage, page_rank, version) VALUES ($1, $2, $3, $4, $5, $6)",
				pkgName, pkgInfo.DependsCount, pkgInfo.Description, pkgInfo.URL, pkgInfo.PageRank, pkgInfo.Version)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec("UPDATE aur_packages SET depends_count = $1, description = $2, homepage = $3, page_rank = $4, version = $5 WHERE package = $6",
				pkgInfo.DependsCount, pkgInfo.Description, pkgInfo.URL, pkgInfo.PageRank, pkgInfo.Version, pkgName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ac *AurCollector) storeDependenciesInDatabase() {
	db, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for pkgName, pkgInfo := range ac.PkgInfoMap {
		for _, dep := range pkgInfo.Depends {
			_, err := db.Exec("INSERT INTO aur_relationships (frompackage, topackage) VALUES ($1, $2)", pkgName, dep)
			if err != nil {
				if isUniqueViolation(err) {
					continue
				}
				fmt.Printf("Error storing dependencies for package %s: %v\n", pkgName, err)
			}
		}
	}
}

func (ac *AurCollector) generateDependencyGraph(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("digraph {\n")

	packageIndices := make(map[string]int)
	index := 0

	for pkgName, pkgInfo := range ac.PkgInfoMap {
		packageIndices[pkgName] = index
		label := fmt.Sprintf("%s@%s", pkgName, pkgInfo.Description)
		writer.WriteString(fmt.Sprintf("  %d [label=\"%s\"];\n", index, label))
		index++
	}

	for pkgName, pkgInfo := range ac.PkgInfoMap {
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

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
