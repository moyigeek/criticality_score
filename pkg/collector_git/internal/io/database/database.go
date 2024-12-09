/*
 * @Date: 2024-09-07 20:46:04
 * @LastEditTime: 2024-12-09 19:10:51
 * @Description:
 */
package database

import (
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
)

const (
	SQLITE_MAX_IDLE_CONNS int = 10
	SQLITE_MAX_OPEN_CONNS int = 100
	BATCH_SIZE                = 500
)

type GitMetrics struct {
	// gorm.Model
	Name             string    `gorm:"column:_name;not null"`
	Owner            string    `gorm:"column:_owner;not null"`
	Source           string    `gorm:"column:_source;not null"`
	URL              string    `gorm:"column:git_link;not null"` //* `gorm:"unique;not null"`
	License          string    `gorm:"column:license;not null"`
	Ecosystems       string    `gorm:"column:ecosystem;not null"`
	Languages        string    `gorm:"column:languages;not null"`
	CreatedSince     time.Time `gorm:"column:created_since;not null"`
	UpdatedSince     time.Time `gorm:"column:updated_since;not null"`
	ContributorCount int       `gorm:"column:contributor_count;not null"`
	OrgCount         int       `gorm:"column:org_count;not null"`
	CommitFrequency  float64   `gorm:"column:commit_frequency;not null"`
	NeedUpdate       bool      `gorm:"column:need_update;not null"`
}

func NewGitMetrics(
	Name, Owner, Source, URL, License, Languages /*, License*/ string,
	Ecosystems string,
	CreatedSince, UpdatedSince time.Time,
	ContributorCount, OrgCount int,
	CommitFrequency float64,
	NeedUpdate bool,
) GitMetrics {
	return GitMetrics{
		// gorm.Model{},
		Name,
		Owner,
		Source,
		URL,
		License,
		Ecosystems,
		Languages,
		CreatedSince,
		UpdatedSince,
		ContributorCount,
		OrgCount,
		CommitFrequency,
		NeedUpdate,
	}
}

func Repo2Metrics(r *git.Repo) GitMetrics {
	return NewGitMetrics(
		r.Name,
		r.Owner,
		r.Source,
		r.URL,
		r.License,
		r.Languages,
		r.Ecosystems,
		r.CreatedSince,
		r.UpdatedSince,
		r.ContributorCount,
		r.OrgCount,
		r.CommitFrequency,
		false,
	)
}
