package repository

import "github.com/HUSTSecLab/criticality_score/pkg/storage"

type LangEcoLinkRepository interface {
	/** QUERY **/
	GetLangEcoInfoByLink(link string) ([]*LangEcoInfo, error)
	/** INSERT/UPDATE **/
	BatchInsertOrUpdateDistLinks(data []*LangEcoInfo) error
}

type langEcoLinkRepository struct {
	appDb storage.AppDatabaseContext
}

type LangEcoInfo struct {
	GitLink  *string
	NpmNum   *int
	GoNum    *int
	MavenNum *int
	PypiNum  *int
	NuGetNum *int
	CargoNum *int
}

func NewLangEcoLinkRepository(appDb storage.AppDatabaseContext) LangEcoLinkRepository {
	return &langEcoLinkRepository{
		appDb: appDb,
	}
}

func (r *langEcoLinkRepository) GetLangEcoInfoByLink(link string) ([]*LangEcoInfo, error) {
	panic("unimplemented")
}

func (r *langEcoLinkRepository) BatchInsertOrUpdateDistLinks(data []*LangEcoInfo) error {
	panic("unimplemented")
}
