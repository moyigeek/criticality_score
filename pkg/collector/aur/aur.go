package aur

import (
	"bufio"
	"bytes"
	"compress/gzip"
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

func Aur(outputPath string) {
	pkgInfoMap, err := getDependencies()
	if err != nil {
		log.Fatal(err)
	}

	depMap := make(map[string][]string)
	for pkgName := range pkgInfoMap {
		visited := make(map[string]bool)
		deps := getAllDep(pkgInfoMap, pkgName, visited, []string{})
		depMap[pkgName] = deps
	}

	countMap := make(map[string]int)
	for _, deps := range depMap {
		for _, dep := range deps {
			countMap[dep]++
		}
	}

	pagerank := calculatePageRank(pkgInfoMap, 20, 0.85)
	for pkgName, pkgInfo := range pkgInfoMap {
		pagerankVal := pagerank[pkgName]
		depCount := countMap[pkgName]
		pkgInfo.PageRank = pagerankVal
		pkgInfo.DependsCount = depCount
		pkgInfoMap[pkgName] = pkgInfo
	}
	err = updateOrInsertDatabase(pkgInfoMap)
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
	for pkgName, pkgInfo := range pkgInfoMap {
		if err := storeDependenciesInDatabase(pkgName, pkgInfo.Depends); err != nil {
			if isUniqueViolation(err) {
				continue
			}
			fmt.Printf("Error storing dependencies for package %s: %v\n", pkgName, err)
		}
	}
	fmt.Println("Database updated successfully.")

	if outputPath != "" {
		err := generateDependencyGraph(pkgInfoMap, outputPath)
		if err != nil {
			fmt.Printf("Error generating dependency graph: %v\n", err)
			return
		}
		fmt.Println("Dependency graph generated successfully.")
	}
}

func getAllDep(packages map[string]PackageInfo, pkgName string, visited map[string]bool, deps []string) []string {
	if visited[pkgName] {
		return deps
	}

	visited[pkgName] = true
	deps = append(deps, pkgName)

	if pkg, ok := packages[pkgName]; ok {
		for _, depName := range pkg.Depends {
			deps = getAllDep(packages, depName, visited, deps)
		}
	}
	return deps
}

func decompressGzip(data []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	uncompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(uncompressedData), nil
}

func getDependencies() (map[string]PackageInfo, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var packages []PackageInfo
	err = json.Unmarshal([]byte(body), &packages)
	if err != nil {
		return nil, err
	}

	pkgMap := make(map[string]PackageInfo)
	for _, pkg := range packages {
		pkgMap[pkg.Name] = pkg
	}
	return pkgMap, nil
}

func calculatePageRank(pkgInfoMap map[string]PackageInfo, iterations int, dampingFactor float64) map[string]float64 {
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

		for pkgName, pkgInfo := range pkgInfoMap {
			var depNum int
			for _, depName := range pkgInfo.Depends {
				if _, exists := pkgInfoMap[depName]; exists {
					depNum++
				}
			}
			for _, depName := range pkgInfo.Depends {
				if _, exists := pkgInfoMap[depName]; exists {
					newPageRank[depName] += dampingFactor * (pageRank[pkgName] / float64(depNum))
				}
			}
		}
		pageRank = newPageRank
	}
	return pageRank
}

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}

func updateOrInsertDatabase(pkgInfoMap map[string]PackageInfo) error {
	db, err := storage.GetDefaultAppDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	for pkgName, pkgInfo := range pkgInfoMap {
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

func storeDependenciesInDatabase(pkgName string, dependencies []string) error {
	db, err := storage.GetDefaultAppDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	for _, dep := range dependencies {
		_, err := db.Exec("INSERT INTO aur_relationships (frompackage, topackage) VALUES ($1, $2)", pkgName, dep)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateDependencyGraph(pkgInfoMap map[string]PackageInfo, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("digraph {\n")

	packageIndices := make(map[string]int)
	index := 0

	for pkgName, pkgInfo := range pkgInfoMap {
		packageIndices[pkgName] = index
		label := fmt.Sprintf("%s@%s", pkgName, pkgInfo.Description)
		writer.WriteString(fmt.Sprintf("  %d [label=\"%s\"];\n", index, label))
		index++
	}

	for pkgName, pkgInfo := range pkgInfoMap {
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
