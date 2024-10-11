package debian

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	_ "github.com/lib/pq" // Assuming PostgreSQL, adjust as needed
)

var cacheDir = "/tmp/cloc-debian-cache"

// Updated DepInfo struct to include Description and Homepage
type DepInfo struct {
	Name        string
	Arch        string
	Version     string
	Description string
	Homepage    string
}

func updateOrInsertDatabase(pkgInfoMap map[string]PackageInfo) error {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	for pkgName, pkgInfo := range pkgInfoMap {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM debian_packages WHERE package = $1)", pkgName).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			// TODO: if the package does not exist, notify user to update git_link
			_, err := db.Exec("INSERT INTO debian_packages (package, depends_count, description, homepage) VALUES ($1, $2, $3, $4)",
				pkgName, pkgInfo.DependsCount, pkgInfo.Description, pkgInfo.Homepage)
			if err != nil {
				return err
			}
		} else {

			_, err := db.Exec("UPDATE debian_packages SET depends_count = $1, description = $2, homepage = $3 WHERE package = $4",
				pkgInfo.DependsCount, pkgInfo.Description, pkgInfo.Homepage, pkgName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getMirrorFile(path string) []byte {
	resp, _ := http.Get("https://mirrors.hust.edu.cn/debian/" + path)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func getDecompressedFile(path string) string {
	file := getMirrorFile(path)
	reader, _ := gzip.NewReader(strings.NewReader(string(file)))
	defer reader.Close()

	decompressed, _ := ioutil.ReadAll(reader)
	return string(decompressed)
}

func getPackageList() string {
	return getDecompressedFile("dists/stable/main/binary-amd64/Packages.gz")
}

func parseList() map[string]map[string]interface{} {
	content := getPackageList()
	lists := strings.Split(content, "\n\n")
	packages := make(map[string]map[string]interface{})

	for _, packageStr := range lists {
		if strings.TrimSpace(packageStr) == "" {
			continue
		}
		pkg := make(map[string]interface{})
		pkg["__raw"] = packageStr
		var currentKey string

		lines := strings.Split(packageStr, "\n")
		for _, line := range lines {
			if matched, _ := regexp.MatchString(".+:.+", line); matched {
				parts := strings.SplitN(line, ":", 2)
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				pkg[key] = value
				currentKey = key
			} else if matched, _ := regexp.MatchString(".+:\\s*", line); matched {
				currentKey = strings.Split(line, ":")[0]
				pkg[currentKey] = ""
			} else if matched, _ := regexp.MatchString(" .+", line); matched {
				// 类型断言，确保 currentKey 对应的值是字符串
				if currentValue, ok := pkg[currentKey].(string); ok {
					pkg[currentKey] = currentValue + " " + strings.TrimSpace(line)
				}
			}
		}

		// 后处理
		if depends, ok := pkg["Depends"].(string); ok {
			depList := strings.Split(depends, ",")
			var depStrings []interface{}
			for _, dep := range depList {
				depInfo := toDep(strings.TrimSpace(dep), packageStr)
				depStrings = append(depStrings, depInfo)
			}
			pkg["Depends"] = depStrings
		}

		if packageName, ok := pkg["Package"].(string); ok {
			packages[packageName] = pkg
		}
	}

	return packages
}

// Updated toDep function to extract additional fields
func toDep(dep string, rawContent string) DepInfo {
	re := regexp.MustCompile(`^(.+?)(:.+?)?(\s\((.+)\))?(\s\|.+)?$`)
	matches := re.FindStringSubmatch(dep)

	// Initialize DepInfo with default values
	depInfo := DepInfo{Name: dep, Arch: "", Version: "", Description: "", Homepage: ""}

	if matches != nil {
		depInfo.Name = matches[1]
		if matches[2] != "" {
			depInfo.Arch = strings.TrimSpace(matches[2])
		}
		if matches[4] != "" {
			depInfo.Version = strings.TrimSpace(matches[4])
		}
	}

	// Extract Description and Homepage from rawContent
	descriptionRegex := regexp.MustCompile(`(?m)^Description:\s*(.*)$`)
	homepageRegex := regexp.MustCompile(`(?m)^Homepage:\s*(.*)$`)

	// Extract Description
	if descMatches := descriptionRegex.FindStringSubmatch(rawContent); len(descMatches) > 1 {
		depInfo.Description = descMatches[1]
	}

	// Extract Homepage
	if homeMatches := homepageRegex.FindStringSubmatch(rawContent); len(homeMatches) > 1 {
		depInfo.Homepage = homeMatches[1]
	}

	return depInfo
}

func generateDependencyGraph(packages map[string]map[string]interface{}, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("digraph {\n")

	// Create a map to store package indices
	packageIndices := make(map[string]int)
	index := 0

	// Assign an index to each package and write the node definitions
	for pkgName, pkgInfo := range packages {
		packageIndices[pkgName] = index
		label := fmt.Sprintf("%s@%s", pkgName, pkgInfo["Version"].(string))
		writer.WriteString(fmt.Sprintf("  %d [label=\"%s\"];\n", index, label))
		index++
	}

	// Write the edges (dependencies)
	for pkgName, pkgInfo := range packages {
		pkgIndex := packageIndices[pkgName]
		if depends, ok := pkgInfo["Depends"].([]interface{}); ok {
			for _, depInterface := range depends {
				if depInfo, ok := depInterface.(DepInfo); ok {
					if depIndex, ok := packageIndices[depInfo.Name]; ok {
						writer.WriteString(fmt.Sprintf("  %d -> %d [label=\"%s\"];\n", pkgIndex, depIndex, depInfo.Version))
					}
				}
			}
		}
	}

	writer.WriteString("}\n")
	writer.Flush()
	return nil
}

func getAllDep(packages map[string]map[string]interface{}, pkgName string, deps []string) []string {
	deps = append(deps, pkgName)
	if pkg, ok := packages[pkgName]; ok {
		if depends, ok := pkg["Depends"].([]interface{}); ok {
			for _, depInterface := range depends {
				if depMap, ok := depInterface.(DepInfo); ok {
					pkgname := depMap.Name
					if !contains(deps, pkgname) {
						deps = getAllDep(packages, pkgname, deps)
					}
				}
			}
		}
	}
	return deps
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func Debian(outputPath string) {
	fmt.Println("Getting package list...")
	packages := parseList()
	fmt.Printf("Done, total: %d packages.\n", len(packages))
	fmt.Println("Building dependencies graph...")

	keys := make([]string, 0, len(packages))
	for k := range packages {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// fmt.Println(keys)
	depMap := make(map[string][]string)
	for _, pkgName := range keys {
		deps := getAllDep(packages, pkgName, []string{})
		depMap[pkgName] = deps
	}
	fmt.Println("Calculating dependencies count...")
	countMap := make(map[string]int)
	for _, deps := range depMap {
		for _, dep := range deps {
			countMap[dep]++
		}
	}

	// Create a map to hold package information
	pkgInfoMap := make(map[string]PackageInfo)

	// Populate the pkgInfoMap with counts, descriptions, and homepages
	for pkgName, pkgInfo := range packages {
		depCount := countMap[pkgName] // Get the dependency count

		// Safely extract Description and Homepage, defaulting to empty string if not present
		description, ok := pkgInfo["Description"].(string)
		if !ok {
			description = "" // Set to empty string if not found
		}

		homepage, ok := pkgInfo["Homepage"].(string)
		if !ok {
			homepage = "" // Set to empty string if not found
		}

		pkgInfoMap[pkgName] = PackageInfo{
			DependsCount: depCount,
			Description:  description,
			Homepage:     homepage,
		}
	}

	// Update database with package information
	err := updateOrInsertDatabase(pkgInfoMap)
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
	fmt.Println("Database updated successfully.")

	if outputPath != "" {
		err := generateDependencyGraph(packages, outputPath)
		if err != nil {
			fmt.Printf("Error generating dependency graph: %v\n", err)
			return
		}
		fmt.Println("Dependency graph generated successfully.")
	}
}

type PackageInfo struct {
	DependsCount int
	Description  string
	Homepage     string
}
