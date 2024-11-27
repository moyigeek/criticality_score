/*
 * @Date: 2024-09-07 20:07:35
 * @LastEditTime: 2024-11-27 20:27:00
 * @Description:
 */
package sqlite

import (
	"fmt"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"

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
	db.Where(&database.GitMetrics{URL: metrics.URL}).Assign(database.GitMetrics{
		CreatedSince:     metrics.CreatedSince,
		UpdatedSince:     metrics.UpdatedSince,
		ContributorCount: metrics.ContributorCount,
		OrgCount:         metrics.OrgCount,
		CommitFrequency:  metrics.CommitFrequency,
		Name:             metrics.Name,
		Owner:            metrics.Owner,
		Source:           metrics.Source,
		URL:              metrics.URL,
		Ecosystems:       metrics.Ecosystems,
		NeedUpdate:       metrics.NeedUpdate,
	}).FirstOrCreate(metrics)
}

func BatchInsertMetrics(db *gorm.DB, metrics [database.BATCH_SIZE]database.GitMetrics) error {

	tx := db.Begin()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(&metrics).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
