package repository

import (
	"iter"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/sqlutil"
	"github.com/samber/lo"
)

type ScoreRepository interface {
	/** QUERY **/

	Query() (iter.Seq[*Score], error)
	GetByGitLink(distID int64) (*Score, error)

	/** INSERT/UPDATE **/

	// When inserting, make sure all score is properly calculated
	// Any data will not be copied from old data
	// update_time will be updated automatically
	InsertOrUpdate(score *Score) error
	// When inserting, make sure all score is properly calculated
	// Any data will not be copied from old data
	// update_time will be updated automatically
	BatchInsertOrUpdate(scores []*Score) error
}

type Score struct {
	ID         *int64 `pk:"true" generated:"true"`
	GitLink    *string
	DistID     *int64
	DistScore  *float64
	DepsDevID  *int64
	DevScore   *float64
	GitID      *int64
	GitScore   *float64
	Score      *float64
	UpdateTime *time.Time
}

const ScoreTableName = "scores"

var _ ScoreRepository = (*scoreRepository)(nil)

type scoreRepository struct {
	appDb storage.AppDatabaseContext
}

// BatchInsertOrUpdate implements ScoreRepository.
func (s *scoreRepository) BatchInsertOrUpdate(scores []*Score) error {
	for _, score := range scores {
		score.UpdateTime = lo.ToPtr(time.Now())
	}

	return sqlutil.BatchInsert(s.appDb, ScoreTableName, scores)
}

// BatchIntLink implements ScoreRepository.
func (s *scoreRepository) GetByGitLink(distID int64) (*Score, error) {
	return sqlutil.QueryCommonFirst[Score](s.appDb, ScoreTableName,
		`WHERE git_link = $1 ORDER BY DESC ID`,
		distID)
}

// InsertOrUpdate implements ScoreRepository.
func (s *scoreRepository) InsertOrUpdate(score *Score) error {
	score.UpdateTime = lo.ToPtr(time.Now())

	return sqlutil.Insert(s.appDb, ScoreTableName, score)
}

// Query implements ScoreRepository.
func (s *scoreRepository) Query() (iter.Seq[*Score], error) {
	subQuery := `(SELECT DISTINCT ON (git_link) * FROM ` + ScoreTableName + ` ORDER BY git_link, id DESC)`
	return sqlutil.QueryCommon[Score](s.appDb, subQuery, "")
}

func NewScoreRepository(appDb storage.AppDatabaseContext) ScoreRepository {
	return &scoreRepository{
		appDb: appDb,
	}
}
