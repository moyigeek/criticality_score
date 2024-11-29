package gitlab

type Response []RequestElement

// Request
type RequestElement struct {
	AvatarURL         *string   `json:"avatar_url"`
	CreatedAt         string    `json:"created_at"`
	DefaultBranch     string    `json:"default_branch"`
	Description       *string   `json:"description"`
	ForksCount        int64     `json:"forks_count"`
	HTTPURLToRepo     string    `json:"http_url_to_repo"`
	ID                int64     `json:"id"`
	LastActivityAt    string    `json:"last_activity_at"`
	Name              string    `json:"name"`
	NameWithNamespace string    `json:"name_with_namespace"`
	Namespace         Namespace `json:"namespace"`
	Path              string    `json:"path"`
	PathWithNamespace string    `json:"path_with_namespace"`
	ReadmeURL         *string   `json:"readme_url"`
	SSHURLToRepo      string    `json:"ssh_url_to_repo"`
	StarCount         int64     `json:"star_count"`
	TagList           []string  `json:"tag_list"`
	Topics            []string  `json:"topics"`
	WebURL            string    `json:"web_url"`
}

type Namespace struct {
	AvatarURL *string `json:"avatar_url"`
	FullPath  string  `json:"full_path"`
	ID        int64   `json:"id"`
	Kind      string  `json:"kind"`
	Name      string  `json:"name"`
	ParentID  *int64  `json:"parent_id"`
	Path      string  `json:"path"`
	WebURL    string  `json:"web_url"`
}
