package score

import (
	"log"
	"math"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/repository"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/sqlutil"
)

type LinkScore struct {
	GitMetrics       []*repository.GitMetric
	GitMetadataScore GitMetadataScore
	LangEcosystems   []*repository.LangEcosystem
	LangEcoScore     LangEcoScore
	DistDependencies []*repository.DistDependency
	DistScore        DistScore
	Score            float64
}

type GitMetadata struct {
	Id               int64
	CreatedSince     time.Time
	UpdatedSince     time.Time
	ContributorCount int
	CommitFrequency  float64
	Org_Count        int
}

type GitMetadataScore struct {
	GitMetrics       []*repository.GitMetric
	GitMetadataScore float64
}

type DistMetadata struct {
	Id        int64
	DepImpact float64
	DepCount  int
	PageRank  float64
	Type      repository.DistType
}

type LangEcoMetadata struct {
	Id            int64
	Type          repository.LangEcosystemType
	LangEcoImpact float64
	DepCount      int
}

type DistScore struct {
	DistDependencies []*repository.DistDependency
	DistImpact       float64
	DistPageRank     float64
	DistScore        float64
}

type LangEcoScore struct {
	LangEcosystems  []*repository.LangEcosystem
	LangEcoImpact   float64
	LangEcoPageRank float64
	LangEcoScore    float64
}

// Define weights (αi) and max thresholds (Ti)
var weights = map[string]map[string]float64{
	"gitMetadataScore": {
		"created_since":     1,
		"updated_since":     -1,
		"contributor_count": 2,
		"commit_frequency":  1,
		"org_count":         1,
		"gitMetadataScore":  1,
	},
	"distScore": {
		"dist_impact":   1,
		"dist_pagerank": 1,
		"distScore":     5,
	},
	"langEcoScore": {
		"lang_eco_impact": 1,
		"langEcoScore":    5,
	},
}

var thresholds = map[string]map[string]float64{
	"gitMetadataScore": {
		"created_since":     120,
		"updated_since":     120,
		"contributor_count": 40000,
		"commit_frequency":  1000,
		"org_count":         8400,
		"gitMetadataScore":  1,
	},
	"distScore": {
		"dist_impact":   1,
		"dist_pagerank": 1,
		"distScore":     1,
	},
	"langEcoScore": {
		"lang_eco_impact": 1,
		"langEcoScore":    1,
	},
}

var PackageList = map[repository.DistType]int{
	repository.Debian:   0,
	repository.Arch:     0,
	repository.Nix:      0,
	repository.Homebrew: 0,
	repository.Gentoo:   0,
	repository.Alpine:   0,
	repository.Fedora:   0,
	repository.Ubuntu:   0,
	repository.Deepin:   0,
	repository.Aur:      0,
	repository.Centos:   0,
}

var PackageCounts = map[repository.LangEcosystemType]int{
	repository.Npm:   3.37e6,
	repository.Go:    1.29e6,
	repository.Maven: 668e3,
	repository.Pypi:  574e3,
	repository.NuGet: 430e3,
	repository.Cargo: 168e3,
}

func (langEcoMetadata *LangEcoMetadata) ParseLangEcoMetadata(langEcosystem *repository.LangEcosystem) {
	langEcoMetadata.Id = *langEcosystem.ID
	langEcoMetadata.Type = *langEcosystem.Type
	langEcoMetadata.DepCount = *langEcosystem.DepCount
	langEcoMetadata.LangEcoImpact = *langEcosystem.LangEcoImpact
}

func (distMetadata *DistMetadata) PraseDistMetadata(distLink *repository.DistDependency) {
	distMetadata.Id = *distLink.ID
	distMetadata.DepCount = *distLink.DepCount
	distMetadata.PageRank = *distLink.PageRank
	distMetadata.Type = *distLink.Type
}

func (gitMetadata *GitMetadata) ParseMetadata(gitMetic *repository.GitMetric) {
	gitMetadata.Id = *gitMetic.ID
	if !sqlutil.IsNull(gitMetic.CreatedSince) {
		gitMetadata.CreatedSince = **gitMetic.CreatedSince
	}
	if !sqlutil.IsNull(gitMetic.UpdatedSince) {
		gitMetadata.UpdatedSince = **gitMetic.UpdatedSince
	}
	if !sqlutil.IsNull(gitMetic.ContributorCount) {
		gitMetadata.ContributorCount = **gitMetic.ContributorCount
	}
	if !sqlutil.IsNull(gitMetic.CommitFrequency) {
		gitMetadata.CommitFrequency = **gitMetic.CommitFrequency
	}
	if !sqlutil.IsNull(gitMetic.OrgCount) {
		gitMetadata.Org_Count = **gitMetic.OrgCount
	}
}

func (langEcoScore *LangEcoScore) CalculateLangEcoScore() {
	langEcoScore.LangEcoScore = weights["lang_eco_score"]["lang_eco_impact"] * langEcoScore.LangEcoImpact
}

func NewLangEcoScore() *LangEcoScore {
	return &LangEcoScore{}
}

func (gitMetadataScore *GitMetadataScore) CalculateGitMetadataScore(gitMetadata *GitMetadata) {
	var score float64
	var createdSinceScore, updatedSinceScore, contributorCountScore, commitFrequencyScore, orgCountScore float64

	monthsSinceCreation := time.Since(gitMetadata.CreatedSince).Hours() / (24 * 30)
	normalized := math.Log(monthsSinceCreation+1) / math.Log(math.Max(monthsSinceCreation, thresholds["gitMetadataScore"]["created_since"])+1)
	createdSinceScore = weights["gitMetadataScore"]["created_since"] * normalized
	score += createdSinceScore

	monthsSinceUpdate := time.Since(gitMetadata.UpdatedSince).Hours() / (24 * 30)
	normalized = math.Log(monthsSinceUpdate+1) / math.Log(math.Max(monthsSinceUpdate, thresholds["gitMetadataScore"]["updated_since"])+1)
	updatedSinceScore = weights["gitMetadataScore"]["updated_since"] * normalized
	score += updatedSinceScore

	normalized = math.Log(float64(gitMetadata.ContributorCount)+1) / math.Log(math.Max(float64(gitMetadata.ContributorCount), thresholds["gitMetadataScore"]["contributor_count"])+1)
	contributorCountScore = weights["gitMetadataScore"]["contributor_count"] * normalized
	score += contributorCountScore

	normalized = math.Log(gitMetadata.CommitFrequency+1) / math.Log(math.Max(gitMetadata.CommitFrequency, thresholds["gitMetadataScore"]["commit_frequency"])+1)
	commitFrequencyScore = weights["gitMetadataScore"]["commit_frequency"] * normalized
	score += commitFrequencyScore

	normalized = math.Log(float64(gitMetadata.Org_Count)+1) / math.Log(math.Max(float64(gitMetadata.Org_Count), thresholds["gitMetadataScore"]["org_count"])+1)
	orgCountScore = weights["gitMetadataScore"]["org_count"] * normalized
	score += orgCountScore

	gitMetadataScore.GitMetadataScore = score
	gitMetadataScore.GitMetrics = []*repository.GitMetric{
		{
			CreatedSince:     sqlutil.ToNullable(gitMetadata.CreatedSince),
			UpdatedSince:     sqlutil.ToNullable(gitMetadata.UpdatedSince),
			ContributorCount: sqlutil.ToNullable(gitMetadata.ContributorCount),
			CommitFrequency:  sqlutil.ToNullable(gitMetadata.CommitFrequency),
			OrgCount:         sqlutil.ToNullable(gitMetadata.Org_Count),
		},
	}
}

func NewGitMetadata() *GitMetadata {
	return &GitMetadata{}
}

func (distScore *DistScore) CalculateDistScore() {
	distScore.DistScore = weights["distScore"]["dist_impact"]*distScore.DistImpact + weights["distScore"]["dist_pagerank"]*distScore.DistPageRank
}
func (linkScore *LinkScore) CalculateScore() {
	score := 0.0

	score += weights["gitMetadataScore"]["gitMetadataScore"] * linkScore.GitMetadataScore.GitMetadataScore

	score += weights["lang_eco_impact"]["lang_eco_impact"] * linkScore.LangEcoScore.LangEcoScore

	score += weights["distScore"]["distScore"] * linkScore.DistScore.DistScore

	var totalnum float64
	for nameScore, value := range weights {
		for nameSubScore := range value {
			if nameSubScore != nameScore {
				totalnum += weights["gitMetadataScore"][nameSubScore]
			}
		}
	}
	linkScore.Score = score / totalnum
}

func NewGitMetadataScore() *GitMetadataScore {
	return &GitMetadataScore{}
}

func NewLangEcoMetadata() *LangEcoMetadata {
	return &LangEcoMetadata{}
}

func NewDistMetadata() *DistMetadata {
	return &DistMetadata{}
}

func NewLinkScore(gitMetadataScore *GitMetadataScore, distScore *DistScore, langEcoScore *LangEcoScore) *LinkScore {
	return &LinkScore{
		LangEcoScore:     *langEcoScore,
		DistScore:        *distScore,
		GitMetadataScore: *gitMetadataScore,
	}
}

func NewDistScore() *DistScore {
	return &DistScore{}
}

func LogNormalize(value, threshold float64) float64 {
	return math.Log(value+1) / math.Log(math.Max(value, threshold)+1)
}

func FetchGitMetrics(ac storage.AppDatabaseContext) map[string]*GitMetadata {
	repo := repository.NewGitMetricsRepository(ac)
	linksIter, err := repo.Query()
	linksMap := make(map[string]*GitMetadata)
	if err != nil {
		log.Fatalf("Failed to fetch git links: %v", err)
	}
	for link := range linksIter {
		gitMetadata := NewGitMetadata()
		gitMetadata.ParseMetadata(link)
		linksMap[*link.GitLink] = gitMetadata
	}
	return linksMap
}

func FetchLangEcoMetadata(ac storage.AppDatabaseContext) map[string]*LangEcoScore {
	repo := repository.NewLangEcoLinkRepository(ac)
	LangEcoMap := make(map[string]*LangEcoScore)
	linksIter, err := repo.Query()
	if err != nil {
		log.Fatalf("Failed to fetch lang eco links: %v", err)
	}
	for link := range linksIter {
		langEcoMetadata := NewLangEcoMetadata()
		langEcoMetadata.ParseLangEcoMetadata(link)
		if exists, ok := LangEcoMap[*link.GitLink]; ok && exists != nil {
			LangEcoMap[*link.GitLink].LangEcosystems = append(LangEcoMap[*link.GitLink].LangEcosystems, link)
			LangEcoMap[*link.GitLink].LangEcoImpact += langEcoMetadata.LangEcoImpact
		} else {
			LangEcoMap[*link.GitLink] = &LangEcoScore{LangEcosystems: []*repository.LangEcosystem{link}, LangEcoImpact: langEcoMetadata.LangEcoImpact}
		}
	}
	return LangEcoMap
}

func FetchDistMetadata(ac storage.AppDatabaseContext) map[string]*DistScore {
	repo := repository.NewDistDependencyRepository(ac)
	distMap := make(map[string]*DistScore)
	linksIter, err := repo.Query()
	if err != nil {
		log.Fatalf("Failed to fetch dist links: %v", err)
	}
	for link := range linksIter {
		distMetadata := NewDistMetadata()
		distMetadata.PraseDistMetadata(link)
		if exists, ok := distMap[*link.GitLink]; ok && exists != nil {
			distMap[*link.GitLink].DistDependencies = append(distMap[*link.GitLink].DistDependencies, link)
			distMap[*link.GitLink].DistImpact += distMetadata.DepImpact
			distMap[*link.GitLink].DistPageRank += distMetadata.PageRank
		} else {
			distMap[*link.GitLink] = &DistScore{DistDependencies: []*repository.DistDependency{link}, DistImpact: distMetadata.DepImpact, DistPageRank: distMetadata.PageRank}
		}
	}
	return distMap
}
func FetchGitLink(ac storage.AppDatabaseContext) []string {
	repo := repository.NewAllGitLinkRepository(ac)
	linksIter, err := repo.Query()
	if err != nil {
		log.Fatalf("Failed to fetch git links: %v", err)
	}
	links := []string{}
	for link := range linksIter {
		links = append(links, link)
	}
	return links
}

func UpdatePackageList(ac storage.AppDatabaseContext) {
	repo := repository.NewDistDependencyRepository(ac)
	for distType := range PackageList {
		count, err := repo.QueryDistCountByType(distType)
		if err != nil {
			log.Fatalf("Failed to fetch dist links: %v", err)
		}
		PackageList[distType] = count
	}
}

func UpdateScore(ac storage.AppDatabaseContext, packageScore map[string]*LinkScore) {
	repo := repository.NewScoreRepository(ac)
	scores := []*repository.Score{}
	for link, linkScore := range packageScore {
		score := repository.Score{
			Score:            &linkScore.Score,
			GitLink:          &link,
			DistDependencies: linkScore.DistDependencies,
			GitMetrics:       linkScore.GitMetrics,
			LangEcosystems:   linkScore.LangEcosystems,
			DistScore:        &linkScore.DistScore.DistScore,
			LangScore:        &linkScore.LangEcoScore.LangEcoScore,
			GitScore:         &linkScore.GitMetadataScore.GitMetadataScore,
		}
		scores = append(scores, &score)
	}
	if err := repo.BatchInsertOrUpdate(scores); err != nil {
		log.Fatalf("Failed to update score: %v", err)
	}
}

func FetchDistMetadataSingle(ac storage.AppDatabaseContext, link string) map[string]*DistScore {
	repo := repository.NewDistDependencyRepository(ac)
	linksMap := []*repository.DistDependency{}
	distMap := make(map[string]*DistScore)
	for PackageType := range PackageList {
		distInfo, err := repo.GetByLink(link, int(PackageType))
		if err != nil {
			log.Fatalf("Failed to fetch dist links: %v", err)
		}
		linksMap = append(linksMap, distInfo)
	}
	for _, link := range linksMap {
		distMetadata := NewDistMetadata()
		distMetadata.PraseDistMetadata(link)
		if exists, ok := distMap[*link.GitLink]; ok && exists != nil {
			distMap[*link.GitLink].DistDependencies = append(distMap[*link.GitLink].DistDependencies, link)
			distMap[*link.GitLink].DistImpact += distMetadata.DepImpact
			distMap[*link.GitLink].DistPageRank += distMetadata.PageRank
		} else {
			distMap[*link.GitLink] = &DistScore{DistDependencies: []*repository.DistDependency{link}, DistImpact: distMetadata.DepImpact, DistPageRank: distMetadata.PageRank}
		}
	}
	return distMap
}

func FetchLangEcoMetadataSingle(ac storage.AppDatabaseContext, link string) map[string]*LangEcoScore {
	repo := repository.NewLangEcoLinkRepository(ac)
	langEcoMap := make(map[string]*LangEcoScore)
	linksIter, err := repo.QueryByLink(link)
	if err != nil {
		log.Fatalf("Failed to fetch lang eco links: %v", err)
	}
	for link := range linksIter {
		langEcoMetadata := NewLangEcoMetadata()
		langEcoMetadata.ParseLangEcoMetadata(link)
		if exists, ok := langEcoMap[*link.GitLink]; ok && exists != nil {
			langEcoMap[*link.GitLink].LangEcosystems = append(langEcoMap[*link.GitLink].LangEcosystems, link)
			langEcoMap[*link.GitLink].LangEcoImpact += langEcoMetadata.LangEcoImpact
		} else {
			langEcoMap[*link.GitLink] = &LangEcoScore{LangEcosystems: []*repository.LangEcosystem{link}, LangEcoImpact: langEcoMetadata.LangEcoImpact}
		}
	}
	return langEcoMap
}
func FetchGitMetricsSingle(ac storage.AppDatabaseContext, link string) map[string]*GitMetadata {
	repo := repository.NewGitMetricsRepository(ac)
	linkInfo, err := repo.QueryByLink(link)
	linksMap := make(map[string]*GitMetadata)
	if err != nil {
		log.Fatalf("Failed to fetch git links: %v", err)
	}
	gitMetadata := NewGitMetadata()
	gitMetadata.ParseMetadata(linkInfo)
	linksMap[*linkInfo.GitLink] = gitMetadata
	return linksMap
}
