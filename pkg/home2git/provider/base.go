package provider

import (
	"fmt"
	"regexp"
)

type QueryResultItem struct {
	GitURL     string
	Confidence int // Confidence level of the result, 0-10^4
}

type QueryResult struct {
	Items    []QueryResultItem
	NeedNext bool
}

type QueryFunction func(homePage string, packageName string) (*QueryResult, error)

var gitLinkPatterns = []struct {
	Pattern    *regexp.Regexp
	Confidence int
}{
	{Pattern: regexp.MustCompile(`git://[^\s'";,#\\]+`), Confidence: 5000},
	{Pattern: regexp.MustCompile(`https?://[^\s'";,#\\]+\.git`), Confidence: 5000},
	{Pattern: regexp.MustCompile(`https?://github.com/[^\s'";,#\\]+`), Confidence: 4000},
}

func matchGitLink(text string) bool {
	for _, pattern := range gitLinkPatterns {
		if pattern.Pattern.MatchString(text) {
			return true
		}
	}
	return false
}

func getMatchedLinks(text string, depth int) []QueryResultItem {
	var links []QueryResultItem
	for _, pattern := range gitLinkPatterns {
		matches := pattern.Pattern.FindAllString(text, -1)
		for _, match := range matches {
			fmt.Println(depth)
			confidence := pattern.Confidence / (3 - depth)
			links = append(links, QueryResultItem{
				GitURL:     match,
				Confidence: confidence,
			})
		}
	}
	return links
}
