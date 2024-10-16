/*
 * @Date: 2024-09-07 20:07:35
 * @LastEditTime: 2024-09-29 17:10:46
 * @Description:
 */
package sqlite

import (
	"fmt"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
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
	utils.CheckIfError(err)
	sqlDB, err := db.DB()
	utils.CheckIfError(err)
	sqlDB.SetMaxIdleConns(database.SQLITE_MAX_IDLE_CONNS)
	sqlDB.SetMaxOpenConns(database.SQLITE_MAX_OPEN_CONNS)

	return db
}

func CreateTable(db *gorm.DB) {
	err := db.AutoMigrate(&database.GitMetrics{})
	utils.CheckIfError(err)
}

func InsertTable(db *gorm.DB, metrics *database.GitMetrics) {
	db.Where(&database.GitMetrics{URL: metrics.URL}).Assign(database.GitMetrics{
		UpdatedSince:     metrics.UpdatedSince,
		ContributorCount: metrics.ContributorCount,
		OrgCount:         metrics.OrgCount,
		CommitFrequency:  metrics.CommitFrequency,
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
