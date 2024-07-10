package provider

import "fmt"

func DetermineGitLink(homepage string, _ string) (*QueryResult, error) {
	results := getMatchedLinks(homepage)

	if len(results) == 0 {
		return nil, fmt.Errorf("no git link found")
	}

	return &QueryResult{
		Items:    results,
		NeedNext: false,
	}, nil
}
