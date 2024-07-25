package home2git

import "regexp"

// 定义匹配Git链接的正则表达式，确保正确终止
var gitLinkPatterns = []*regexp.Regexp{
	regexp.MustCompile(`https?://github\.com/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://gitlab\.com/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://bitbucket\.org/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://gitlab\.org/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://gitee\.com/[^,)#\s'"]+`),
}

// GitLink 包含链接和匹配的正则表达式
type GitLink struct {
	URL     string
	Pattern *regexp.Regexp
}

// CheckIfGitLink 检查一个字符串是否为Git链接
func CheckIfGitLink(url string) *GitLink {
	for _, pattern := range gitLinkPatterns {
		if pattern.MatchString(url) {
			return &GitLink{URL: url, Pattern: pattern}
		}
	}
	return nil
}
