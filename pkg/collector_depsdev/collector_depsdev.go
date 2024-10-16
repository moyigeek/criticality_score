package collector_depsdev

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/google/go-github/v63/github"
	_ "github.com/lib/pq" // Assuming PostgreSQL, adjust as needed
	"golang.org/x/oauth2"
)

type DependentInfo struct {
	DependentCount         int `json:"dependentCount"`
	DirectDependentCount   int `json:"directDependentCount"`
	IndirectDependentCount int `json:"indirectDependentCount"`
}

func updateDatabase(link, projectName string, dependentCount int, projectType string) error {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE git_metrics SET depsdev_count = $1 and pkg_manager = $2 WHERE git_link = $3", dependentCount, projectType, link) // Adjust the WHERE clause as needed
	return err
}

func Run(configPath string) {
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		fmt.Errorf("Error initializing database: %v\n", err)
		return
	}
	defer db.Close()

	// Query git_links from git_metrics table
	rows, err := db.Query("SELECT git_link FROM git_metrics")
	if err != nil {
		fmt.Println("Error querying git_metrics:", err)
		return
	}
	defer rows.Close()

	var gitLinks []string
	for rows.Next() {
		var gitLink string
		if err := rows.Scan(&gitLink); err != nil {
			fmt.Println("Error scanning git_link:", err)
			return
		}
		gitLinks = append(gitLinks, gitLink)
	}

	// Process each git link
	for _, link := range gitLinks {
		if strings.HasPrefix(link, "https://github.com/") || strings.HasPrefix(link, "https://gitlab.com/") || strings.HasPrefix(link, "https://bitbucket.org/") {
			parts := strings.Split(link, "/")
			if len(parts) >= 5 {
				owner := parts[3]
				repo := parts[4]
		
				// Remove the .git suffix if it exists
				if strings.HasSuffix(repo, ".git") {
					repo = strings.TrimSuffix(repo, ".git")
				}
		
				projectType := ""
				if strings.HasPrefix(link, "https://github.com/") {
					projectType = getGitHubProjectType(owner, repo)
				} else if strings.HasPrefix(link, "https://gitlab.com/") {
					projectType = getGitLabProjectType(owner, repo)
				} else if strings.HasPrefix(link, "https://bitbucket.org/") {
					projectType = getBitbucketProjectType(owner, repo)
				}
		
				if projectType != "" {
					latestVersion := getLatestVersion(owner, repo, projectType)
					if latestVersion != "" {
						queryDepsDev(link, projectType, repo, latestVersion)
					}
				}
			}
		}		
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error reading rows:", err)
	}
}

func getGitHubProjectType(owner, repo string) string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: storage.GetGlobalConfig().GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

attempt:
	// 获取仓库内容
	_, dirContent, resp, err := client.Repositories.GetContents(ctx, owner, repo, "", nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			fmt.Println("Error fetching repository contents:", err)
			body := err.Error()
			re := regexp.MustCompile(`rate reset in (\d+)m(\d+)s`)
			matches := re.FindStringSubmatch(body)
			if len(matches) == 3 {
				minutes, _ := strconv.Atoi(matches[1])
				seconds, _ := strconv.Atoi(matches[2])
				waitTime := time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
				fmt.Printf("Rate limit exceeded. Waiting %v to retry...\n", waitTime)
				time.Sleep(waitTime)
				goto attempt
			}
		}
		fmt.Println("Error fetching repository contents:", err)
		return ""
	}

	// 判断项目类型
	for _, file := range dirContent {
		if file.Name != nil {
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
	}
	return ""
}

func getGitLabProjectType(owner, repo string) string {
    url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s%%2F%s/repository/tree", owner, repo)
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Error fetching GitLab repository contents:", err)
        return ""
    }
    defer resp.Body.Close()

    var files []map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&files)

    for _, file := range files {
        if name, ok := file["name"].(string); ok {
            switch name {
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
    }
    return ""
}

func getBitbucketProjectType(owner, repo string) string {
    // Initial API request
    url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/src", owner, repo)

    // Create the request
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        fmt.Println("Error creating request:", err)
        return ""
    }

    // Set headers
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse // Don't follow redirects automatically
        },
    }

    // Send request
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error fetching Bitbucket repository contents:", err)
        return ""
    }
    defer resp.Body.Close()
	var body []byte
    // Check if there's a redirect
    if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
        redirectURL := resp.Header.Get("Location")
        // fmt.Println("Redirecting to:", redirectURL)

        // Fetch from the new location
        req, err = http.NewRequest("GET", redirectURL, nil)
        if err != nil {
            fmt.Println("Error creating new request:", err)
            return ""
        }
		req.Header.Set("User-Agent", "curl/7.68.0")
		resp, err = client.Do(req)
        if err != nil {
            fmt.Println("Error following redirect:", err)
            return ""
        }
        defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
    }
    // Now read the response from the new location
    var result map[string][]map[string]interface{}
	err = json.NewDecoder(strings.NewReader(string(body))).Decode(&result)

	files := result["values"]
	for _, file := range files {
		if path, ok := file["path"].(string); ok {
			switch path {
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
	}

	return ""
}

type VersionInfo struct {
	VersionKey struct {
		Version string `json:"version"`
	} `json:"versionKey"`
	PublishedAt time.Time `json:"publishedAt"`
}

type PackageInfo struct {
	Versions []VersionInfo `json:"versions"`
}

func getLatestVersion(owner, repo, projectType string) string {
	ctx := context.Background()

	// 构造请求URL
	url := fmt.Sprintf("https://api.deps.dev/v3alpha/systems/%s/packages/%s", projectType, repo)

	// 发起HTTP GET请求
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Println("Error fetching package information:", err)
		return ""
	}
	defer resp.Body.Close()

	// 解析响应体
	body, _ := io.ReadAll(resp.Body)
	var result PackageInfo
	json.Unmarshal(body, &result)

	// 寻找最新版本
	var latestVersion string
	var latestDate time.Time
	for _, version := range result.Versions {
		if version.PublishedAt.After(latestDate) {
			latestDate = version.PublishedAt
			latestVersion = version.VersionKey.Version
		}
	}

	return latestVersion
}

func queryDepsDev(link, projectType, projectName, version string) {
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

	err = updateDatabase(link, projectName, info.DependentCount, projectType)
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
}
