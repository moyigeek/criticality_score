/*
 * @Date: 2024-09-07 20:46:04
 * @LastEditTime: 2024-09-29 14:31:08
 * @Description:
 */
package database

import (
	"time"

	"gorm.io/gorm"
)

const (
	SQLITE_MAX_IDLE_CONNS int = 10
	SQLITE_MAX_OPEN_CONNS int = 100
	BATCH_SIZE                = 500
)

type GitMetrics struct {
	gorm.Model
	Name             string    `gorm:"column:_name;not null"`
	Owner            string    `gorm:"column:_owner;not null"`
	Source           string    `gorm:"column:_source;not null"`
	URL              string    `gorm:"column:git_link;not null"` //* `gorm:"unique;not null"`
	Ecosystems       []string  `gorm:"column:ecosystem;not null"`
	CreatedSince     time.Time `gorm:"column:created_since;not null"`
	UpdatedSince     time.Time `gorm:"column:updated_since;not null"`
	ContributorCount int       `gorm:"column:contributor_count;not null"`
	OrgCount         int       `gorm:"column:org_count;not null"`
	CommitFrequency  float64   `gorm:"column:commit_frequency;not null"`
	// License          string    `gorm:"not null"`
	//Languages []string `gorm:"not null"`
}

func NewGitMetrics(
	Name, Owner, Source, URL /*, License*/ string,
	Ecosystems []string,
	CreatedSince, UpdatedSince time.Time,
	ContributorCount, OrgCount int,
	CommitFrequency float64,
) GitMetrics {
	return GitMetrics{
		gorm.Model{},
		Name,
		Owner,
		Source,
		URL,
		Ecosystems,
		CreatedSince,
		UpdatedSince,
		ContributorCount,
		OrgCount,
		CommitFrequency,
		//		License,
	}
}
