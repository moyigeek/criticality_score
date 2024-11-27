/*
 * @Author: 7erry
 * @Date: 2024-09-29 14:41:35
 * @LastEditTime: 2024-11-27 20:09:54
 * @Description:
 */
package collector

import (
	"testing"

	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/stretchr/testify/require"
)

func TestCollect(t *testing.T) {
	tests := []struct {
		input    string
		expected error
	}{
		{input: "https://github.com/gin-gonic/gin123456", expected: nil},
		{input: "https://gitee.com/teocloud/teo-docs-search-engine.git", expected: nil},
		{input: "https://gitlab.com/Sasha-Zayets/nx-ci-cd.git", expected: nil},
		{input: "https://salsa.debian.org/med-team/kmer.git", expected: nil},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			u := url.ParseURL(test.input)
			_, err := Collect(&u)
			require.Equal(t, test.expected, err)
		})
	}
}
