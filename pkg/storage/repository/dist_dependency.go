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

	QueryByType(distType int) (iter.Seq[*DistLinkInfo], error)
	GetByLink(packageName string, distType int) (*DistLinkInfo, error)

	/** INSERT/UPDATE **/
	// update_time will be updated automatically
	InsertOrUpdate(packageInfo *DistLinkInfo) error
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

type DistLinkInfo struct {
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

// GetByLink implements DistributionDependencyRepository.
func (r *distLinkRepository) GetByLink(packageName string, distType int) (*DistLinkInfo, error) {
	return sqlutil.QueryCommonFirst[DistLinkInfo](r.ctx, DistDependencyTableName,
		`WHERE git_link = $1 and type = $2 ORDER BY id DESC`, packageName, distType)
}

// InsertOrUpdate implements DistributionDependencyRepository.
func (r *distLinkRepository) InsertOrUpdate(packageInfo *DistLinkInfo) error {
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
func (r *distLinkRepository) QueryByType(distType int) (iter.Seq[*DistLinkInfo], error) {
	return sqlutil.QueryCommon[DistLinkInfo](r.ctx, DistDependencyTableName,
		"where type = $1", distType)
}
