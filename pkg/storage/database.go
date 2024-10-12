package storage

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	_ "github.com/lib/pq"
)

var globalConfig Config

type Config struct {
	Database    string `json:"database"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	GitHubToken string `json:"GitHubToken"`
}

func loadConfig(configPath string) (Config, error) {
	var config Config
	file, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	return config, err
}

func InitializeDatabase(configPath string) error {
	flag.Parse()
	var err error
	globalConfig, err = loadConfig(configPath)
	if err != nil {
		fmt.Errorf("Failed to load config:", err)
		return err
	}
	return nil
}

func GetDatabaseConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		globalConfig.Host, globalConfig.Port, globalConfig.User, globalConfig.Password, globalConfig.Database)
	db, err := sql.Open("postgres", connStr)
	return db, err
}

func GetGlobalConfig() Config {
	return globalConfig
}
