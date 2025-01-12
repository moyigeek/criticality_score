package gitlab

import "time"

type Response []GitLabResposneElement

type GitLabResposneElement struct {
	ID                int64          `json:"id"`
	Description       *string        `json:"description"`
	Name              string         `json:"name"`
	NameWithNamespace string         `json:"name_with_namespace"`
	Path              string         `json:"path"`
	PathWithNamespace string         `json:"path_with_namespace"`
	CreatedAt         time.Time      `json:"created_at"`
	DefaultBranch     *DefaultBranch `json:"default_branch,omitempty"`
	TagList           []string       `json:"tag_list"`
	Topics            []string       `json:"topics"`
	SSHURLToRepo      string         `json:"ssh_url_to_repo"`
	HTTPURLToRepo     string         `json:"http_url_to_repo"`
	WebURL            string         `json:"web_url"`
	ReadmeURL         *string        `json:"readme_url"`
	ForksCount        *int64         `json:"forks_count,omitempty"`
	AvatarURL         *string        `json:"avatar_url"`
	StarCount         int64          `json:"star_count"`
	LastActivityAt    time.Time      `json:"last_activity_at"`
	Namespace         Namespace      `json:"namespace"`
}

type Namespace struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	Path      string  `json:"path"`
	Kind      Kind    `json:"kind"`
	FullPath  string  `json:"full_path"`
	ParentID  *int64  `json:"parent_id"`
	AvatarURL *string `json:"avatar_url"`
	WebURL    string  `json:"web_url"`
}

type DefaultBranch string

const (
	Deprecated     DefaultBranch = "deprecated"
	Dev            DefaultBranch = "dev"
	Develop        DefaultBranch = "develop"
	Main           DefaultBranch = "main"
	Master         DefaultBranch = "master"
	Next           DefaultBranch = "next"
	Production     DefaultBranch = "production"
	Release        DefaultBranch = "release"
	The123StableZh DefaultBranch = "12-3-stable-zh"
)

type Kind string

const (
	Group Kind = "group"
	User  Kind = "user"
)
