package repositories

import (
	"fmt"
	"reflect"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type GitMetricsFrom int

type GitMetricsRepository interface {
	/** QUERY **/

	GetGitMetricsRankingLatest(take int, skip int) ([]*GitMetrics, error)
	GetGitMetricsRankingUntil(take int, skip int, until time.Time) ([]*GitMetrics, error)
	GetGitMetricsHistoryByLink(gitLink string) ([]*GitMetrics, error)

	GetGitMetricsByLinkIncludingDeleted(gitLink string) (*GitMetrics, error)
	GetGitMetricsByLink(gitLink string) (*GitMetrics, error)

	/** INSERT/UPDATE **/

	InsertOrUpdateGitMetrics(data *GitMetrics) error
	// this function will only insert the data, but not update
	BatchInsertGitMetrics(data []*GitMetrics) error

	/** DELETE **/

	DeleteGitMetricsByLink(gitLink string) error
	// this function will not examine whether the data exists
	BatchDeleteGitMetricsByLink(data []string) error
}

type gitmetricsRepository struct {
	appDb storage.AppDatabaseContext
}

// const (
// 	GitMetricsFromPacakge   GitMetricsFrom = 0
// 	GitMetricsFromGithub    GitMetricsFrom = 1
// 	GitMetricsFromGitlab    GitMetricsFrom = 2
// 	GitMetricsFromBitbucket GitMetricsFrom = 3
// )

type GitMetrics struct {
	SeqId            *int64          `column:"id"`
	GitLink          *string         `column:"git_link"`
	EcoSystem        *string         `column:"ecosystem"`
	CreatedSince     *time.Time      `column:"created_since"`
	UpdatedSince     *time.Time      `column:"updated_since"`
	ContributorCount *int            `column:"contributor_count"`
	CommitFrequency  *float64        `column:"commit_frequency"`
	DepsDevCount     *int            `column:"depsdev_count"`
	DepsDistro       *string         `column:"deps_distro"`
	OrgCount         *int            `column:"org_count"`
	License          *string         `column:"license"`
	Language         *pq.StringArray `column:"language"`
	CloneValid       *bool           `column:"clone_valid"`
	DepsDevPageRank  *float64        `column:"depsdev_pagerank"`
	Scores           *float64        `column:"scores"`
	IsDeleted        *bool           `column:"is_deleted"`

	// follwing is update by repository layer, any change by caller will be ignored
	UpdateTimeGitMetadata  *time.Time `column:"update_time_git_metadata"`
	UpdateTimeDepsDev      *time.Time `column:"update_time_deps_dev"`
	UpdateTimeDistribution *time.Time `column:"update_time_distribution"`
	UpdateTimeScores       *time.Time `column:"update_time_scores"`
	UpdateTime             *time.Time `column:"update_time"`
}

func NewGitMetricsRepository(appDb storage.AppDatabaseContext) GitMetricsRepository {
	return &gitmetricsRepository{appDb: appDb}
}

func (r *gitmetricsRepository) GetGitMetricsHistoryByLink(gitLink string) ([]*GitMetrics, error) {
	panic("unimplemented")
}

func (r *gitmetricsRepository) GetGitMetricsRankingLatest(take int, skip int) ([]*GitMetrics, error) {
	panic("unimplemented")
}

func (r *gitmetricsRepository) GetGitMetricsRankingUntil(take int, skip int, until time.Time) ([]*GitMetrics, error) {
	panic("unimplemented")
}

func (r *gitmetricsRepository) GetGitMetricsByLinkIncludingDeleted(gitLink string) (*GitMetrics, error) {
	rows, err := getDataFromTable[GitMetrics](r.appDb, "git_metrics_history", "WHERE git_link = $1 ORDER BY id DESC LIMIT 1", gitLink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rows.Next()
}

func (r *gitmetricsRepository) GetGitMetricsByLink(gitLink string) (*GitMetrics, error) {
	rows, err := getDataFromTable[GitMetrics](r.appDb, "git_metrics_history", "WHERE git_link = $1 and is_deleted = FALSE ORDER BY id DESC LIMIT 1", gitLink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rows.Next()
}

func (r *gitmetricsRepository) InsertOrUpdateGitMetrics(data *GitMetrics) error {
	if data.GitLink == nil {
		return fmt.Errorf("GitLink is required")
	}

	oldData, err := r.GetGitMetricsByLink(*data.GitLink)

	if err != nil {
		return err
	}

	if oldData != nil {
		// use reflection to update the map
		reflectType := reflect.TypeOf(*data)
		dataReflectVal := reflect.ValueOf(data).Elem()
		oldDataReflectVal := reflect.ValueOf(oldData).Elem()

		for i := 0; i < reflectType.NumField(); i++ {
			// field := reflectType.Field(i)
			if dataReflectVal.Field(i).IsNil() && !oldDataReflectVal.Field(i).IsNil() {
				val := oldDataReflectVal.Field(i).Interface()
				dataReflectVal.Field(i).Set(reflect.ValueOf(val))
			}
		}
	}

	if data.EcoSystem != nil || data.CreatedSince != nil || data.UpdatedSince != nil || data.ContributorCount != nil || data.CommitFrequency != nil || data.OrgCount != nil || data.License != nil || data.Language != nil || data.CloneValid != nil {
		t := time.Now()
		data.UpdateTimeGitMetadata = &t
	}

	t := time.Now()
	if data.DepsDistro != nil {
		data.UpdateTimeDistribution = &t
	}

	if data.DepsDevCount != nil || data.DepsDevPageRank != nil {
		data.UpdateTimeDepsDev = &t
	}

	isDeleted := false
	data.IsDeleted = &isDeleted
	data.UpdateTime = &t
	data.SeqId = nil

	return insertDataIntoTable(r.appDb, "git_metrics_history", data)
}

func (r *gitmetricsRepository) DeleteGitMetricsByLink(gitLink string) error {
	link, _ := r.GetGitMetricsByLink(gitLink)
	if link == nil {
		return fmt.Errorf("GitMetrics not found")
	}

	*link.IsDeleted = true
	link.SeqId = nil
	return insertDataIntoTable(r.appDb, "git_metrics_history", link)
}

func (r *gitmetricsRepository) BatchInsertGitMetrics(data []*GitMetrics) error {
	for _, d := range data {
		d.SeqId = nil
	}
	return batchInsertDataIntoTable(r.appDb, "git_metrics_history", data)
}

// NOTE: this function will not examine whether the data exists
func (r *gitmetricsRepository) BatchDeleteGitMetricsByLink(data []string) error {
	toDelete := make([]*GitMetrics, 0)
	for _, d := range data {
		toDelete = append(toDelete, &GitMetrics{GitLink: &d, IsDeleted: lo.ToPtr(true)})
	}
	return r.BatchInsertGitMetrics(toDelete)
}
