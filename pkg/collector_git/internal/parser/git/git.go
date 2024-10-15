package git

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	parser "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"

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
	Languages []string
	Metrics   *RepoMetrics
}

func NewRepo() Repo {
	return Repo{
		Name:      "",
		Owner:     "",
		Source:    "",
		URL:       "",
		Languages: nil,
		Metrics:   &RepoMetrics{},
	}
}

func (r *Repo) Show() {
	utils.Info("Repo Name: %s", r.Name)
	utils.Info("Repo Owner: %s", r.Owner)
	utils.Info("Repo Source: %s", r.Source)
	utils.Info("Repo URL: %s", r.URL)
	utils.Info("Repo Languages: %s", r.Languages)
	r.Metrics.Show()
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

func (rm *RepoMetrics) Show() {
	utils.Info("Created Since: %s", rm.CreatedSince)
	utils.Info("Updated Since: %s", rm.UpdatedSince)
	utils.Info("Contributors Count: %d", rm.ContributorCount)
	utils.Info("Orgs Count: %d", rm.OrgCount)
	utils.Info("Commit Frequency: %f", rm.CommitFrequency)
}

func GetBlobs(r *git.Repository) *[]*object.Blob {
	bIter, err := r.BlobObjects()
	utils.CheckIfError(err)
	var blobs []*object.Blob
	err = bIter.ForEach(func(b *object.Blob) error {
		// fmt.Println(b)
		blobs = append(blobs, b)
		return nil
	})
	utils.CheckIfError(err)
	return &blobs
}

func GetBranches(r *git.Repository) *[]*plumbing.Reference {
	rIter, err := r.Branches()
	utils.CheckIfError(err)
	var refs []*plumbing.Reference
	err = rIter.ForEach(func(r *plumbing.Reference) error {
		// fmt.Println(r)
		refs = append(refs, r)
		return nil
	})
	utils.CheckIfError(err)
	return &refs
}

func GetCommits(r *git.Repository) *[]*object.Commit {
	cIter, err := r.CommitObjects()
	utils.CheckIfError(err)
	var commits []*object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		// fmt.Println(c)
		commits = append(commits, c)
		return nil
	})
	utils.CheckIfError(err)
	return &commits
}

func GetConfig(r *git.Repository) *gitconfig.Config {
	c, err := r.Config()
	utils.CheckIfError(err)
	return c
}

func GetObjects(r *git.Repository) *[]*object.Object {
	oIter, err := r.Objects()
	utils.CheckIfError(err)
	var objs []*object.Object
	err = oIter.ForEach(func(o object.Object) error {
		// fmt.Println(o)
		objs = append(objs, &o)
		return nil
	})
	utils.CheckIfError(err)
	return &objs
}

func GetReferences(r *git.Repository) *[]*plumbing.Reference {
	rIter, err := r.References()
	utils.CheckIfError(err)
	var refs []*plumbing.Reference
	err = rIter.ForEach(func(r *plumbing.Reference) error {
		// fmt.Println(r)
		refs = append(refs, r)
		return nil
	})
	utils.CheckIfError(err)
	return &refs
}

func GetRemotes(r *git.Repository) *[]*git.Remote {
	remotes, err := r.Remotes()
	utils.CheckIfError(err)
	return &remotes
}

func GetTags(r *git.Repository) *[]*object.Tag {
	tIter, err := r.TagObjects()
	utils.CheckIfError(err)
	var tags []*object.Tag
	err = tIter.ForEach(func(t *object.Tag) error {
		// fmt.Println(t)
		tags = append(tags, t)
		return nil
	})
	utils.CheckIfError(err)
	return &tags
}

func GetTagRefs(r *git.Repository) *[]*plumbing.Reference {
	rIter, err := r.Tags()
	utils.CheckIfError(err)
	var refs []*plumbing.Reference
	err = rIter.ForEach(func(r *plumbing.Reference) error {
		// fmt.Println(r)
		refs = append(refs, r)
		return nil
	})
	utils.CheckIfError(err)
	return &refs
}

func GetTrees(r *git.Repository) *[]*object.Tree {
	tIter, err := r.TreeObjects()
	utils.CheckIfError(err)
	var trees []*object.Tree
	err = tIter.ForEach(func(t *object.Tree) error {
		// fmt.Println(t)
		trees = append(trees, t)
		return nil
	})
	utils.CheckIfError(err)
	return &trees
}

func GetWorkTree(r *git.Repository) *git.Worktree {
	wt, err := r.Worktree()
	utils.CheckIfError(err)
	return wt
}

func GetURL(r *git.Repository) string {
	//? Git Fetch 和 Git Push 的 Remote URL 大部分情况下应该是相同的
	//? 但 Git Fetch 的 Remote URL 会被优先考虑作为这一仓库的 URL
	remotes := GetRemotes(r)
	if len(*remotes) == 0 {
		return ""
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
	return u
}

func GetLicense(r *git.Repository) string {

	// ToDO 尝试实现为 https://github.com/licensee/licensee

	//* Maybe multiple licenses
	//* l := []string{}
	l := parser.UNKNOWN_LICENSE
	ref, err := r.Head()
	utils.CheckIfError(err)
	commit, err := r.CommitObject(ref.Hash())
	utils.CheckIfError(err)
	tree, err := commit.Tree()
	utils.CheckIfError(err)

	f, err := tree.File("LICENSE")
	utils.CheckIfError(err)
	content, err := f.Lines()
	utils.CheckIfError(err)

	for k, v := range parser.LICENSE_KEYWORD {
		if strings.Contains(content[0], k) {
			//* l = append(l, v)
			l = v
		}
	}
	return l
}

func GetLanguages(r *git.Repository) *[]string {
	l := map[string]int{}
	ref, err := r.Head()
	utils.CheckIfError(err)
	commit, err := r.CommitObject(ref.Hash())
	utils.CheckIfError(err)
	tree, err := commit.Tree()
	utils.CheckIfError(err)
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
	utils.CheckIfError(err)
	var languages []string
	for k, v := range l {
		if v >= parser.LANGUAGE_THRESHOLD {
			languages = append(languages, k)
		}
	}
	return &languages
}

func GetMetrics(r *git.Repository) *RepoMetrics {
	//* 不直接使用 GetCommits 的目的在于 Log 可选择起止时间
	cIter, err := r.Log(&git.LogOptions{
		//* From:  ref.Hash(),
		All:   true,
		Since: &parser.BEGIN_TIME,
		Until: &parser.END_TIME,
		Order: git.LogOrderCommitterTime,
	})
	utils.CheckIfError(err)

	metrics := NewRepoMetrics()
	contributors := map[string]int{}
	orgs := map[string]int{}
	var commit_count float64 = 0

	latest_commit, err := cIter.Next()
	/*
		如果下载的仓库不完整就会导致迭代器为空的错误
			if err == io.EOF {
				return nil
			}
	*/
	utils.CheckIfError(err)

	author := fmt.Sprintf("%s(%s)", latest_commit.Author.Name, latest_commit.Author.Email)
	e := strings.Split(latest_commit.Author.Email, "@")
	org := e[len(e)-1]

	metrics.UpdatedSince = latest_commit.Author.When
	contributors[author]++
	orgs[org]++

	if latest_commit.Author.When.After(parser.LAST_YEAR) {
		commit_count++
	}

	flag := true
	created_since := time.Time{}

	err = cIter.ForEach(func(c *object.Commit) error {
		author := fmt.Sprintf("%s(%s)", c.Author.Name, c.Author.Email)
		e = strings.Split(c.Author.Email, "@")
		org := e[len(e)-1]

		created_since = c.Author.When
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
	utils.CheckIfError(err)

	metrics.CreatedSince = created_since
	metrics.ContributorCount = len(contributors)
	metrics.OrgCount = len(orgs)
	metrics.CommitFrequency = commit_count / 52

	return &metrics
}

func ParseGitRepo(r *git.Repository) *Repo {

	repo := NewRepo()

	u := GetURL(r)
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

	languages := &[]string{}
	/*
		languages := GetLanguages(r)
		if languages == nil {
			languages = &[]string{}
		}
	*/

	metrics := GetMetrics(r)
	if metrics == nil {
		m := NewRepoMetrics()
		metrics = &m
	}

	license := ""
	//license = GetLicense(r)

	repo.Name = name
	repo.Owner = owner
	repo.Source = source
	repo.URL = u
	repo.Languages = *languages
	repo.License = license
	repo.Metrics = metrics

	return &repo
}
