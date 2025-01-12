package repository

import (
	"iter"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/lib/pq"
)

type ResultRepository interface {
	/** QUERY **/
	QueryOrderByScoreWithPaging(take int, skip int) (iter.Seq[*Result], error)
	QueryByUntilOrderByScoreWithPaging(take int, skip int, until time.Time) (iter.Seq[*Result], error)

	QueryHistoryByLink(gitLink string) (iter.Seq[*Result], error)
	GetByLink(gitLink string) (*Result, error)
}

type Result struct {
	SeqId                  *int64          `column:"id"`
	GitLink                *string         `column:"git_link"`
	EcoSystem              *string         `column:"ecosystem"`
	CreatedSince           *time.Time      `column:"created_since"`
	UpdatedSince           *time.Time      `column:"updated_since"`
	ContributorCount       *int            `column:"contributor_count"`
	CommitFrequency        *float64        `column:"commit_frequency"`
	DepsDevCount           *int            `column:"depsdev_count"`
	DepsDistro             *string         `column:"deps_distro"`
	OrgCount               *int            `column:"org_count"`
	License                *string         `column:"license"`
	Language               *pq.StringArray `column:"language"`
	CloneValid             *bool           `column:"clone_valid"`
	DepsDevPageRank        *float64        `column:"depsdev_pagerank"`
	Scores                 *float64        `column:"scores"`
	IsDeleted              *bool           `column:"is_deleted"`
	UpdateTimeGitMetadata  *time.Time      `column:"update_time_git_metadata"`
	UpdateTimeDepsDev      *time.Time      `column:"update_time_deps_dev"`
	UpdateTimeDistribution *time.Time      `column:"update_time_distribution"`
	UpdateTimeScores       *time.Time      `column:"update_time_scores"`
	UpdateTime             *time.Time      `column:"update_time"`
}

type resultRepository struct {
	ctx storage.AppDatabaseContext
}

var _ ResultRepository = (*resultRepository)(nil)

// GetByLink implements ResultRepository.
func (r *resultRepository) GetByLink(gitLink string) (*Result, error) {
	panic("unimplemented")
}

// QueryByUntilOrderByScoreWithPaging implements ResultRepository.
func (r *resultRepository) QueryByUntilOrderByScoreWithPaging(take int, skip int, until time.Time) (iter.Seq[*Result], error) {
	panic("unimplemented")
}

// QueryHistoryByLink implements ResultRepository.
func (r *resultRepository) QueryHistoryByLink(gitLink string) (iter.Seq[*Result], error) {
	panic("unimplemented")
}

// QueryOrderByScoreWithPaging implements ResultRepository.
func (r *resultRepository) QueryOrderByScoreWithPaging(take int, skip int) (iter.Seq[*Result], error) {
	panic("unimplemented")
}

func NewResultRepository(appDb storage.AppDatabaseContext) ResultRepository {
	return &resultRepository{ctx: appDb}
}
