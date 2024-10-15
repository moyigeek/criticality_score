/*
 * @Author: 7erry
 * @Date: 2024-08-31 03:50:13
 * @LastEditTime: 2024-09-29 16:38:22
 * @Description:
 */
package git

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestGetMetrics(t *testing.T) {
	tests := []struct {
		input    string
		expected RepoMetrics
	}{{
		input: "https://gitee.com/teocloud/teo-docs-search-engine.git",
		expected: RepoMetrics{
			CreatedSince:     time.Date(2024, 2, 7, 23, 31, 55, 0, time.UTC),
			UpdatedSince:     time.Date(2024, 2, 8, 2, 45, 14, 0, time.UTC),
			ContributorCount: 1,
			OrgCount:         1,
			CommitFrequency:  0.17307692307692307,
		},
	},
		{
			input: "https://gitee.com/Open-Brother/pzstudio.git",
			expected: RepoMetrics{
				CreatedSince:     time.Date(2022, 6, 20, 21, 55, 57, 0, time.UTC),
				UpdatedSince:     time.Date(2023, 2, 2, 7, 30, 36, 0, time.UTC),
				ContributorCount: 3,
				OrgCount:         1,
				CommitFrequency:  0,
			},
		},
		{
			input: "https://gitee.com/mirrors/Proxy-Go.git",
			expected: RepoMetrics{
				CreatedSince:     time.Date(2021, 5, 11, 13, 40, 18, 0, time.UTC),
				UpdatedSince:     time.Date(2024, 5, 24, 23, 22, 19, 0, time.UTC),
				ContributorCount: 3,
				OrgCount:         2,
				CommitFrequency:  0.40384615384615385,
			},
		},
		{
			input: "https://gitcode.com/lovinpanda/DirectX.git",
			expected: RepoMetrics{
				CreatedSince:     time.Date(2024, 8, 24, 14, 37, 42, 0, time.UTC),
				UpdatedSince:     time.Date(2024, 9, 26, 18, 20, 37, 0, time.UTC),
				ContributorCount: 2,
				OrgCount:         1,
				CommitFrequency:  0.096154,
			},
		},
		{
			input: "https://gitlab.com/Sasha-Zayets/nx-ci-cd.git",
			expected: RepoMetrics{
				CreatedSince:     time.Date(2024, 8, 28, 19, 18, 24, 0, time.UTC),
				UpdatedSince:     time.Date(2024, 8, 30, 07, 23, 15, 0, time.UTC),
				ContributorCount: 2,
				OrgCount:         1,
				CommitFrequency:  0.326923,
			},
		},
		{
			input: "https://salsa.debian.org/med-team/kmer.git",
			expected: RepoMetrics{
				CreatedSince:     time.Date(2015, 5, 6, 19, 39, 42, 0, time.UTC),
				UpdatedSince:     time.Date(2024, 9, 10, 10, 35, 3, 0, time.UTC),
				ContributorCount: 0,
				OrgCount:         0,
				CommitFrequency:  0,
			},
		},
	}
	for n, test := range tests {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			u := url.ParseURL(test.input)
			r, err := collector.EzCollect(&u)
			err = utils.HandleErr(err, u.URL)
			utils.CheckIfError(err)
			m := GetMetrics(r)
			require.Equal(t, test.expected, *m)
		})
	}
}

func TestParseGitRepo(t *testing.T) {
	tests := []struct {
		input    string
		expected Repo
	}{
		{},
		{},
		{},
		{},
		{},
	}
	for n, test := range tests {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			u := url.ParseURL(test.input)
			r, err := collector.EzCollect(&u)
			err = utils.HandleErr(err, u.URL)
			utils.CheckIfError(err)
			repo := ParseGitRepo(r)
			require.Equal(t, test.expected, *repo)
		})
	}
}

func TestGetLicense(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "https://gitee.com/goldpankit/goldpankit.git",
			expected: []string{"GPL-2.0"},
		}, /*
			{
				input:    "https://github.com/gin-gonic/gin.git",
				expected: []string{"MIT license"},
			},
			{
				input:    "https://bitbucket.org/evolution536/crysearch-memory-scanner.git",
				expected: []string{"MIT license"},
			},*/
	}
	for n, test := range tests {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			u := url.ParseURL(test.input)
			r, err := collector.EzCollect(&u)
			err = utils.HandleErr(err, u.URL)
			utils.CheckIfError(err)
			l := GetLicense(r)
			fmt.Println(l)
			require.Equal(t, test.expected, l)
		})
	}
}

func TestGetLanguages(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "https://gitee.com/goldpankit/goldpankit.git",
			expected: []string{"Vue", "JavaScript", "SCSS", "HTML"},
		}, /*
			{
				input:    "https://github.com/gin-gonic/gin.git",
				expected: []string{"MIT license"},
			},
			{
				input:    "https://github.com/appleboy/gin-jwt.git",
				expected: []string{"MIT license"},
			},*/
	}
	for n, test := range tests {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			u := url.ParseURL(test.input)
			r, err := collector.EzCollect(&u)
			err = utils.HandleErr(err, u.URL)
			utils.CheckIfError(err)
			l := GetLanguages(r)
			require.Equal(t, test.expected, *l)
		})
	}
}

func TestGetURL(t *testing.T) {
	tests := []struct {
		url string
	}{
		{"https://gitee.com/teocloud/teo-docs-search-engine.git"},
		{"https://gitee.com/Open-Brother/pzstudio.git"},
		{"https://gitee.com/mirrors/Proxy-Go.git"},
		{"https://gitcode.com/lovinpanda/DirectX.git"},
		{"https://gitlab.com/Sasha-Zayets/nx-ci-cd.git"},
		{"https://salsa.debian.org/med-team/kmer.git"},
	}
	for n, test := range tests {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			u := url.ParseURL(test.url)
			r, err := collector.EzCollect(&u)
			err = utils.HandleErr(err, u.URL)
			utils.CheckIfError(err)
			result := GetURL(r)
			require.Equal(t, test.url, result)
		})
	}
}
