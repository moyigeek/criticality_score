package storage

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"context"
	_ "github.com/lib/pq"
	"github.com/go-redis/redis/v8"
)

var globalConfig Config

type Config struct {
	Database    string `json:"database"`
	User        string `json:"user"`
	Password    string `json:"password"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	GitHubToken string `json:"GitHubToken"`
	GitLabToken string `json:"GitLabToken"`
	Redispass   string `json:"redispass"`
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

func InitRedis()(*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, 
	})
	return rdb, nil
}

func GetDatabaseConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		globalConfig.Host, globalConfig.Port, globalConfig.User, globalConfig.Password, globalConfig.Database)
	db, err := sql.Open("postgres", connStr)
	return db, err
}

func SetKeyValue(rdb *redis.Client, key, value string) error {
	err := rdb.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("could not set key '%s': %v", key, err)
	}
	return nil
}

func GetKeyValue(rdb *redis.Client, key string) (string, error) {
	val, err := rdb.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key '%s' does not exist", key)
		}
		return "", fmt.Errorf("could not get key '%s': %v", key, err)
	}
	return val, nil
}

func PersistData(rdb *redis.Client) error {
	err := rdb.BgSave(context.Background()).Err()
	if err != nil {
		return fmt.Errorf("could not trigger RDB save: %v", err)
	}
	fmt.Println("RDB persistence triggered")
	return nil
}

func GetGlobalConfig() Config {
	return globalConfig
}

func InitDatabase(configPath string) error {
	var err error
	globalConfig, err = loadConfig(configPath)
	if err != nil {
		fmt.Errorf("Failed to load config:", err)
		return err
	}
	return nil
}