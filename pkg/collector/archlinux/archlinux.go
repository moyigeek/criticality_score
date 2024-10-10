package archlinux

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	_ "github.com/lib/pq" // Assuming PostgreSQL, adjust as needed
)

type Config struct {
	Database    string `json:"database"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	GitHubToken string `json:"GitHubToken"`
}

func loadConfig(configPath string) (Config, error) {
	var config Config
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	return config, err
}

func updateDatabase(pkgInfoMap map[string]DepInfo, config Config) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	for pkgName, pkgInfo := range pkgInfoMap {
		_, err := db.Exec("UPDATE arch_packages SET depends_count = $1, description = $2, homepage = $3 WHERE package = $4",
			pkgInfo.DependsCount, pkgInfo.Description, pkgInfo.Homepage, pkgName)
		if err != nil {
			return err
		}
	}
	return nil
}

type DepInfo struct {
	Name         string
	Arch         string
	Version      string
	Description  string // New field for package description
	Homepage     string // New field for package homepage
	DependsCount int    // New field for dependency count
}

func toDep(dep string, rawContent string) DepInfo {
	re := regexp.MustCompile(`^([^=><!]+?)(?:([=><!]+)([^:]+))?(?::(.+?))?(?:\s*\((.+)\))?$`)
	matches := re.FindStringSubmatch(dep)

	// Initialize DepInfo with default values
	depInfo := DepInfo{Name: dep, Arch: "", Version: "", Description: "", Homepage: ""}

	if matches != nil {
		depInfo.Name = matches[1]
		depInfo.Version = matches[2] + matches[3]
		depInfo.Arch = matches[4]
	}

	// Extract Description and Homepage from rawContent
	descriptionRegex := regexp.MustCompile(`(?m)^%DESC%\s*(.+)$`)
	homepageRegex := regexp.MustCompile(`(?m)^%URL%\s*(.+)$`)

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

func extractTarGz(gzipStream io.Reader, dest string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)
	hasFiles := false

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		hasFiles = true
		target := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				file.Close()
				return err
			}
			file.Close()
		}
	}

	if !hasFiles {
		return fmt.Errorf("empty tar archive")
	}

	return nil
}

func readDescFile(descPath string) (DepInfo, []DepInfo, error) {
	file, err := os.Open(descPath)
	if err != nil {
		return DepInfo{}, nil, err
	}
	defer file.Close()

	var pkgInfo DepInfo
	var dependencies []DepInfo
	var inPackageSection, inDependSection bool
	var rawContent strings.Builder // 用于构建完整的原始内容
	var expectNextLine string      // 用于存储期待的下一行内容

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "%NAME%" {
			inPackageSection = true
			continue
		}
		if line == "%DEPENDS%" {
			inDependSection = true
			inPackageSection = false
			continue
		}
		if strings.HasPrefix(line, "%") {
			inPackageSection = false
			inDependSection = false
		}

		if inPackageSection && line != "" {
			rawContent.WriteString(line + "\n")        // 将当前行添加到原始内容中
			pkgInfo = toDep(line, rawContent.String()) // 传递完整的原始内容
		}

		if inDependSection && line != "" {
			rawContent.WriteString(line + "\n")                                   // 将当前行添加到原始内容中
			dependencies = append(dependencies, toDep(line, rawContent.String())) // 传递完整的原始内容
		}

		// 处理特定的标记
		if line == "%URL%" {
			expectNextLine = "url" // 标记期待下一行是 URL
		} else if line == "%DESC%" {
			expectNextLine = "desc" // 标记期待下一行是描述
		} else if expectNextLine == "url" {
			rawContent.WriteString("%URL%\n" + line + "\n") // 将URL行添加到原始内容中
			expectNextLine = ""                             // 重置标记
		} else if expectNextLine == "desc" {
			rawContent.WriteString("%DESC%\n" + line + "\n") // 将描述行添加到原始内容中
			expectNextLine = ""                              // 重置标记
		}
	}

	if err := scanner.Err(); err != nil {
		return DepInfo{}, nil, err
	}
	return pkgInfo, dependencies, nil
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
		label := fmt.Sprintf("%s@%s", pkgName, pkgInfo["Info"].(DepInfo).Version)
		writer.WriteString(fmt.Sprintf("  %d [label=\"%s\"];\n", index, label))
		index++
	}

	// Write the edges (dependencies)
	for pkgName, pkgInfo := range packages {
		pkgIndex := packageIndices[pkgName]
		if depends, ok := pkgInfo["Depends"].([]DepInfo); ok {
			for _, dep := range depends {
				if depIndex, ok := packageIndices[dep.Name]; ok {
					writer.WriteString(fmt.Sprintf("  %d -> %d [label=\"%s\"];\n", pkgIndex, depIndex, dep.Version))
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
		if depends, ok := pkg["Depends"].([]DepInfo); ok {
			for _, dep := range depends {
				pkgname := dep.Name
				if !contains(deps, pkgname) {
					deps = getAllDep(packages, pkgname, deps)
				}
			}
		}
	}
	return deps
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func Archlinux(outputPath string) {
	downloadDir := "./download"

	// Check if download directory exists
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		fmt.Println("Download directory not found, starting download...")
		DownloadFiles()
	}

	fmt.Println("Getting package list...")
	extractDir := "./extracted"
	packages := make(map[string]map[string]interface{})
	packageNamePattern := regexp.MustCompile(`^([a-zA-Z0-9\-_]+)-([0-9\._]+)`)

	// Create extract directory if it doesn't exist
	if _, err := os.Stat(extractDir); os.IsNotExist(err) {
		err := os.Mkdir(extractDir, 0755)
		if err != nil {
			fmt.Printf("Error creating extract directory: %v\n", err)
			return
		}
	}

	// Walk through download directory
	err := filepath.Walk(downloadDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), ".tar.gz") {
			// Extract tar.gz file
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			err = extractTarGz(file, extractDir)
			if err != nil {
				if err.Error() == "empty tar archive" {
					fmt.Printf("Skipping empty tar archive: %s\n", path)
					return nil // Skip empty tar archive and continue
				}
				fmt.Printf("Error extracting %s: %v\n", path, err)
				return nil // Skip this file and continue
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking through download directory: %v\n", err)
		return
	}

	// Walk through extracted directory
	err = filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), "desc") {
			// Parse desc file
			packageName := packageNamePattern.FindStringSubmatch(filepath.Base(filepath.Dir(path)))
			if packageName != nil {
				pkgInfo, dependencies, err := readDescFile(path)
				if err != nil {
					return err
				}
				if _, ok := packages[pkgInfo.Name]; !ok {
					packages[pkgInfo.Name] = make(map[string]interface{})
				}
				packages[pkgInfo.Name]["Depends"] = dependencies
				packages[pkgInfo.Name]["Info"] = pkgInfo
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking through extracted directory: %v\n", err)
		return
	}
	fmt.Printf("Done, total: %d packages.\n", len(packages))

	if outputPath != "" {
		err := generateDependencyGraph(packages, outputPath)
		if err != nil {
			fmt.Printf("Error generating dependency graph: %v\n", err)
			return
		}
		fmt.Println("Dependency graph generated successfully.")
	}
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}
	fmt.Println("Building dependencies graph...")
	keys := make([]string, 0, len(packages))
	for k := range packages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

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
	pkgInfoMap := make(map[string]DepInfo)

	// Populate pkgInfoMap with counts, descriptions, and homepages
	for pkgName, pkgInfo := range packages {
		depCount := countMap[pkgName] // Get the dependency count

		// Safely extract Description and Homepage, defaulting to empty string if not present
		var description, homepage string

		if info, ok := pkgInfo["Info"].(DepInfo); ok {
			description = info.Description
			homepage = info.Homepage
		} else {
			description = "" // Set to empty string if Info is not of type DepInfo
			homepage = ""    // Set to empty string if Info is not of type DepInfo
		}

		pkgInfoMap[pkgName] = DepInfo{
			Name:         pkgName,
			DependsCount: depCount,
			Description:  description,
			Homepage:     homepage,
		}
	}

	// Update database with package information
	err = updateDatabase(pkgInfoMap, config) // Pass the pkgInfoMap
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
	fmt.Println("Database updated successfully.")
}
