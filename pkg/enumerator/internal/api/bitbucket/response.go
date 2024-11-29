package bitbucket

// Request

type Response struct {
	Next    string  `json:"next"`
	Pagelen int64   `json:"pagelen"`
	Values  []Value `json:"values"`
}

type Value struct {
	CreatedOn        string           `json:"created_on"`
	Description      string           `json:"description"`
	ForkPolicy       string           `json:"fork_policy"`
	FullName         string           `json:"full_name"`
	HasIssues        bool             `json:"has_issues"`
	HasWiki          bool             `json:"has_wiki"`
	IsPrivate        bool             `json:"is_private"`
	Language         string           `json:"language"`
	Links            ValueLinks       `json:"links"`
	Mainbranch       Mainbranch       `json:"mainbranch"`
	Name             string           `json:"name"`
	OverrideSettings OverrideSettings `json:"override_settings"`
	Owner            Owner            `json:"owner"`
	Parent           interface{}      `json:"parent"`
	Project          Project          `json:"project"`
	SCM              string           `json:"scm"`
	Size             int64            `json:"size"`
	Slug             string           `json:"slug"`
	Type             string           `json:"type"`
	UpdatedOn        string           `json:"updated_on"`
	UUID             string           `json:"uuid"`
	Website          string           `json:"website"`
	Workspace        Workspace        `json:"workspace"`
}

type ValueLinks struct {
	Avatar       PurpleAvatar `json:"avatar"`
	Branches     Branches     `json:"branches"`
	Clone        []Clone      `json:"clone"`
	Commits      Commits      `json:"commits"`
	Downloads    Downloads    `json:"downloads"`
	Forks        Forks        `json:"forks"`
	Hooks        Hooks        `json:"hooks"`
	HTML         PurpleHTML   `json:"html"`
	Issues       *Issues      `json:"issues,omitempty"`
	Pullrequests Pullrequests `json:"pullrequests"`
	Self         PurpleSelf   `json:"self"`
	Source       Source       `json:"source"`
	Tags         Tags         `json:"tags"`
	Watchers     Watchers     `json:"watchers"`
}

type PurpleAvatar struct {
	Href string `json:"href"`
}

type Branches struct {
	Href string `json:"href"`
}

type Clone struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type Commits struct {
	Href string `json:"href"`
}

type Downloads struct {
	Href string `json:"href"`
}

type Forks struct {
	Href string `json:"href"`
}

type PurpleHTML struct {
	Href string `json:"href"`
}

type Hooks struct {
	Href string `json:"href"`
}

type Issues struct {
	Href string `json:"href"`
}

type Pullrequests struct {
	Href string `json:"href"`
}

type PurpleSelf struct {
	Href string `json:"href"`
}

type Source struct {
	Href string `json:"href"`
}

type Tags struct {
	Href string `json:"href"`
}

type Watchers struct {
	Href string `json:"href"`
}

type Mainbranch struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type OverrideSettings struct {
	BranchingModel       bool `json:"branching_model"`
	DefaultMergeStrategy bool `json:"default_merge_strategy"`
}

type Owner struct {
	AccountID   string     `json:"account_id"`
	DisplayName string     `json:"display_name"`
	Links       OwnerLinks `json:"links"`
	Nickname    string     `json:"nickname"`
	Type        string     `json:"type"`
	Username    *string    `json:"username,omitempty"`
	UUID        string     `json:"uuid"`
}

type OwnerLinks struct {
	Avatar FluffyAvatar `json:"avatar"`
	HTML   FluffyHTML   `json:"html"`
	Self   FluffySelf   `json:"self"`
}

type FluffyAvatar struct {
	Href string `json:"href"`
}

type FluffyHTML struct {
	Href string `json:"href"`
}

type FluffySelf struct {
	Href string `json:"href"`
}

type Project struct {
	Key   string       `json:"key"`
	Links ProjectLinks `json:"links"`
	Name  string       `json:"name"`
	Type  string       `json:"type"`
	UUID  string       `json:"uuid"`
}

type ProjectLinks struct {
	Avatar TentacledAvatar `json:"avatar"`
	HTML   TentacledHTML   `json:"html"`
	Self   TentacledSelf   `json:"self"`
}

type TentacledAvatar struct {
	Href string `json:"href"`
}

type TentacledHTML struct {
	Href string `json:"href"`
}

type TentacledSelf struct {
	Href string `json:"href"`
}

type Workspace struct {
	Links WorkspaceLinks `json:"links"`
	Name  string         `json:"name"`
	Slug  string         `json:"slug"`
	Type  string         `json:"type"`
	UUID  string         `json:"uuid"`
}

type WorkspaceLinks struct {
	Avatar StickyAvatar `json:"avatar"`
	HTML   StickyHTML   `json:"html"`
	Self   StickySelf   `json:"self"`
}

type StickyAvatar struct {
	Href string `json:"href"`
}

type StickyHTML struct {
	Href string `json:"href"`
}

type StickySelf struct {
	Href string `json:"href"`
}
