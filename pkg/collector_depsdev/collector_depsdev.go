package collector_depsdev

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

type DependentInfo struct {
	DependentCount         int `json:"dependentCount"`
	DirectDependentCount   int `json:"directDependentCount"`
	IndirectDependentCount int `json:"indirectDependentCount"`
}

func Run(inputName, outputName string) {
	file, err := os.Open(inputName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "https://github.com/") {
			parts := strings.Split(line, "/")
			if len(parts) >= 5 {
				owner := parts[3]
				repo := parts[4]
				projectType := getProjectType(owner, repo)
				if projectType != "" {
					latestVersion := getLatestVersion(owner, repo, projectType)
					if latestVersion != "" {
						queryDepsDev(projectType, repo, latestVersion, outputName)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func getProjectType(owner, repo string) string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 获取仓库内容
	_, dirContent, _, err := client.Repositories.GetContents(ctx, owner, repo, "", nil)
	if err != nil {
		fmt.Println("Error fetching repository contents:", err)
		return ""
	}

	// 判断项目类型
	for _, file := range dirContent {
		switch *file.Name {
		case "package.json":
			return "npm"
		case "setup.py":
			return "pypi"
		case "Cargo.toml":
			return "cargo"
		case "pom.xml":
			return "maven"
		case "build.gradle":
			return "gradle"
		case "go.mod":
			return "go"
		}
	}
	return ""
}

func getLatestVersion(owner, repo, projectType string) string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_AUTH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 获取最新的发布版本
	release, _, err := client.Repositories.GetLatestRelease(ctx, owner, repo)
	if err != nil {
		fmt.Println("Error fetching latest release:", err)
		return ""
	}
	return release.GetTagName()
}

func queryDepsDev(projectType, projectName, version, outputFile string) {
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}
	url := fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s/versions/%s:dependents", projectType, projectName, version)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error querying deps.dev:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code")
		return
	}

	var info DependentInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 将输出重定向到文件
	fmt.Fprintf(file, "Project: %s, Version: %s\n", projectName, version)
	fmt.Fprintf(file, "Dependent Count: %d, Direct Dependent Count: %d, Indirect Dependent Count: %d\n",
		info.DependentCount, info.DirectDependentCount, info.IndirectDependentCount)
}
