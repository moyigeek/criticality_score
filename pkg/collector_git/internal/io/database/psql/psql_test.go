/*
 * @Author: 7erry
 * @Date: 2024-10-18 20:26:31
 * @LastEditTime: 2024-11-27 21:24:10
 * @Description:
 */
package psql

import (
	"strconv"
	"testing"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
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
			if err != nil {
				panic(err)
			}
			repo, err := git.ParseGitRepo(r)
			if err != nil {
				t.Fatal(err)
			}
			db, err := InitDB()
			if err != nil {
				panic(err)
			}
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
