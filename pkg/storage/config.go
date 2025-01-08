package storage

import (
	"encoding/json"
	"os"
)

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
