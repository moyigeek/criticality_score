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

const GitMetricTableName = "git_metrics"

type gitmetricsRepository struct {
	appDb storage.AppDatabaseContext
}

var _ GitMetricsRepository = (*gitmetricsRepository)(nil)

// BatchInsertOrUpdate implements GitMetricsRepository.
func (g *gitmetricsRepository) BatchInsertOrUpdate(data []*GitMetric) error {
	for _, d := range data {
		d.UpdateTime = sqlutil.ToNullable(time.Now())
	}
	return sqlutil.BatchInsert(g.appDb, string(GitMetricTableName), data)
}

// InsertOrUpdate implements GitMetricsRepository.
func (g *gitmetricsRepository) InsertOrUpdate(data *GitMetric) error {
	oldData, err := g.QueryByLink(*data.GitLink)
	if err != nil {
		sqlutil.MergeStruct(oldData, data)
	}
	data.UpdateTime = sqlutil.ToNullable(time.Now())
	return sqlutil.Insert(g.appDb, string(GitMetricTableName), data)
}

// Query implements GitMetricsRepository.
func (g *gitmetricsRepository) Query() (iter.Seq[*GitMetric], error) {
	subQuery := fmt.Sprintf(`(SELECT DISTINCT ON (git_link)
	 * 
	FROM %s
	ORDER BY git_link, id DESC)`, GitMetricTableName)
	return sqlutil.QueryCommon[GitMetric](g.appDb, subQuery, "")
}

// QueryByLink implements GitMetricsRepository.
func (g *gitmetricsRepository) QueryByLink(link string) (*GitMetric, error) {
	return sqlutil.QueryCommonFirst[GitMetric](g.appDb, GitMetricTableName, "WHERE git_link = $1 ORDER BY id DESC", link)
}

func NewGitMetricsRepository(appDb storage.AppDatabaseContext) GitMetricsRepository {
	return &gitmetricsRepository{appDb: appDb}
}
