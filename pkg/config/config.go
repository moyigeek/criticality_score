package config

import (
	"os"

	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/spf13/viper"
)

func readPasswordFromFile(file string) string {
	content, err := os.ReadFile(file)
	if err != nil {
		return ""
	}
	return string(content)
}

func GetDatabaseConfig() *storage.Config {
	if viper.GetString("db.password") == "" && viper.GetString("db.password-file") != "" {
		viper.Set("db.password", readPasswordFromFile(viper.GetString("db.password-file")))
	}

	return &storage.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Database: viper.GetString("db.database"),
		UseSSL:   viper.GetBool("db.use-ssl"),
	}
}

func GetLogConfig() *logger.AppLoggerConfig {
	var level logger.LoggerLevel
	var format logger.LoggerFormatType
	var ttype logger.LoggerOutput

	viperLevel := viper.GetString("log.level")
	switch viperLevel {
	case "trace":
		level = logger.LoggerLevelTrace
	case "debug":
		level = logger.LoggerLevelDebug
	case "info":
		level = logger.LoggerLevelInfo
	case "warn":
		level = logger.LoggerLevelWarn
	case "error":
		level = logger.LoggerLevelError
	case "fatal":
		level = logger.LoggerLevelFatal
	case "panic":
		level = logger.LoggerLevelPanic
	default:
		level = logger.LoggerLevelInfo
	}

	viperFormat := viper.GetString("log.format")

	switch viperFormat {
	case "text":
		format = logger.LoggerFormatText
	case "cli":
		format = logger.LoggerFormatCliTool
	case "json":
		format = logger.LoggerFormatJSON
	default:
		format = logger.LoggerFormatJSON
	}

	viperType := viper.GetString("log.type")

	switch viperType {
	case "console":
		ttype = logger.LoggerOutputStderr
	case "file":
		ttype = logger.LoggerOutputFile
	case "es":
		ttype = logger.LoggerOutputElasticSearch
	default:
		ttype = logger.LoggerOutputStderr
	}

	return &logger.AppLoggerConfig{
		Level:         level,
		FormatType:    format,
		Output:        ttype,
		OutputPath:    viper.GetString("log.path"),
		OutputEsURL:   viper.GetString("log.es-url"),
		OutputEsIndex: viper.GetString("log.es-index"),
	}

}

func GetGithubToken() string {
	return viper.GetString("token.github")
}

func GetGitStoragePath() string {
	return viper.GetString("git.storage")
}
