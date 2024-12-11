package home2git

import (
	"encoding/csv"
	"fmt"
	"os"
	"net/http"
	"time"
	"io"
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var visitedLinks = make(map[string]bool)


func Home2git(flagConfigPath string, outputCSV string, url string){
	err := storage.InitializeDatabase(flagConfigPath)
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	if url == ""{
		fmt.Println("Please provide a LLM URL.")
		return
	}
	var PackageList [][]string
	PackageList, _ = FetchAllLinks()

	outFile, err := os.Create(outputCSV)
	if err != nil {
		fmt.Printf("Failed to create output file %s: %v\n", outputCSV, err)
		return
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)

	if err := writer.Write([]string{"PackageName", "Homepage", "GitHub Repository"}); err != nil {
		fmt.Printf("Error writing header: %v\n", err)
		return
	}
	writer.Flush()

	for i := 0; i < len(PackageList); i++ {
		homepageURL := PackageList[i][0]
		packageName := PackageList[i][1]

		htmlContent, err := DownloadHTML(homepageURL)
		if err != nil {
			continue
		}

		links, _ := FindLinksInHTML(homepageURL, htmlContent, 1)
		githubURL := ProcessHomepage(packageName, links, homepageURL, url)
		res := Check(githubURL)
		if res == "" {
			continue
		}

		if err := writer.Write([]string{packageName, homepageURL, res}); err != nil {
			fmt.Printf("Error writing row: %v\n", err)
			continue
		}
		writer.Flush()
	}

	fmt.Println("Processing complete. Results saved to", outputCSV)
}
func Check(githubURL string)string{
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(githubURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return githubURL
	}
	return ""
}

func FetchAllLinks() ([][]string, error) {
	db, err := storage.GetDatabaseConnection()
	rows, err := db.Query("SELECT package, homepage FROM gentoo_packages union select package, homepage from homebrew_packages union select package, homepage from nix_packages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links [][]string
	for rows.Next() {
		var packageName, link string
		if err := rows.Scan(&packageName, &link); err != nil {
			return nil, err
		}
		links = append(links, []string{link, packageName})
	}
	fmt.Println("Fetched", len(links), "links")
	return links, nil
}

func DownloadHTML(url string) (string, error) {
	visitedLinks[url] = true

	// 创建一个 http.Client 对象
	client := &http.Client{}

	// 创建一个 http.Request 对象
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// 设置 User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func FindLinksInHTML(url string, htmlContent string, depth int) ([]string, error) {
	// if depth == 0 {
	// 	return nil, nil
	// }
	// doc, err := html.Parse(strings.NewReader(htmlContent))
	// if err != nil {
	// 	return nil, fmt.Errorf("Error parsing HTML: %v", err)
	// }

	var links []string
	// var currentDomain = strings.Split(url, "/")[2] // Get the domain from the URL
	for _, pattern := range gitLinkPatterns {
		if matches := pattern.FindAllString(htmlContent, -1); matches != nil {
			for _, match := range matches {
				if !visitedLinks[match] {
					visitedLinks[match] = true
					links = append(links, match)
				}
			}
		}
	}
	// var traverse func(*html.Node)
	// traverse = func(n *html.Node) {
	// 	if n.Type == html.TextNode {
	// 		for _, pattern := range gitLinkPatterns {
	// 			if matches := pattern.FindAllString(htmlContent, -1); matches != nil {
	// 				for _, match := range matches {
	// 					if !visitedLinks[match] {
	// 						visitedLinks[match] = true
	// 						links = append(links, match)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
	// 		traverse(c)
	// 	}
	// 	if n.Type == html.ElementNode && n.Data == "a" {
	// 		for _, a := range n.Attr {
	// 			if a.Key == "href" {
	// 				href := strings.TrimSpace(a.Val)
	// 				if strings.HasPrefix(href, "/") {
	// 					href = strings.Split(url, "/")[0] + "//" + currentDomain + href
	// 				} else if !strings.HasPrefix(href, "http://") &&
	// 					!strings.HasPrefix(href, "https://") &&
	// 					!strings.HasPrefix(href, "//") {
	// 					href = url + href
	// 				} else {
	// 					linkDomain := strings.Split(href, "/")[2]
	// 					if linkDomain != currentDomain {
	// 						continue // Skip non-same-origin full URLs
	// 					}
	// 				}
	// 				if !visitedLinks[href] {
	// 					visitedLinks[href] = true
	// 					deeperContent, err := DownloadHTML(href)
	// 					if err == nil && depth > 1 {
	// 						FindLinksInHTML(href, deeperContent, depth-1)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// traverse(doc)
	return links, nil
}

func FindGitRepository(homepageURL string, links []string, packageName string, attempts int, url string) string {
	var prompt string
	linksString := strings.Join(links, ", ")
	if len(links) > 0 {
		prompt = fmt.Sprintf(PROMPT["home2git_link"], linksString, packageName)
	} else {
		prompt = fmt.Sprintf(PROMPT["home2git_nolink"], homepageURL)
	}
	fullResponse := InvokeModel(prompt, attempts, url)
	if strings.Contains(fullResponse, "does not exist") {
		if attempts < 3 {
			return FindGitRepository(homepageURL, links, packageName, attempts+1, url)
		}
		return "does not exist"
	} else if strings.Contains(fullResponse, "URL is:") {
		potentialURL := strings.Split(strings.Split(fullResponse, "URL is:")[1], "\n")[0]
		patten := gitLinkPatterns[0]
		if matches := patten.FindStringSubmatch(potentialURL); len(matches) > 0 {
			return matches[0]
		}
		return potentialURL
	}
	return "does not exist"
}

func ProcessHomepage(packageName string, links []string, homepageURL string,url string) string {
	githubURL := FindGitRepository(homepageURL, links, packageName, 3, url)
	return githubURL
}
