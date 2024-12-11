package home2git

import "regexp"


var gitLinkPatterns = []*regexp.Regexp{
	regexp.MustCompile(`https?://github\.com/([A-Za-z0-9]+)/([A-Za-z0-9]+)`),
	regexp.MustCompile(`https?://gitlab\.com/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://bitbucket\.org/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://gitlab\.org/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://gitee\.com/[^,)#\s'"]+`),
}

type GitLink struct {
	URL     string
	Pattern *regexp.Regexp
}

func CheckIfGitLink(url string) *GitLink {
	for _, pattern := range gitLinkPatterns {
		if pattern.MatchString(url) {
			return &GitLink{URL: url, Pattern: pattern}
		}
	}
	return nil
}

var PROMPT = map[string]string{
	"home2git_link": "Given the list of repository links [%s] and the homepage URL '%s', select the most likely repository link that matches the homepage. If a direct match is found, return the URL in the format 'URL is: [matched_url]'. If no direct match is found, check other repositories on GitHub, GitLab, or Gitee. If a related repository exists, respond with 'URL is: [url]'. If no relevant repository can be identified, respond with 'does not exist'.",
	"home2git_nolink": "Check if there is a git repository for %s hosted on platforms like GitHub, GitLab, or Gitee. If it exists, respond in the format 'URL is: [url]'. If no repository exists, respond with 'does not exist'.",
}