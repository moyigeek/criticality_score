/*
 * @Author: 7erry
 * @Date: 2024-10-18 20:26:31
 * @LastEditTime: 2025-01-07 19:03:45
 * @Description:
 */
package sqlite

import (
	"strconv"
	"testing"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/database"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/parser/git"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/parser/url"
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
			repo, err := git.ParseRepo(r)
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
				License:          repo.License,
				Ecosystems:       repo.Ecosystems,
				Languages:        repo.Languages,
				CreatedSince:     repo.CreatedSince,
				UpdatedSince:     repo.UpdatedSince,
				ContributorCount: repo.ContributorCount,
				OrgCount:         repo.OrgCount,
				CommitFrequency:  repo.CommitFrequency,
				NeedUpdate:       false,
			})
		})
	}
}
