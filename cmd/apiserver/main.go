package main

import (
	"github.com/HUSTSecLab/criticality_score/cmd/apiserver/internal/server"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/spf13/pflag"
)

var flagConfigPath = pflag.String("config", "config.json", "path to the config file")

func main() {
	pflag.Parse()
	logger.Config(&logger.AppLoggerConfig{
		Level:      logger.LoggerLevelInfo,
		FormatType: logger.LoggerFormatJSON,
	})
	storage.BindDefaultConfigPath("config")

	server.RegisterService()
	server.StartWebServer("0.0.0.0", 8080)
}
