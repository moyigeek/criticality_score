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

type Metrics struct {
	gorm.Model
	Name             string    `gorm:"not null"`
	Owner            string    `gorm:"not null"`
	Source           string    `gorm:"not null"`
	URL              string    `gorm:"not null"` //* `gorm:"unique;not null"`
	CreatedSince     time.Time `gorm:"not null"`
	UpdatedSince     time.Time `gorm:"not null"`
	ContributorCount int       `gorm:"not null"`
	OrgCount         int       `gorm:"not null"`
	CommitFrequency  float64   `gorm:"not null"`
}

func NewMetrics(
	Name, Owner, Source, URL string,
	CreatedSince, UpdatedSince time.Time,
	ContributorCount, OrgCount int,
	CommitFrequency float64,
) Metrics {
	return Metrics{
		gorm.Model{},
		Name,
		Owner,
		Source,
		URL,
		CreatedSince,
		UpdatedSince,
		ContributorCount,
		OrgCount,
		CommitFrequency,
	}
}
