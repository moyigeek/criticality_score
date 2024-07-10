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
)

var cacheDir = "/tmp/cloc-debian-cache"

type DepInfo struct {
	Name    string
	Arch    string
	Version string
}

func getMirrorFile(path string) []byte {
	resp, _ := http.Get("http://mirrors.hust.college/debian/" + path)
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
				depInfo := toDep(strings.TrimSpace(dep))
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

func toDep(dep string) DepInfo {
	re := regexp.MustCompile(`^(.+?)(:.+?)?(\s\((.+)\))?(\s\|.+)?$`)
	matches := re.FindStringSubmatch(dep)
	if matches != nil {
		return DepInfo{Name: matches[1], Arch: matches[2], Version: matches[4]}
	}
	return DepInfo{Name: dep, Arch: "", Version: ""}
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

func Debian() {
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

	fmt.Println("Writing result...")
	file, _ := os.Create("result_debian.csv")
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("name,refcount\n")
	for key, count := range countMap {
		writer.WriteString(fmt.Sprintf("%s,%d\n", key, count))
	}
	writer.Flush()
}
