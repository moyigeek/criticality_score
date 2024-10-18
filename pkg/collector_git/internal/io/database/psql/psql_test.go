package psql

import (
	"strconv"
	"testing"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"
)

func TestInsertTable(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "https://github.com/Oztechan/CCC.git",
		},
	}
	for n, test := range tests {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			u := url.ParseURL(test.input)
			r, err := collector.EzCollect(&u)
			utils.CheckIfError(err)
			repo := git.ParseGitRepo(r)
			db := InitDB()
			CreateTable(db)
			InsertTable(db, &database.GitMetrics{
				Name:             repo.Name,
				Owner:            repo.Owner,
				Source:           repo.Source,
				URL:              repo.URL,
				Ecosystems:       repo.Ecosystems,
				CreatedSince:     repo.Metrics.CreatedSince,
				UpdatedSince:     repo.Metrics.UpdatedSince,
				ContributorCount: repo.Metrics.ContributorCount,
				OrgCount:         repo.Metrics.OrgCount,
				CommitFrequency:  repo.Metrics.CommitFrequency,
			})
		})
	}
}
