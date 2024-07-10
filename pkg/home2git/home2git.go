package home2git

import (
	"fmt"

	"github.com/HUSTSeclab/criticality_score/pkg/home2git/provider"
)

var providers = []provider.QueryFunction{
	provider.DetermineGitLink,
	provider.QeuryByBrowse,
}

func HomepageToGit(homepage string, packageName string) (*provider.QueryResultItem, error) {
	var results []provider.QueryResultItem
	for _, p := range providers {
		result, err := p(homepage, packageName)
		if err != nil {
			continue
		}

		results = append(results, result.Items...)

		if !result.NeedNext {
			break
		}
	}

	// return the most confident result
	maxConfidence := 0
	var maxResult *provider.QueryResultItem
	for _, result := range results {
		if result.Confidence > maxConfidence {
			maxConfidence = result.Confidence
			maxResult = &result
		}
	}

	if maxResult == nil {
		return nil, fmt.Errorf("no git link found")
	}

	return maxResult, nil
}
