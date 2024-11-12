/*
 * @Date: 2024-09-07 16:30:21
 * @LastEditTime: 2024-09-29 14:29:34
 * @Description:
 */
package psql

import (
	"fmt"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDBFromStorageConfig() *gorm.DB {
	config := storage.GetGlobalConfig()
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.Database,
		config.Port,
		"disable",
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	utils.CheckIfError(err)
	return db
}

func InitDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.PSQL_HOST,
		config.PSQL_USER,
		config.PSQL_PASSWORD,
		config.PSQL_DATABASE_NAME,
		config.PSQL_PORT,
		config.PSQL_SSL_MODE,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	utils.CheckIfError(err)
	return db
}

func CreateTable(db *gorm.DB) {
	err := db.AutoMigrate(&database.GitMetrics{})
	utils.CheckIfError(err)
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
		NeedUpdate:	      metrics.NeedUpdate,
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
