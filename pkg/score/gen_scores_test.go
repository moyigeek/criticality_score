package score

import (
	"math"
	"testing"
	"time"
)

func TestCalculateGitMetadataScore(t *testing.T) {
	createdSince := time.Now().AddDate(-1, 0, 0)
	updatedSince := time.Now().AddDate(0, -6, 0)
	contributorCount := 100
	commitFrequency := 50.0
	orgCount := 10

	gitMetadata := &GitMetadata{
		CreatedSince:     &createdSince,
		UpdatedSince:     &updatedSince,
		ContributorCount: &contributorCount,
		CommitFrequency:  &commitFrequency,
		Org_Count:        &orgCount,
	}

	expectedScore := 1.8338566950193282
	gitMetadata.CalculateGitMetadataScore()

	if gitMetadata.GitMetadataScore != expectedScore {
		t.Errorf("Expected score %v, but got %v", expectedScore, gitMetadata.GitMetadataScore)
	}
}

func TestCalculateDistScore(t *testing.T) {
	distScore := &DistScore{
		DistImpact:   0.5,
		DistPageRank: 0.5,
	}

	expectedScore := (weights["distScore"]["dist_impact"] * distScore.DistImpact) + (weights["distScore"]["dist_pagerank"] * distScore.DistPageRank)
	distScore.CalculateDistScore()

	if distScore.DistScore != expectedScore {
		t.Errorf("Expected score %v, but got %v", expectedScore, distScore.DistScore)
	}
}

func TestLogNormalize(t *testing.T) {
	value := 10.0
	threshold := 100.0

	expected := math.Log(value+1) / math.Log(math.Max(value, threshold)+1)
	actual := LogNormalize(value, threshold)

	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}
}
