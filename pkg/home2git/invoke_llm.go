package home2git

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Response struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// findGitRepository tries to find a Git repository for the given homepage URL. If none is found after three attempts, it returns "does not exist".
func FindGitRepository(homepageURL string, links []string, packageName string, attempts int) string {
	var prompt string
	linksString := strings.Join(links, ", ")
	if len(links) > 0 {
		prompt = fmt.Sprintf("Given the list of repository links [%s] and the homepage URL '%s', select the most likely repository link that matches the homepage. If a direct match is found, return the URL in the format 'URL is: [matched_url]'. If no direct match is found, check other repositories on GitHub, GitLab, or Gitee. If a related repository exists, respond with 'URL is: [url]'. If no relevant repository can be identified, respond with 'does not exist'.", linksString, packageName)
	} else {
		prompt = fmt.Sprintf("Check if there is a git repository for %s hosted on platforms like GitHub, GitLab, or Gitee. If it exists, respond in the format 'URL is: [url]'. If no repository exists, respond with 'does not exist'.", homepageURL)
	}

	url := "http://222.20.126.129:11434/api/generate"
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	data := map[string]interface{}{
		"model":  "llama3:70b",
		"prompt": prompt,
		"stream": true,
	}
	jsonData, _ := json.Marshal(data)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "HTTP request failed"
	}
	defer response.Body.Close()

	reader := bufio.NewReader(response.Body)
	fullResponse := ""

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		var jsonResponse Response
		if err := json.Unmarshal(line, &jsonResponse); err != nil {
			return "JSON decoding error"
		}
		fullResponse += jsonResponse.Response
		if jsonResponse.Done {
			break
		}
	}

	if strings.Contains(fullResponse, "does not exist") {
		if attempts < 3 {
			return FindGitRepository(homepageURL, links, packageName, attempts+1)
		}
		return "does not exist"
	} else if strings.Contains(fullResponse, "URL is:") {
		potentialURL := strings.Split(strings.Split(fullResponse, "URL is:")[1], "\n")[0]
		patten := gitLinkPatterns[0]
		if matches := patten.FindStringSubmatch(potentialURL); len(matches) > 0 {
			// 如果找到匹配，返回匹配的URL
			return matches[0]
		}
		return potentialURL
	}
	return "does not exist"
}

func ProcessHomepage(packageName string, links []string, homepageURL string) string {
	// 调用已经存在的 FindGitRepository 函数，并直接返回其结果
	githubURL := FindGitRepository(homepageURL, links, packageName, 3)
	return githubURL
}
