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
	DepCount   *int
	PageRank   *float64
	UpdateTime *time.Time
}

func NewDistDependencyRepository(appDb storage.AppDatabaseContext) DistDependencyRepository {
	return &distLinkRepository{ctx: appDb}
}

// Query implements DistributionDependencyRepository.
func (r *distLinkRepository) Query() (iter.Seq[*DistDependency], error) {
	panic("unimplemented")
}

// QueryDistCountByType implements DistributionDependencyRepository.
func (r *distLinkRepository) QueryDistCountByType(distType int) (int, error) {
	panic("unimplemented")
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
