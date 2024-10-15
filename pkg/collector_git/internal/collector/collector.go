package collector

import (
	config "github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	parser "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func Collect(u *url.RepoURL) (*gogit.Repository, error) {
	r, err := Clone(u)
	if err == gogit.ErrRepositoryAlreadyExists {
		r, err = Update(u)
	}
	return r, err
}

func BriefCollect(u *url.RepoURL) (*gogit.Repository, error) {
	r, err := BriefClone(u)
	/*
		if err == gogit.ErrRepositoryAlreadyExists {
			r, err = Update(u)
		}
	*/
	return r, err
}

func EzCollect(u *url.RepoURL) (*gogit.Repository, error) {
	r, err := MemClone(u)
	return r, err
}

func Clone(u *url.RepoURL) (*gogit.Repository, error) {
	path := config.STORAGE_PATH + u.Pathname
	r, err := gogit.PlainClone(path, false, &gogit.CloneOptions{
		URL: u.URL,
		// Progress:     os.Stdout,
		SingleBranch: false,
	})
	return r, err
}

func BriefClone(u *url.RepoURL) (*gogit.Repository, error) {
	path := config.STORAGE_PATH + u.Pathname
	r, err := gogit.PlainClone(path, true, &gogit.CloneOptions{
		URL: u.URL,
		// Progress:     os.Stdout,
		SingleBranch: false,
	})
	return r, err
}

func MemClone(u *url.RepoURL) (*gogit.Repository, error) {
	r, err := gogit.Clone(memory.NewStorage(), nil, &gogit.CloneOptions{
		URL: u.URL,
		// Progress:     os.Stdout,
		SingleBranch: false,
	})
	return r, err
}

func Open(path string) (*gogit.Repository, error) {
	r, err := gogit.PlainOpen(path)
	return r, err
}

func Pull(r *gogit.Repository, path string) error {
	wt, err := r.Worktree()
	utils.CheckIfError(err)
	remotes, err := r.Remotes()
	utils.CheckIfError(err)
	var remote, u string

	if len(remotes) > 0 {
		remote = (remotes)[0].Config().Name
		urls := (remotes)[0].Config().URLs
		if len(urls) > 0 {
			u = urls[0]
		}
	}

	if remote == "" {
		remote = parser.DEFAULT_REMOTE_NAME
	}

	if u == "" {
		u = "https://" + parser.DEFAULT_SOURCE + path
	}

	err = wt.Pull(&gogit.PullOptions{
		RemoteName:   remote,
		RemoteURL:    u,
		SingleBranch: true,
		Force:        true,
	})

	return err
}

/*
func Fetch(r *gogit.Repository, path string) error {
	var u string

	remotes := git.GetRemotes(r)
	if len(*remotes) > 0 {
		us := (*remotes)[0].Config().URLs
		if len(us) > 0 {
			u = us[0]
		}
	}

	if u == "" {
		u = "https://" + parser.DEFAULT_SOURCE + path
	}

	err := r.Fetch(&gogit.FetchOptions{
		RemoteURL: u,
		RefSpecs:  []gogitconfig.RefSpec{"refs/*:refs/*", "HEAD:ref/heads/HEAD"},
		// Progress:  os.Stdout,
	})
	return err
}
*/

func Update(u *url.RepoURL) (*gogit.Repository, error) {
	r, err := Open(config.STORAGE_PATH + u.Pathname)
	if err != nil {
		return r, err
	}
	err = Pull(r, u.Pathname)
	// err := fetch(r)
	if err == gogit.NoErrAlreadyUpToDate {
		err = nil
	}
	return r, err
}
