package packagist

type Response []RequestElement

// Request
type RequestElement struct {
	Authors          string      `json:"authors"`
	BugTrackerURI    *string     `json:"bug_tracker_uri"`
	ChangelogURI     *string     `json:"changelog_uri"`
	DocumentationURI *string     `json:"documentation_uri"`
	Downloads        int64       `json:"downloads"`
	FundingURI       interface{} `json:"funding_uri"`
	GemURI           string      `json:"gem_uri"`
	HomepageURI      *string     `json:"homepage_uri"`
	Info             string      `json:"info"`
	Licenses         []string    `json:"licenses"`
	MailingListURI   *string     `json:"mailing_list_uri"`
	Metadata         Metadata    `json:"metadata"`
	Name             string      `json:"name"`
	Platform         string      `json:"platform"`
	ProjectURI       string      `json:"project_uri"`
	SHA              string      `json:"sha"`
	SourceCodeURI    *string     `json:"source_code_uri"`
	Version          string      `json:"version"`
	VersionDownloads int64       `json:"version_downloads"`
	WikiURI          *string     `json:"wiki_uri"`
}

type Metadata struct {
	AllowedPushHost *string `json:"allowed_push_host,omitempty"`
	ChangelogURI    *string `json:"changelog_uri,omitempty"`
	FundingURI      *string `json:"funding-uri,omitempty"`
	HomepageURI     string  `json:"homepage_uri"`
	SourceCodeURI   *string `json:"source_code_uri,omitempty"`
}
