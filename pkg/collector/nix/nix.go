package nix

import (
	"bytes"
	"encoding/gob"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

// DepInfo struct to store package information
type DepInfo struct {
	Name        string
	Version     string
	Homepage    string
	Description string
	GitLink     string
	DepCount    int // 新增的依赖计数字段
}

// isValidNixIdentifier checks if a string is a valid Nix identifier
func isValidNixIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	first := s[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_') {
		return false
	}
	for _, c := range s[1:] {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}

// attributePathToNixExpression converts an attribute path to a valid Nix expression
func attributePathToNixExpression(attributePath string) string {
	components := strings.Split(attributePath, ".")
	expr := "pkgs"
	for _, comp := range components {
		if isValidNixIdentifier(comp) {
			expr += "." + comp
		} else {
			expr += `."` + comp + `"`
		}
	}
	return expr
}

// getAllNixPackages retrieves all Nix packages as DepInfo and their dependencies
func GetAllNixPackages() (map[DepInfo][]DepInfo, error) {
	cmd := exec.Command("nix-env", "-qaP")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("Error running nix-env command: %v", err)
	}

	packages := make(map[DepInfo][]DepInfo)
	lines := strings.Split(string(out), "\n")

	re := regexp.MustCompile(`^nixpkgs\.(.+?)\s+([^\s]+)$`)

	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.Contains(line, "evaluation warning") {
			continue
		}

		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			attributePath := matches[1]
			packageFullName := matches[2]

			parts := strings.Split(packageFullName, "-")

			versionIndex := -1
			for i := 1; i < len(parts); i++ {
				if len(parts[i]) > 0 && unicode.IsDigit(rune(parts[i][0])) {
					versionIndex = i
					break
				}
			}

			var packageName, packageVersion string
			if versionIndex != -1 {
				packageName = strings.Join(parts[:versionIndex], "-")
				packageVersion = strings.Join(parts[versionIndex:], "-")
			} else {
				packageName = packageFullName
				packageVersion = ""
			}

			packageInfo, err := GetNixPackageInfo(attributePath)
			if err != nil {
				fmt.Printf("Error getting info for %s: %v\n", attributePath, err)
				continue
			}

			pkgDepInfo := DepInfo{
				Name:        packageName,
				Version:     packageVersion,
				Homepage:    packageInfo.Homepage,
				Description: packageInfo.Description,
				GitLink:     packageInfo.GitLink,
			}

			dependencies, err := GetNixPackageDependencies(attributePath)
			if err != nil {
				fmt.Printf("Error getting dependencies for %s: %v\n", attributePath, err)
				continue
			}

			packages[pkgDepInfo] = dependencies
			fmt.Println(pkgDepInfo)
			fmt.Println(dependencies)
		}
	}

	return packages, nil
}

// getNixPackageInfo retrieves the package information using attribute path
func GetNixPackageInfo(attributePath string) (DepInfo, error) {
	nixPkgExpression := attributePathToNixExpression(attributePath)

	expr := fmt.Sprintf(`
let
  pkgs = import <nixpkgs> {};
  pkg = %s;
  pname = if pkg ? pname then pkg.pname else if pkg ? name then pkg.name else "";
  version = if pkg ? version then pkg.version else "unknown";
  meta = if pkg ? meta then pkg.meta else {};
  homepage = if meta ? homepage then meta.homepage else "";
  description = if meta ? description then meta.description else "";
  srcUrl = if pkg ? src then
    if pkg.src ? url then pkg.src.url else if pkg.src ? urls then builtins.elemAt pkg.src.urls 0 else ""
  else "";
  passthruUrl = if pkg ? passthru && pkg.passthru ? updateScript && pkg.passthru.updateScript ? url then
    pkg.passthru.updateScript.url
  else "";
  gitLink = if srcUrl != "" then srcUrl else passthruUrl;
in
{
  name = pname;
  version = version;
  homepage = homepage;
  description = description;
  gitLink = gitLink;
}
`, nixPkgExpression)

	cmd := exec.Command("nix", "eval", "--impure", "--expr", expr, "--extra-experimental-features", "nix-command", "--json")
	var out bytes.Buffer
	var outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	err := cmd.Run()
	if err != nil {
		return DepInfo{}, fmt.Errorf("Error running nix eval for package '%s': %v\nNix error output:\n%s", attributePath, err, outErr.String())
	}

	var result map[string]string
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		return DepInfo{}, fmt.Errorf("Error parsing JSON for package '%s': %v", attributePath, err)
	}

	depInfo := DepInfo{
		Name:        result["name"],
		Version:     result["version"],
		Homepage:    result["homepage"],
		Description: result["description"],
		GitLink:     result["gitLink"],
	}

	depInfo.GitLink = processGitLink(depInfo.GitLink)

	return depInfo, nil
}

func GetNixPackageDependencies(attributePath string) ([]DepInfo, error) {
    nixPkgExpression := attributePathToNixExpression(attributePath)

	exprTemplate := `
	let
	  pkgs = import <nixpkgs> {};
	  pkg = %s;
	in {
		buildInputs = map (x: if x ? pname then x.pname else if x ? name then x.name else "") (pkg.buildInputs or []);
	}
	`	
    evalExpr := fmt.Sprintf(exprTemplate, nixPkgExpression)
    results, err := nixEval(evalExpr)
    if err != nil {
        return nil, fmt.Errorf("Error getting dependencies for %s: %v", attributePath, err)
    }

    buildInputNames := results["buildInputs"]
    finalInputs := []DepInfo{}
    for _, name := range buildInputNames {
        finalInputs = append(finalInputs, DepInfo{Name: name.Name})
    }

    return finalInputs, nil
}

// nixEval executes a Nix expression and parses the JSON output into []DepInfo
func nixEval(expr string) (map[string][]DepInfo, error) {
    cmd := exec.Command("nix", "eval", "--impure", "--expr", expr, "--extra-experimental-features", "nix-command", "--json")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return nil, fmt.Errorf("Error running nix eval: %v", err)
    }

    // 修改为 map[string][]string 以匹配 JSON 结构
    var result map[string][]string
    // fmt.Println(string(out.Bytes())) // 打印输出以便调试
    if err := json.Unmarshal(out.Bytes(), &result); err != nil {
        return nil, fmt.Errorf("Error parsing JSON: %v", err)
    }

    depsMap := make(map[string][]DepInfo)
    
    // 构建依赖映射
    for key, depList := range result {
        for _, depName := range depList {
            depInfo := DepInfo{
                Name: depName, // 只存储名称
            }
            depsMap[key] = append(depsMap[key], depInfo)
        }
    }

    return depsMap, nil
}



// processGitLink processes the gitLink to ensure it points to a git repository
func processGitLink(gitLink string) string {
	if gitLink == "" {
		return ""
	}

	parsedURL, err := url.Parse(gitLink)
	if err != nil {
		return ""
	}

	codeHostingDomains := []string{
		"github.com",
		"gitlab.com",
		"bitbucket.org",
	}

	for _, domain := range codeHostingDomains {
		if strings.Contains(parsedURL.Host, domain) {
			return gitLink
		}
	}

	return ""
}

func mergeDependencies(packages map[DepInfo][]DepInfo) map[DepInfo][]DepInfo {
	mergedPackages := make(map[DepInfo][]DepInfo)

	for pkg, deps := range packages {
		merged := false

		// 遍历已合并的包，查找相同包名
		for existingPkg := range mergedPackages {
			if existingPkg.Name == pkg.Name {
				// 合并版本信息
				versionSet := make(map[string]struct{})
				versionSet[existingPkg.Version] = struct{}{}
				versionSet[pkg.Version] = struct{}{}
				pkg.Version = strings.Join(getKeys(versionSet), ",") // 更新版本信息

				// 合并依赖
				mergedDeps := make(map[string]DepInfo)
				for _, dep := range mergedPackages[existingPkg] {
					mergedDeps[dep.Name] = dep
				}
				for _, dep := range deps {
					mergedDeps[dep.Name] = dep // 保留最新的依赖信息
				}

				// 更新合并后的依赖列表
				mergedDepsList := make([]DepInfo, 0, len(mergedDeps))
				for _, dep := range mergedDeps {
					mergedDepsList = append(mergedDepsList, dep)
				}
				mergedPackages[existingPkg] = mergedDepsList
				merged = true
				break
			}
		}

		// 如果没有找到相同的包名，则直接添加
		if !merged {
			mergedPackages[pkg] = deps
		}
	}

	return mergedPackages
}
// writeCSV writes the collected package information to a CSV file
// 修改后的 getAllDep 函数
func getAllDep(packages map[DepInfo][]DepInfo, pkgName string, deps []string) []string {
	deps = append(deps, pkgName) // 添加当前包名到依赖列表
	// fmt.Println(pkgName)
	for pkg, depsList := range packages {
		if pkg.Name == pkgName { // 比较包名
			for _, dep := range depsList {
				pkgname := dep.Name
				if !contains(deps, pkgname) { // 检查依赖是否已在列表中
					deps = getAllDep(packages, pkgname, deps) // 递归调用
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

// 修改后的 writeCSV 函数
func writeCSV(packages map[DepInfo][]DepInfo, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入 CSV 头部
	if err := writer.Write([]string{"Package Name", "Version", "Homepage", "Description", "Git Link", "Dependencies", "Dependency Count"}); err != nil {
		return fmt.Errorf("Error writing header to CSV: %v", err)
	}

	// 合并包信息
	mergedPackages := mergeDependencies(packages)
	// fmt.Println(mergedPackages)
	// 写入合并后的包数据
	for pkg, deps := range mergedPackages {
		// fmt.Println(pkg)
		allDeps := getAllDep(mergedPackages, pkg.Name, []string{}) // 使用新的 getAllDep 逻辑

		// 计算依赖数
		pkg.DepCount = len(allDeps)

		dependencies := make([]string, len(deps))
		for i, dep := range deps {
			dependencies[i] = dep.Name
		}
		sort.Strings(dependencies)

		// 使用双引号包裹每个字段以处理逗号
		if err := writer.Write([]string{
			fmt.Sprintf("\"%s\"", pkg.Name),
			fmt.Sprintf("\"%s\"", pkg.Version),
			fmt.Sprintf("\"%s\"", pkg.Homepage),
			fmt.Sprintf("\"%s\"", pkg.Description),
			fmt.Sprintf("\"%s\"", pkg.GitLink),
			fmt.Sprintf("\"%s\"", strings.Join(dependencies, ", ")),
			fmt.Sprintf("\"%d\"", pkg.DepCount),
		}); err != nil {
			return fmt.Errorf("Error writing package data to CSV: %v", err)
		}
	}

	return nil
}


// 辅助函数：获取map的键
func getKeys(set map[string]struct{}) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	return keys
}

func SavePackage(packages map[DepInfo][]DepInfo) error {
    file, err := os.Create("packages.gob")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := gob.NewEncoder(file)
    return encoder.Encode(packages)
}

func LoadPackage() (map[DepInfo][]DepInfo, error) {
    file, err := os.Open("packages.gob")
    if err != nil {
        if os.IsNotExist(err) {
            return nil, nil // 文件不存在，返回 nil
        }
        return nil, err
    }
    defer file.Close()

    var packages map[DepInfo][]DepInfo
    decoder := gob.NewDecoder(file)
    err = decoder.Decode(&packages)
    if err != nil {
        return nil, err
    }

    return packages, nil
}

func GetNixPackageList() ([]DepInfo, error) {
    cmd := exec.Command("nix-env", "-qaP")
    out, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("Error running nix-env command: %v", err)
    }

    var packages []DepInfo
    lines := strings.Split(string(out), "\n")

    re := regexp.MustCompile(`^nixpkgs\.(.+?)\s+([^\s]+)$`)
    for _, line := range lines {
        if strings.TrimSpace(line) == "" || strings.Contains(line, "evaluation warning") {
            continue
        }

        matches := re.FindStringSubmatch(line)
        if len(matches) == 3 {
            packages = append(packages, DepInfo{Name: matches[1], Version: matches[2]})
        }
    }
    return packages, nil
}

func reverseDependencies(deps map[DepInfo][]DepInfo) map[DepInfo][]DepInfo {
    reversed := make(map[DepInfo][]DepInfo)
    for key, values := range deps {
        for _, dep := range values {
            simplifiedDep := DepInfo{Name: dep.Name} // 只保留 Name 字段
			// fmt.Println(key)
            reversed[simplifiedDep] = append(reversed[simplifiedDep], key)
        }
    }
    return reversed
}

func Nix() {
    // 检查是否存在缓存文件
    packages, err := LoadPackage()
    if err != nil {
        fmt.Printf("Error loading package list: %v\n", err)
        return
    }

    if packages == nil {
        // 如果没有缓存，则获取新的包列表
        packages, err = GetAllNixPackages()
        if err != nil {
            fmt.Printf("Error retrieving Nix packages: %v\n", err)
            return
        }

        // 存储到文件
        if err := SavePackage(packages); err != nil {
            fmt.Printf("Error saving package list: %v\n", err)
            return
        }
    }

    // 反转依赖关系
    simplifiedDep := reverseDependencies(packages)

    // 更新 simplifiedDep 的键名
    for dep, deps := range simplifiedDep {
        if fullDep, exists := findDepInfoByName(packages, dep.Name); exists {
            // 删除旧键
            delete(simplifiedDep, dep)

            // 使用完整信息更新键
            simplifiedDep[fullDep] = deps
        }
    }

	for dep := range packages {
        if _, exists := simplifiedDep[dep]; !exists {
            simplifiedDep[dep] = nil  // Add packages with no dependencies
        }
    }

    // 将包信息存储到数据库
    if err := updateOrInsertNixPackages(simplifiedDep); err != nil {
        fmt.Printf("Error updating or inserting Nix packages into database: %v\n", err)
        return
    }

    fmt.Println("Successfully updated package information in the database")
}

// 新增函数：将 Nix 包信息存储到数据库
func updateOrInsertNixPackages(packages map[DepInfo][]DepInfo) error {
    db, err := storage.GetDatabaseConnection()
    if err != nil {
        return err
    }
    defer db.Close()

    for pkg, deps := range packages {
        // 假设我们只存储包名和依赖数量
        var exists bool
        err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM nix_packages WHERE package = $1)", pkg.Name).Scan(&exists)
        if err != nil {
            return err
        }

        if !exists {
            // 插入新包信息
            _, err := db.Exec("INSERT INTO nix_packages (package, version, homepage, description, git_link, depends_count) VALUES ($1, $2, $3, $4, $5, $6)",
                pkg.Name, pkg.Version, pkg.Homepage, pkg.Description, normalizeGitLink(pkg.GitLink), len(deps))
            if err != nil {
                return err
            }
        } else {
            // 检查当前 git_link 是否为空
            var currentGitLink *string
            err := db.QueryRow("SELECT git_link FROM nix_packages WHERE package = $1", pkg.Name).Scan(&currentGitLink)
            if err != nil {
                return err
            }

            // 更新其他字段，如果 currentGitLink 为空则更新 git_link
            // if currentGitLink == nil || *currentGitLink == "" {
                _, err = db.Exec("UPDATE nix_packages SET version = $1, homepage = $2, description = $3, git_link = $4, depends_count = $5 WHERE package = $6",
                    pkg.Version, pkg.Homepage, pkg.Description, normalizeGitLink(pkg.GitLink), len(deps), pkg.Name)
                if err != nil {
                    return err
                }
        //     } else {
        //         // 只更新其他字段，不更新 git_link
        //         _, err := db.Exec("UPDATE nix_packages SET version = $1, homepage = $2, description = $3, depends_count = $4 WHERE name = $5",
        //             pkg.Version, pkg.Homepage, pkg.Description, len(deps), pkg.Name)
        //         if err != nil {
        //             return err
        //         }
        //     }
        }
    }
    return nil
}

// findDepInfoByName 根据包名在包列表中查找对应的 DepInfo
func findDepInfoByName(packages map[DepInfo][]DepInfo, name string) (DepInfo, bool) {
    for dep := range packages {
        if dep.Name == name {
            return dep, true
        }
    }
    return DepInfo{}, false
}

func normalizeGitLink(link string) string {
	// 检查并提取组织名和仓库名
	var orgName, repoName string

	if strings.HasPrefix(link, "https://github.com/") {
		parts := strings.Split(link, "/")
		if len(parts) >= 5 {
			orgName = parts[3]
			repoName = parts[4]
			return fmt.Sprintf("https://github.com/%s/%s.git", orgName, repoName)
		}
	} else if strings.HasPrefix(link, "http://github.com/") {
		parts := strings.Split(link, "/")
		if len(parts) >= 5 {
			orgName = parts[3]
			repoName = parts[4]
			return fmt.Sprintf("http://github.com/%s/%s.git", orgName, repoName)
		}
	} else if strings.HasPrefix(link, "https://gitlab.com/") {
		parts := strings.Split(link, "/")
		if len(parts) >= 5 {
			orgName = parts[3]
			repoName = parts[4]
			return fmt.Sprintf("https://gitlab.com/%s/%s.git", orgName, repoName)
		}
	} else if strings.HasPrefix(link, "http://gitlab.com/") {
		parts := strings.Split(link, "/")
		if len(parts) >= 5 {
			orgName = parts[3]
			repoName = parts[4]
			return fmt.Sprintf("http://gitlab.com/%s/%s.git", orgName, repoName)
		}
	} else if strings.HasPrefix(link, "https://gitee.com/") {
		parts := strings.Split(link, "/")
		if len(parts) >= 5 {
			orgName = parts[3]
			repoName = parts[4]
			return fmt.Sprintf("https://gitee.com/%s/%s.git", orgName, repoName)
		}
	} else if strings.HasPrefix(link, "http://gitee.com/") {
		parts := strings.Split(link, "/")
		if len(parts) >= 5 {
			orgName = parts[3]
			repoName = parts[4]
			return fmt.Sprintf("http://gitee.com/%s/%s.git", orgName, repoName)
		}
	} else if strings.HasPrefix(link, "https://bitbucket.org/") {
		parts := strings.Split(link, "/")
		if len(parts) >= 5 {
			orgName = parts[3]
			repoName = parts[4]
			return fmt.Sprintf("https://bitbucket.org/%s/%s.git", orgName, repoName)
		}
	} else if strings.HasPrefix(link, "http://bitbucket.org/") {
		parts := strings.Split(link, "/")
		if len(parts) >= 5 {
			orgName = parts[3]
			repoName = parts[4]
			return fmt.Sprintf("http://bitbucket.org/%s/%s.git", orgName, repoName)
		}
	}

	return "" // 如果不符合任何协议，返回空字符串
}
