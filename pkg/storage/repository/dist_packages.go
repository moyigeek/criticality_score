package repository

import "github.com/HUSTSecLab/criticality_score/pkg/storage"

type DistPackageRepository interface {
	/** QUERY **/
	GetDistInfoByPackage(packageName string) ([]*PackageInfo, error)
	/** INSERT/UPDATE **/
	InsertOrUpdateDistPackage(packageInfo *PackageInfo) error
}

type DistLinkTablePrefix string

const (
	DistLinkTablePrefixAlpine    DistLinkTablePrefix = "alpine"
	DistLinkTablePrefixArchlinux                     = "arch"
	DistLinkTablePrefixAur                           = "aur"
	DistLinkTablePrefixCentos                        = "centos"
	DistLinkTablePrefixDebian                        = "debian"
	DistLinkTablePrefixDeepin                        = "deepin"
	DistLinkTablePrefixFedora                        = "fedora"
	DistLinkTablePrefixGentoo                        = "gentoo"
	DistLinkTablePrefixHomebrew                      = "homebrew"
	DistLinkTablePrefixNix                           = "nix"
	DistLinkTablePrefixUbuntu                        = "ubuntu"
)

type distPackageRepository struct {
	AppDb storage.AppDatabaseContext
	Dist  DistLinkTablePrefix
}

type PackageInfo struct {
	Package      *string
	HomePage     *string
	Description  *string
	DependsCount *int
	PageRank     *float64
	Version      *string
}

func NewDistPackageRepository(appDb storage.AppDatabaseContext, dist DistLinkTablePrefix) DistPackageRepository {
	return &distPackageRepository{
		AppDb: appDb,
		Dist:  dist,
	}
}

func (r *distPackageRepository) InsertOrUpdateDistPackage(packageInfo *PackageInfo) error {
	panic("unimplemented")
}

func (r *distPackageRepository) GetDistInfoByPackage(packageName string) ([]*PackageInfo, error) {
	panic("unimplemented")
}
