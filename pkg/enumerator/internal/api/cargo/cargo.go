package cargo

type Response struct {
	Crates []Crate `json:"crates"`
	Meta   Meta    `json:"meta"`
}

type Crate struct {
	Badges           []string    `json:"badges"`
	Categories       interface{} `json:"categories"`
	CreatedAt        string      `json:"created_at"`
	DefaultVersion   string      `json:"default_version"`
	Description      string      `json:"description"`
	Documentation    *string     `json:"documentation"`
	Downloads        int64       `json:"downloads"`
	ExactMatch       bool        `json:"exact_match"`
	Homepage         *string     `json:"homepage"`
	ID               string      `json:"id"`
	Keywords         interface{} `json:"keywords"`
	Links            Links       `json:"links"`
	MaxStableVersion string      `json:"max_stable_version"`
	MaxVersion       string      `json:"max_version"`
	Name             string      `json:"name"`
	NewestVersion    string      `json:"newest_version"`
	RecentDownloads  int64       `json:"recent_downloads"`
	Repository       string      `json:"repository"`
	UpdatedAt        string      `json:"updated_at"`
	Versions         interface{} `json:"versions"`
	Yanked           bool        `json:"yanked"`
}

type Links struct {
	OwnerTeam           string `json:"owner_team"`
	OwnerUser           string `json:"owner_user"`
	Owners              string `json:"owners"`
	ReverseDependencies string `json:"reverse_dependencies"`
	VersionDownloads    string `json:"version_downloads"`
	Versions            string `json:"versions"`
}

type Meta struct {
	NextPage string      `json:"next_page"`
	PrevPage interface{} `json:"prev_page"`
	Total    int64       `json:"total"`
}
