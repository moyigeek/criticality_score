package gentoo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

// PackageInfo struct to store package information
type PackageInfo struct {
	Name         string
	Version      string
	Description  string
	Homepage     string
	Depends      []string
	DependsCount int
	URL          string
	GitRepo      string
}

// extractNameAndVersion extracts the package name and version from the file name
func extractNameAndVersion(fileName string) (string, string) {
	lastDashIndex := strings.LastIndex(fileName, "-") // 找到最后一个 `-` 的位置
	if lastDashIndex != -1 {
		// 检查版本号的第一个字符是否是 'r'
		versionPart := fileName[lastDashIndex+1:]
		if strings.HasPrefix(versionPart, "r") {
			// 如果是 'r'，则向前查找一个 `-` 作为版本号的分隔符
			secondLastDashIndex := strings.LastIndex(fileName[:lastDashIndex], "-")
			if secondLastDashIndex != -1 {
				name := fileName[:secondLastDashIndex] // 包名是最后一个 `-` 前的部分
				version := strings.TrimSuffix(fileName[secondLastDashIndex+1:], ".ebuild") // 版本号去掉后缀
				return name, version
			}
		} else {
			name := fileName[:lastDashIndex] // 包名是最后一个 `-` 前的部分
			versionParts := strings.Split(fileName[lastDashIndex+1:], ".") // 版本号是最后一个 `-` 后的部分
			if len(versionParts) > 0 {
				version := strings.TrimSuffix(versionParts[0], ".ebuild") // 版本号去掉后缀
				return name, version
			}
		}
	}
	return "", ""
}

// ParseEbuild parses a Gentoo ebuild file and returns a PackageInfo struct
func ParseEbuild(filePath string) (PackageInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return PackageInfo{}, fmt.Errorf("Error opening ebuild file: %v", err)
	}
	defer file.Close()

	var pkgInfo PackageInfo
	scanner := bufio.NewScanner(file)

	// 从文件名中提取包名和版本
	fileName := filepath.Base(filePath)
	pkgInfo.Name, pkgInfo.Version = extractNameAndVersion(fileName)

	// Regular expressions for parsing
	reDescription := regexp.MustCompile(`^DESCRIPTION="(.+)"$`)
	reHomepage := regexp.MustCompile(`^HOMEPAGE="(.+)"$`)
	reRDEPEND := regexp.MustCompile(`^RDEPEND="(.+)"$`)
	reDEPEND := regexp.MustCompile(`^DEPEND="(.+)"$`)
	reBDEPEND := regexp.MustCompile(`^BDEPEND="(.+)"$`)
	reURL := regexp.MustCompile(`^SRC_URI="(.+)"$`)

	var currentDependencyType string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if reDescription.MatchString(line) {
			matches := reDescription.FindStringSubmatch(line)
			pkgInfo.Description = matches[1]
		} else if reHomepage.MatchString(line) {
			matches := reHomepage.FindStringSubmatch(line)
			pkgInfo.Homepage = matches[1]
		} else if reRDEPEND.MatchString(line) {
			matches := reRDEPEND.FindStringSubmatch(line)
			currentDependencyType = "RDEPEND"
			pkgInfo.Depends = append(pkgInfo.Depends, parseDependencies(matches[1])...)
		} else if reDEPEND.MatchString(line) {
			matches := reDEPEND.FindStringSubmatch(line)
			currentDependencyType = "DEPEND"
			pkgInfo.Depends = append(pkgInfo.Depends, parseDependencies(matches[1])...)
		} else if reBDEPEND.MatchString(line) {
			matches := reBDEPEND.FindStringSubmatch(line)
			currentDependencyType = "BDEPEND"
			pkgInfo.Depends = append(pkgInfo.Depends, parseDependencies(matches[1])...)
		} else if reURL.MatchString(line) {
			matches := reURL.FindStringSubmatch(line)
			pkgInfo.URL = matches[1]
			if strings.HasPrefix(pkgInfo.URL, "https://github.com/") {
				parts := strings.Split(pkgInfo.URL, "/")
				if len(parts) >= 5 {
					orgName := parts[3]
					repoName := parts[4]
					// 替换 ${PN} 为包名
					if strings.Contains(orgName, "${PN}") {
						orgName = strings.Replace(orgName, "${PN}", pkgInfo.Name, -1)
					}
					if strings.Contains(repoName, "${PN}") {
						repoName = strings.Replace(repoName, "${PN}", pkgInfo.Name, -1)
					}
					pkgInfo.GitRepo = fmt.Sprintf("https://github.com/%s/%s.git", orgName, repoName)
				}
			}
		} else if currentDependencyType != "" {
			// 如果当前依赖类型不为空，继续解析依赖项
			pkgInfo.Depends = append(pkgInfo.Depends, parseDependencies(line)...)
		}
	}

	if err := scanner.Err(); err != nil {
		return PackageInfo{}, fmt.Errorf("Error reading ebuild file: %v", err)
	}

	pkgInfo.DependsCount = len(pkgInfo.Depends)
	return pkgInfo, nil
}

// parseDependencies parses a string of dependencies and returns a slice of package names
func parseDependencies(dependencyStr string) []string {
	var dependencies []string
	// 以空格分割依赖项
	dependencyLines := strings.Split(dependencyStr, "\n")
	for _, line := range dependencyLines {
		// 去除前后的空格
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 去除后缀的 "=" 和其他字符
		if strings.Contains(line, "/") {
			// 提取包名
			parts := strings.Split(line, "/")
			if len(parts) >= 2 {
				packageName := parts[1]
				// 去除后面的 "="
				if strings.Contains(packageName, ":=") {
					packageName,_ = extractNameAndVersion(strings.Split(packageName, ":=")[0])
				}
				dependencies = append(dependencies, packageName)
			}
		}
	}
	return dependencies
}

// FetchAndParseEbuildFiles fetches and parses all ebuild files in a given directory
func FetchAndParseEbuildFiles(directory string) (map[string]PackageInfo, error) {
	pkgInfoMap := make(map[string]PackageInfo)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the current item is a directory
		if info.IsDir() {
			return nil // Continue walking through subdirectories
		}
		// Check if the file has a .ebuild extension
		if strings.HasSuffix(info.Name(), ".ebuild") {
			pkgInfo, err := ParseEbuild(path)
			if err != nil {
				return fmt.Errorf("failed to parse ebuild file %s: %v", path, err)
			}
			if pkgInfo.Name != "" {
				pkgInfoMap[pkgInfo.Name] = pkgInfo
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk through ebuild directory: %v", err)
	}

	return pkgInfoMap, nil
}

// UpdateOrInsertDatabase updates or inserts package information into the database
func UpdateOrInsertDatabase(pkgInfoMap map[string]PackageInfo) error {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	for pkgName, pkgInfo := range pkgInfoMap {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM gentoo_packages WHERE package = $1)", pkgName).Scan(&exists)
		if err != nil {
			return err
		}
		if !exists {
			_, err := db.Exec("INSERT INTO gentoo_packages (package, version, depends_count, description, homepage, git_link) VALUES ($1, $2, $3, $4, $5, $6)",
				pkgName, pkgInfo.Version, pkgInfo.DependsCount, pkgInfo.Description, pkgInfo.Homepage, pkgInfo.GitRepo)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec("UPDATE gentoo_packages SET version = $1, depends_count = $2, description = $3, homepage = $4, git_link = $5 WHERE package = $6",
				pkgInfo.Version, pkgInfo.DependsCount, pkgInfo.Description, pkgInfo.Homepage, pkgInfo.GitRepo, pkgName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Gentoo function to fetch, parse ebuild files and generate dependency graph
func Gentoo(outputPath string) {
	// Step 1: Clone the Gentoo repository if it doesn't exist
	baseDirectory := "gentoo" // Define the base directory for cloning
	err := cloneGentooRepo(baseDirectory)
	if err != nil {
		fmt.Printf("Error cloning Gentoo repository: %v\n", err)
		return
	}

	// Step 2: Fetch and parse ebuild files from the cloned directory
	pkgInfoMap, err := FetchAndParseEbuildFiles(baseDirectory)
	if err != nil {
		fmt.Printf("Error fetching package info: %v\n", err)
		return
	}

	// Step 3: Calculate dependencies using the new logic
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

	for pkgName, pkgInfo := range pkgInfoMap {
		depCount := countMap[pkgName]
		pkgInfo.DependsCount = depCount
		pkgInfoMap[pkgName] = pkgInfo
	}

	// Step 4: Update or insert package information into the database
	err = UpdateOrInsertDatabase(pkgInfoMap)
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
	fmt.Println("Database updated successfully.")

	// Step 5: Generate dependency graph if output path is provided
	if outputPath != "" {
		err := generateDependencyGraph(pkgInfoMap, outputPath)
		if err != nil {
			fmt.Printf("Error generating dependency graph: %v\n", err)
			return
		}
		fmt.Println("Dependency graph generated successfully.")
	}
}

// cloneGentooRepo clones the Gentoo repository from GitHub if it doesn't already exist
func cloneGentooRepo(baseDirectory string) error {
	repoURL := "https://github.com/gentoo/gentoo.git"
	dir := filepath.Join(baseDirectory)

	if _, err := os.Stat(dir); err == nil {
		return nil // Directory already exists
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check directory: %v", err)
	}

	cmd := exec.Command("git", "clone", repoURL, dir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}

	return nil
}

// getAllDep calculates all dependencies for a given package
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

// generateDependencyGraph generates a dependency graph and writes it to a file
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
