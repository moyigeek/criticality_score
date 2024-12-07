/*
 * @Author: 7erry
 * @Date: 2024-08-31 03:50:13
 * @LastEditTime: 2024-12-07 18:43:18
 * @Description:
 */
package git

import (
	"strconv"
	"testing"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/logger"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"

	"github.com/stretchr/testify/require"
)

func TestParseGitRepo(t *testing.T) {
	tests := []struct {
		input    string
		expected Repo
	}{
		{
			input:    "https://github.com/gin-gonic/gin.git",
			expected: Repo{},
		},
		{
			input:    "https://github.com/cider-security-research/cicd-goat.git",
			expected: Repo{},
		},
		{
			input:    "https://github.com/cider-security-research/top-10-cicd-security-risks.git",
			expected: Repo{},
		},
		{
			input:    "https://salsa.debian.org/med-team/kmer.git",
			expected: Repo{},
		},
		{
			input:    "https://gitlab.com/Sasha-Zayets/nx-ci-cd.git",
			expected: Repo{},
		},
	}
	for n, test := range tests {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			u := url.ParseURL(test.input)
			r, err := collector.EzCollect(&u)
			if err != nil {
				t.Fatal(err)
			}
			repo, err := ParseRepo(r)
			if err != nil {
				t.Fatal(err)
			}
			logger.Infof("%++v", *repo)
			//require.Equal(t, test.expected, *repo)
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
			if err != nil {
				t.Fatal(err)
			}
			result, err := GetURL(r)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, test.url, result)
		})
	}
}
