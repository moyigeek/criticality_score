package home2git

import (
	"fmt"
	"io"
	"net/http"
)

var visitedLinks = make(map[string]bool)

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
