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

	"database/sql"
	"io/ioutil"

	"github.com/google/go-github/v63/github"
	_ "github.com/lib/pq" // Assuming PostgreSQL, adjust as needed
	"golang.org/x/oauth2"
)

type DependentInfo struct {
	DependentCount         int `json:"dependentCount"`
	DirectDependentCount   int `json:"directDependentCount"`
	IndirectDependentCount int `json:"indirectDependentCount"`
}

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

func updateDatabase(link, projectName string, dependentCount int, config Config) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("UPDATE git_metrics SET depsdev_count = $1 WHERE git_link = $2", dependentCount, link) // Adjust the WHERE clause as needed
	return err
}

func Run(configPath string) {
	// Load configuration
	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Connect to the database
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
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
		if strings.HasPrefix(link, "https://github.com/") {
			parts := strings.Split(link, "/")
			if len(parts) >= 5 {
				owner := parts[3]
				repo := parts[4]

				// Remove the .git suffix if it exists
				if strings.HasSuffix(repo, ".git") {
					repo = strings.TrimSuffix(repo, ".git")
				}

				projectType := getProjectType(owner, repo, config)
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

func getProjectType(owner, repo string, config Config) string {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.GitHubToken},
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

	// Update database with dependent count
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	err = updateDatabase(link, projectName, info.DependentCount, config)
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		return
	}
}
