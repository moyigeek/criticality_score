package repository

import "github.com/HUSTSecLab/criticality_score/pkg/storage"

type DistLinkRepository interface {
	/** QUERY **/
	GetDistInfoByLink(packageName string) ([]*DistLinkInfo, error)
	/** INSERT/UPDATE **/
	InsertOrUpdateDistPackage(packageInfo *DistLinkInfo) error
}

type distLinkRepository struct {
	appDb storage.AppDatabaseContext
}

type DistLinkInfo struct {
	GitLink          *string
	AlpineDepCount   *int
	AlpinePageRank   *float64
	ArchDepCount     *int
	ArchPageRank     *float64
	AurDepCount      *int
	AurPageRank      *float64
	CentosDepCount   *int
	CentosPageRank   *float64
	DebianDepCount   *int
	DebianPageRank   *float64
	DeepinDepCount   *int
	DeepinPageRank   *float64
	FedoraDepCount   *int
	FedoraPageRank   *float64
	GentooDepCount   *int
	GentooPageRank   *float64
	HomebrewDepCount *int
	HomebrewPageRank *float64
	NixDepCount      *int
	NixPageRank      *float64
	UbuntuDepCount   *int
	UbuntuPageRank   *float64
}

func NewDistLinkRepository(appDb storage.AppDatabaseContext) DistLinkRepository {
	return &distLinkRepository{appDb: appDb}
}

func (r *distLinkRepository) GetDistInfoByLink(packageName string) ([]*DistLinkInfo, error) {
	panic("unimplemented")
}

func (r *distLinkRepository) InsertOrUpdateDistPackage(packageInfo *DistLinkInfo) error {
	panic("unimplemented")
}
