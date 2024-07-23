package provider

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/net/html"
)

func downloadHTML(url string) (*string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: %s", url, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyStr := string(body)
	return &bodyStr, nil
}

func traverseNode(doc *html.Node, matcher func(node *html.Node) (bool, bool)) (nodes []*html.Node) {
	var keep, exit bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		keep, exit = matcher(n)
		if keep {
			nodes = append(nodes, n)
		}
		if exit {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return nodes
}

func browseHTML(url string, depth int, urlVisited *[]string) ([]QueryResultItem, error) {
	var results []QueryResultItem

	if depth == 0 {
		return nil, fmt.Errorf("depth limit reached")
	}

	if slices.Contains(*urlVisited, url) {
		return nil, fmt.Errorf("url already visited")
	}

	*urlVisited = append(*urlVisited, url)

	htmlstr, err := downloadHTML(url)
	if err != nil {
		return nil, err
	}

	dom, err := html.Parse(strings.NewReader(*htmlstr))
	if err != nil {
		return nil, err
	}

	nodesContainingGitLink := traverseNode(dom, func(node *html.Node) (bool, bool) {
		return node.Type == html.TextNode && matchGitLink(node.Data), false
	})

	for _, node := range nodesContainingGitLink {
		mLinks := getMatchedLinks(node.Data, depth)
		results = append(results, mLinks...)
	}

	allAnchorNodes := traverseNode(dom, func(node *html.Node) (bool, bool) {
		return node.Type == html.ElementNode && node.Data == "a", false
	})

	for _, node := range allAnchorNodes {
		var href string
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				href = attr.Val
				break
			}
		}

		// must be same domain
		currentDomain := strings.Split(url, "/")[2]
		if strings.HasPrefix(href, "/") {
			href = strings.Split(url, "/")[0] + "//" + currentDomain + href
		} else if !strings.HasPrefix(href, "http://") &&
			!strings.HasPrefix(href, "https://") &&
			!strings.HasPrefix(href, "//") {
			href = url + href
		}

		if !strings.Contains(href, currentDomain) {
			continue
		}

		rLinks, err := browseHTML(href, depth-1, urlVisited)
		if err != nil {
			continue
		}
		results = append(results, rLinks...)
	}

	return results, nil
}

func QeuryByBrowse(homepage string, packageName string) (*QueryResult, error) {
	links, err := browseHTML(homepage, 2, &[]string{})
	if err != nil {
		return nil, err
	}

	var items []QueryResultItem
	items = append(items, links...)

	for _, item := range items {
		if strings.Contains(item.GitURL, packageName) {
			item.Confidence += 2000
		}
	}

	return &QueryResult{
		Items:    items,
		NeedNext: true,
	}, nil
}
