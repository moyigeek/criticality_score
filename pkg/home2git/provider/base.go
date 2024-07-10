package provider

import "regexp"

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
	{Pattern: regexp.MustCompile(`git://\S*`), Confidence: 5000},
	{Pattern: regexp.MustCompile(`https?://\S*\.git`), Confidence: 5000},
	{Pattern: regexp.MustCompile(`https?://github.com/\S*`), Confidence: 4000},
}

func matchGitLink(text string) bool {
	for _, pattern := range gitLinkPatterns {
		if pattern.Pattern.MatchString(text) {
			return true
		}
	}
	return false
}

func getMatchedLinks(text string) []QueryResultItem {
	var links []QueryResultItem
	for _, pattern := range gitLinkPatterns {
		matches := pattern.Pattern.FindAllString(text, -1)
		for _, match := range matches {
			links = append(links, QueryResultItem{
				GitURL:     match,
				Confidence: pattern.Confidence,
			})
		}
	}
	return links
}
