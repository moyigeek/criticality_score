package repository

import (
	"iter"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/sqlutil"
	"github.com/samber/lo"
)

const DistDependencyTableName = "distribution_dependencies"

type DistDependencyRepository interface {
	/** QUERY **/

	Query() (iter.Seq[*DistDependency], error) // Query all distribution information.
	QueryByType(distType int) (iter.Seq[*DistDependency], error)
	GetByLink(packageName string, distType int) (*DistDependency, error)
	QueryDistCountByType(distType int) (int, error) // Get the total number of packages in a Distro.

	/** INSERT/UPDATE **/
	// update_time will be updated automatically
	InsertOrUpdate(packageInfo *DistDependency) error
}

type distLinkRepository struct {
	ctx storage.AppDatabaseContext
}

var _ DistDependencyRepository = (*distLinkRepository)(nil)

type DistType int

const (
	Debian DistType = iota
	Arch
	Homebrew
	Nix
	Alpine
	Centos
	Aur
	Deepin
	Fedora
	Gentoo
	Ubuntu
)

type DistDependency struct {
	ID         *int64 `generated:"true"`
	GitLink    *string
	Type       *DistType
	DepImpact  *float64
	DepCount   *int
	PageRank   *float64
	UpdateTime *time.Time
}

func NewDistDependencyRepository(appDb storage.AppDatabaseContext) DistDependencyRepository {
	return &distLinkRepository{ctx: appDb}
}

// Query implements DistributionDependencyRepository.
func (r *distLinkRepository) Query() (iter.Seq[*DistDependency], error) {
	return sqlutil.Query[DistDependency](r.ctx, `SELECT DISTINCT ON (git_link)
		id, git_link, type, dep_impact, dep_count, page_rank, update_time
		FROM distribution_dependencies ORDER BY git_link, id DESC`)
}

// QueryDistCountByType implements DistributionDependencyRepository.
func (r *distLinkRepository) QueryDistCountByType(distType int) (int, error) {
	// FIXME: Delete situation is not considered in this qeury.
	row := r.ctx.QueryRow(`SELECT DISTINCT ON (git_link) SUM(dep_count) FROM distribution_dependencies WHERE git_link = $1 ORDER BY git_link, id DESC`, distType)
	var result int
	err := row.Scan(&result)
	return result, err
}

// GetByLink implements DistributionDependencyRepository.
func (r *distLinkRepository) GetByLink(packageName string, distType int) (*DistDependency, error) {
	return sqlutil.QueryCommonFirst[DistDependency](r.ctx, DistDependencyTableName,
		`WHERE git_link = $1 and type = $2 ORDER BY id DESC`, packageName, distType)
}

// InsertOrUpdate implements DistributionDependencyRepository.
func (r *distLinkRepository) InsertOrUpdate(packageInfo *DistDependency) error {
	if packageInfo.GitLink == nil || packageInfo.Type == nil {
		return ErrInvalidInput
	}

	packageInfo.UpdateTime = lo.ToPtr(time.Now())

	oldInfo, err := r.GetByLink(*packageInfo.GitLink, int(*packageInfo.Type))
	if err != nil {
		return err
	}

	sqlutil.MergeStruct(oldInfo, packageInfo)
	return sqlutil.Insert(r.ctx, DistDependencyTableName, packageInfo)
}

// QueryByType implements DistributionDependencyRepository.
func (r *distLinkRepository) QueryByType(distType int) (iter.Seq[*DistDependency], error) {
	return sqlutil.QueryCommon[DistDependency](r.ctx, DistDependencyTableName,
		"where type = $1", distType)
}
