package repository

import (
	"fmt"
	"iter"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/sqlutil"
	"github.com/lib/pq"
)

type GitMetricsRepository interface {
	/** QUERY **/
	Query() (iter.Seq[*GitMetric], error)
	QueryByLink(link string) (*GitMetric, error)

	/** INSERT/UPDATE **/
	// NOTE: update_time will be updated automatically
	InsertOrUpdate(data *GitMetric) error
	// NOTE: update_time will be updated automatically
	// and the data will not copy from old data
	BatchInsertOrUpdate(data []*GitMetric) error

	// times will be updated automatically
	InsertOrUpdateFailed(data *FailedGitMetric) error
	DeleteFailed(link string) error
}

type GitMetric struct {
	ID               *int64 `generated:"true"`
	GitLink          *string
	CreatedSince     **time.Time
	UpdatedSince     **time.Time
	ContributorCount **int
	CommitFrequency  **float64
	OrgCount         **int
	License          **pq.StringArray
	Language         **pq.StringArray
	CloneValid       **bool
	UpdateTime       **time.Time
}

type FailedGitMetric struct {
	GitLink    *string `pk:"true"`
	Message    **string
	UpdateTime **time.Time
	Times      **int
}

const GitMetricTableName = "git_metrics"
const FailedGitMetricTableName = "failed_git_metrics"

type gitmetricsRepository struct {
	ctx storage.AppDatabaseContext
}

var _ GitMetricsRepository = (*gitmetricsRepository)(nil)

// BatchInsertOrUpdate implements GitMetricsRepository.
func (g *gitmetricsRepository) BatchInsertOrUpdate(data []*GitMetric) error {
	for _, d := range data {
		d.UpdateTime = sqlutil.ToNullable(time.Now())
	}
	return sqlutil.BatchInsert(g.ctx, string(GitMetricTableName), data)
}

// InsertOrUpdate implements GitMetricsRepository.
func (g *gitmetricsRepository) InsertOrUpdate(data *GitMetric) error {
	oldData, err := g.QueryByLink(*data.GitLink)
	if err != nil {
		sqlutil.MergeStruct(oldData, data)
	}
	data.UpdateTime = sqlutil.ToNullable(time.Now())
	return sqlutil.Insert(g.ctx, string(GitMetricTableName), data)
}

// Query implements GitMetricsRepository.
func (g *gitmetricsRepository) Query() (iter.Seq[*GitMetric], error) {
	subQuery := fmt.Sprintf(`(SELECT DISTINCT ON (git_link)
	 * 
	FROM %s
	ORDER BY git_link, id DESC)`, GitMetricTableName)
	return sqlutil.QueryCommon[GitMetric](g.ctx, subQuery, "")
}

// QueryByLink implements GitMetricsRepository.
func (g *gitmetricsRepository) QueryByLink(link string) (*GitMetric, error) {
	return sqlutil.QueryCommonFirst[GitMetric](g.ctx, GitMetricTableName, "WHERE git_link = $1 ORDER BY id DESC", link)
}

// DeleteFailed implements GitMetricsRepository.
func (g *gitmetricsRepository) DeleteFailed(link string) error {
	return sqlutil.Delete(g.ctx, FailedGitMetricTableName, &FailedGitMetric{GitLink: &link})
}

// InsertOrUpdateFailed implements GitMetricsRepository.
func (g *gitmetricsRepository) InsertOrUpdateFailed(data *FailedGitMetric) error {
	_, err := g.ctx.Exec(`INSERT INTO `+FailedGitMetricTableName+` (git_link, message, update_time, times)
		VALUES ($1, $2, $3, 1)
		ON CONFLICT (git_link) DO UPDATE SET message = $2, update_time = $3, times = times + 1`,
		data.GitLink, data.Message, data.UpdateTime)
	return err

}

func NewGitMetricsRepository(appDb storage.AppDatabaseContext) GitMetricsRepository {
	return &gitmetricsRepository{ctx: appDb}
}
