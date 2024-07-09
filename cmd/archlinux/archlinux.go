package archlinux

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type DepInfo struct {
	Name    string
	Arch    string
	Version string
}

func toDep(dep string) DepInfo {
	re := regexp.MustCompile(`^([^=><!]+?)(?:([=><!]+)([^:]+))?(?::(.+?))?(?:\s*\((.+)\))?$`)
	matches := re.FindStringSubmatch(dep)
	if matches != nil {
		return DepInfo{Name: matches[1], Arch: matches[4], Version: matches[2] + matches[3]}
	}
	return DepInfo{Name: dep, Arch: "", Version: ""}
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
			pkgInfo = toDep(line)
		}
		if inDependSection && line != "" {
			dependencies = append(dependencies, toDep(line))
		}
	}

	if err := scanner.Err(); err != nil {
		return DepInfo{}, nil, err
	}
	return pkgInfo, dependencies, nil
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

func Archlinux() {
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

	fmt.Println("Writing result...")
	file, _ := os.Create("result_archlinux.csv")
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString("name,refcount\n")
	for key, count := range countMap {
		writer.WriteString(fmt.Sprintf("%s,%d\n", key, count))
	}
	writer.Flush()
}
