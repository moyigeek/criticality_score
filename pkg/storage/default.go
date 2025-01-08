package storage

import (
	"fmt"

	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ViperDbConfigPathKey = "db_config_path"
	DbConfigEnv          = "APPDB_CONFIG_PATH"
)

var defaultAppDatabase AppDatabaseContext

// BindDefaultConfigPath binds the default config path to the given flag key and
// APPDB_CONFIG_PATH environment variable.
func BindDefaultConfigPath(flagKey string) {
	viper.BindEnv(DbConfigEnv)
	viper.BindPFlag(ViperDbConfigPathKey, pflag.CommandLine.Lookup(flagKey))
}

func GetDefaultConfig() (*Config, error) {
	if defaultAppDatabase == nil {
		return nil, fmt.Errorf("default app database is not initialized")
	}

	config := defaultAppDatabase.GetConfig()
	return &config, nil
}

func initDefault() (AppDatabaseContext, error) {
	var err error
	configPath := viper.GetString(ViperDbConfigPathKey)
	if configPath == "" {
		return nil, fmt.Errorf("db config path is not set")
	}
	defaultAppDatabase, err = NewAppDatabase(configPath)
	if err != nil {
		return nil, err
	}
	return defaultAppDatabase, nil
}

func GetDefaultAppDatabaseContext() AppDatabaseContext {
	if defaultAppDatabase == nil {
		_, err := initDefault()
		if err != nil {
			logger.Panic("Failed to initialize default app database", err)
		}
	}

	return defaultAppDatabase
}
