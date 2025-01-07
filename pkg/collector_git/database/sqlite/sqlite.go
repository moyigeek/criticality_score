/*
 * @Date: 2024-09-07 20:07:35
 * @LastEditTime: 2025-01-07 19:03:54
 * @Description:
 */
package sqlite

import (
	"fmt"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/database"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	var dsn string

	if config.SQLITE_USER == "" && config.SQLITE_PASSWORD == "" {
		dsn = config.SQLITE_DATABASE_PATH
	} else {
		dsn = fmt.Sprintf(
			"%s:%s@%s",
			config.SQLITE_USER,
			config.SQLITE_PASSWORD,
			config.SQLITE_DATABASE_PATH,
		)
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(database.SQLITE_MAX_IDLE_CONNS)
	sqlDB.SetMaxOpenConns(database.SQLITE_MAX_OPEN_CONNS)

	return db, err
}

func CreateTable(db *gorm.DB) error {
	return db.AutoMigrate(&database.GitMetrics{})
}

func InsertTable(db *gorm.DB, metrics *database.GitMetrics) {
	db.Where(&database.GitMetrics{URL: metrics.URL}).
		Assign(database.GitMetrics{
			Name:             metrics.Name,
			Owner:            metrics.Owner,
			Source:           metrics.Source,
			URL:              metrics.URL,
			License:          metrics.License,
			Ecosystems:       metrics.Ecosystems,
			Languages:        metrics.Languages,
			CreatedSince:     metrics.CreatedSince,
			UpdatedSince:     metrics.UpdatedSince,
			ContributorCount: metrics.ContributorCount,
			OrgCount:         metrics.OrgCount,
			CommitFrequency:  metrics.CommitFrequency,
			NeedUpdate:       metrics.NeedUpdate,
		}).FirstOrCreate(metrics)
}
