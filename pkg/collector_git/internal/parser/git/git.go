package git

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/logger"
	parser "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type Repo struct {
	Name    string
	Owner   string
	Source  string
	URL     string
	License string
	// is_maintained bool
	Languages  []string
	Ecosystems string
	Metrics    *RepoMetrics
}

func NewRepo() Repo {
	return Repo{
		Name:       "",
		Owner:      "",
		Source:     "",
		URL:        "",
		License:    "",
		Languages:  []string{},
		Ecosystems: "",
		Metrics:    &RepoMetrics{},
	}
}

type RepoMetrics struct {
	CreatedSince     time.Time
	UpdatedSince     time.Time
	ContributorCount int
	OrgCount         int
	CommitFrequency  float64
}

func NewRepoMetrics() RepoMetrics {
	return RepoMetrics{
		CreatedSince:     time.Time{},
		UpdatedSince:     time.Time{},
		ContributorCount: 0,
		OrgCount:         0,
		CommitFrequency:  0,
	}
}

func GetBlobs(r *git.Repository) (*[]*object.Blob, error) {
	bIter, err := r.BlobObjects()

	if err != nil {
		return nil, err
	}

	blobs := make([]*object.Blob, 0)

	err = bIter.ForEach(func(b *object.Blob) error {
		// fmt.Println(b)
		blobs = append(blobs, b)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &blobs, nil
}

func GetBranches(r *git.Repository) (*[]*plumbing.Reference, error) {
	rIter, err := r.Branches()

	if err != nil {
		return nil, err
	}

	refs := make([]*plumbing.Reference, 0)
	err = rIter.ForEach(func(r *plumbing.Reference) error {
		// fmt.Println(r)
		refs = append(refs, r)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &refs, nil
}

func GetCommits(r *git.Repository) (*[]*object.Commit, error) {
	cIter, err := r.CommitObjects()

	if err != nil {
		return nil, err
	}

	commits := make([]*object.Commit, 0)
	err = cIter.ForEach(func(c *object.Commit) error {
		// fmt.Println(c)
		commits = append(commits, c)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &commits, nil
}

func GetConfig(r *git.Repository) (*gitconfig.Config, error) {
	c, err := r.Config()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetObjects(r *git.Repository) (*[]*object.Object, error) {
	oIter, err := r.Objects()

	if err != nil {
		return nil, err
	}

	objs := make([]*object.Object, 0)
	err = oIter.ForEach(func(o object.Object) error {
		// fmt.Println(o)
		objs = append(objs, &o)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &objs, nil
}

func GetReferences(r *git.Repository) (*[]*plumbing.Reference, error) {
	rIter, err := r.References()

	if err != nil {
		return nil, err
	}

	refs := make([]*plumbing.Reference, 0)
	err = rIter.ForEach(func(r *plumbing.Reference) error {
		// fmt.Println(r)
		refs = append(refs, r)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &refs, nil
}

func GetRemotes(r *git.Repository) (*[]*git.Remote, error) {
	remotes, err := r.Remotes()

	if err != nil {
		return nil, err
	}

	return &remotes, nil
}

func GetTags(r *git.Repository) (*[]*object.Tag, error) {
	tIter, err := r.TagObjects()

	if err != nil {
		return nil, err
	}

	tags := make([]*object.Tag, 0)
	err = tIter.ForEach(func(t *object.Tag) error {
		// fmt.Println(t)
		tags = append(tags, t)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &tags, nil
}

func GetTagRefs(r *git.Repository) (*[]*plumbing.Reference, error) {
	rIter, err := r.Tags()

	if err != nil {
		return nil, err
	}

	refs := make([]*plumbing.Reference, 0)
	err = rIter.ForEach(func(r *plumbing.Reference) error {
		// fmt.Println(r)
		refs = append(refs, r)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &refs, nil
}

func GetTrees(r *git.Repository) (*[]*object.Tree, error) {
	tIter, err := r.TreeObjects()

	if err != nil {
		return nil, err
	}

	trees := make([]*object.Tree, 0)
	err = tIter.ForEach(func(t *object.Tree) error {
		// fmt.Println(t)
		trees = append(trees, t)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &trees, nil
}

func GetWorkTree(r *git.Repository) (*git.Worktree, error) {
	wt, err := r.Worktree()

	if err != nil {
		return nil, err
	}

	return wt, nil
}

func GetURL(r *git.Repository) (string, error) {
	//? In most cases, the Remote URLs of Git Fetch and Git Push are the same, but we take the former one
	remotes, err := GetRemotes(r)

	if err != nil {
		logger.Error(err)
		return "", err
	}

	if len(*remotes) == 0 {
		return "", nil
	}

	var u string

	if len((*remotes)[0].Config().URLs) > 0 {
		u = (*remotes)[0].Config().URLs[0]
	}

	for _, remote := range *remotes {
		if remote.Config().Name == parser.DEFAULT_REMOTE_NAME {
			if len(remote.Config().URLs) > 0 {
				u = remote.Config().URLs[0]
				break
			}
		}
	}

	return u, nil
}

func GetLicense(r *git.Repository) (string, error) {

	// ToDO https://github.com/licensee/licensee

	//* Maybe multiple licenses
	//* l := []string{}
	l := parser.UNKNOWN_LICENSE

	ref, err := r.Head()
	if err != nil {
		return "", err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	f, err := tree.File("LICENSE")
	if err != nil {
		f, err = tree.File("LICENSE.md")
		if err != nil {
			f, err = tree.File("LICENSE.txt")
		}
	}
	if err == nil {
		content, err := f.Lines()
		if err != nil {
			return "", err
		}

		for k, v := range parser.LICENSE_KEYWORD {
			if strings.Contains(content[0], k) {
				//* l = append(l, v)
				l = v
			}
		}
	}
	return l, nil
}

//! GetLanguages and GetEcosystem could be merged into one function if needed

func GetLanguages(r *git.Repository) (*[]string, error) {
	l := make(map[string]int, 0)
	ref, err := r.Head()
	if err != nil {
		return nil, err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	fIter := tree.Files()
	err = fIter.ForEach(func(f *object.File) error {
		filename := filepath.Base(f.Name)
		v, ok := parser.LANGUAGE_FILENAMES[filename]
		if ok {
			l[v] += 1
		} else {
			ex := filepath.Ext(f.Name)
			v, ok = parser.LANGUAGE_EXTENSIONS[ex]
			if ok {
				l[v] += 1
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	languages := make([]string, 0)
	for k, v := range l {
		if v >= parser.LANGUAGE_THRESHOLD {
			languages = append(languages, k)
		}
	}

	return &languages, nil
}

func GetEcosystem(r *git.Repository) (string, error) {
	eco := ""
	ref, err := r.Head()
	if err != nil {
		return "", err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	//* emap := map[string]int{}
	fIter := tree.Files()
	err = fIter.ForEach(func(f *object.File) error {
		filename := filepath.Base(f.Name)
		v, ok := parser.ECOSYSTEM_MAP[filename]
		if ok {
			// emap[v] += 1
			eco = v
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	return eco, nil
}

func GetMetrics(r *git.Repository) (*RepoMetrics, error) {
	cIter, err := r.Log(&git.LogOptions{
		//* From:  ref.Hash(),
		All:   true,
		Since: &parser.BEGIN_TIME,
		Until: &parser.END_TIME,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, err
	}

	metrics := NewRepoMetrics()
	contributors := make(map[string]int, 0)
	orgs := make(map[string]int, 0)
	var commit_count float64 = 0

	latest_commit, err := cIter.Next()
	if err != nil {
		return nil, err
	}

	author := fmt.Sprintf("%s(%s)", latest_commit.Author.Name, latest_commit.Author.Email)
	e := strings.Split(latest_commit.Author.Email, "@")
	org := e[len(e)-1]

	metrics.UpdatedSince = latest_commit.Committer.When
	contributors[author]++
	orgs[org]++

	if latest_commit.Author.When.After(parser.LAST_YEAR) {
		commit_count++
	}

	flag := true
	created_since := latest_commit.Committer.When

	err = cIter.ForEach(func(c *object.Commit) error {
		author := fmt.Sprintf("%s(%s)", c.Author.Name, c.Author.Email)
		e = strings.Split(c.Author.Email, "@")
		org := e[len(e)-1]

		//! It made sense that this if statement is not necessary but sometimes there are errors
		if created_since.After(c.Committer.When) {
			created_since = c.Committer.When
		}
		contributors[author]++
		orgs[org]++

		if flag {
			if created_since.After(parser.LAST_YEAR) {
				commit_count++
			} else {
				flag = false
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	metrics.CreatedSince = created_since
	metrics.ContributorCount = len(contributors)
	metrics.OrgCount = len(orgs)
	metrics.CommitFrequency = commit_count / 52

	return &metrics, nil
}

func ParseGitRepo(r *git.Repository) (*Repo, error) {

	repo := NewRepo()

	u, err := GetURL(r)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var name, owner, source string
	if u == "" {
		source = parser.UNKNOWN_SOURCE
		u = parser.UNKNOWN_URL
		name = parser.UNKNOWN_NAME
		owner = parser.UNKNOWN_OWNER
	} else {
		uu := url.ParseURL(u)
		path := strings.Split(uu.Pathname, "/")
		name = strings.Split(path[len(path)-1], ".")[0]
		owner = path[len(path)-2]
		if uu.Resource == "" {
			source = parser.UNKNOWN_SOURCE
		} else {
			source = uu.Resource
		}
	}

	license := ""
	//license = GetLicense(r)

	languages := make([]string, 0)
	/*
		languages := GetLanguages(r)
		if languages == nil {
			languages = &[]string{}
		}
	*/

	eco, err := GetEcosystem(r)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	metrics, err := GetMetrics(r)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if metrics == nil {
		m := NewRepoMetrics()
		metrics = &m
	}

	repo.Name = name
	repo.Owner = owner
	repo.Source = source
	repo.URL = u
	repo.License = license
	repo.Languages = languages
	repo.Ecosystems = eco
	repo.Metrics = metrics

	return &repo, nil
}
