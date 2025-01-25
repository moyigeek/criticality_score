package repository

import (
	"iter"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

type AllGitLinkRepository interface {
	/** QUERY **/
	Query() (iter.Seq[string], error)
}

type allGitLinkRepository struct {
	ctx storage.AppDatabaseContext
}

var _ AllGitLinkRepository = (*allGitLinkRepository)(nil)

// Query implements AllLinkRepository.
func (a *allGitLinkRepository) Query() (iter.Seq[string], error) {
	rows, err := a.ctx.Query("SELECT link FROM all_links")

	if err != nil {
		return nil, err
	}

	return func(yield func(string) bool) {
		for rows.Next() {
			var link string
			err := rows.Scan(&link)
			if err != nil {
				return
			}

			if !yield(link) {
				return
			}
		}
	}, nil
}
