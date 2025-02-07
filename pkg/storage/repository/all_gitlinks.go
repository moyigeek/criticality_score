package repository

import (
	"iter"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

type AllGitLinkRepository interface {
	/** QUERY **/
	Query() (iter.Seq[string], error)
	QueryCache() (iter.Seq[string], error)
	MakeCache() error
}

type allGitLinkRepository struct {
	ctx storage.AppDatabaseContext
}

// MakeCache implements AllGitLinkRepository.
func (a *allGitLinkRepository) MakeCache() error {
	_, err := a.ctx.Exec(`DROP TABLE IF EXISTS all_gitlinks_cache;
	CREATE TABLE all_gitlinks_cache AS SELECT * FROM all_gitlinks;
	`)
	return err
}

// QueryCache implements AllGitLinkRepository.
func (a *allGitLinkRepository) QueryCache() (iter.Seq[string], error) {
	rows, err := a.ctx.Query("SELECT git_link FROM all_gitlinks_cache")

	if err != nil {
		return nil, err
	}

	return func(yield func(string) bool) {
		defer rows.Close()
		for rows.Next() {
			var link *string
			err := rows.Scan(&link)
			if err != nil {
				return
			}
			if link == nil {
				continue
			}

			if !yield(*link) {
				return
			}
		}
	}, nil
}

var _ AllGitLinkRepository = (*allGitLinkRepository)(nil)

func NewAllGitLinkRepository(appDb storage.AppDatabaseContext) AllGitLinkRepository {
	return &allGitLinkRepository{ctx: appDb}
}

// Query implements AllLinkRepository.
func (a *allGitLinkRepository) Query() (iter.Seq[string], error) {
	rows, err := a.ctx.Query("SELECT git_link FROM all_gitlinks")

	if err != nil {
		return nil, err
	}

	return func(yield func(string) bool) {
		defer rows.Close()
		for rows.Next() {
			var link *string
			err := rows.Scan(&link)
			if err != nil {
				return
			}
			if link == nil {
				continue
			}

			if !yield(*link) {
				return
			}
		}
	}, nil
}
