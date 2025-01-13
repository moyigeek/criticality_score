package main

import (
	"github.com/HUSTSecLab/criticality_score/cmd/apiserver/internal/server"
	"github.com/HUSTSecLab/criticality_score/pkg/config"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/spf13/pflag"
)

func main() {
	config.RegistCommonFlags(pflag.CommandLine)
	config.ParseFlags(pflag.CommandLine)

	logger.Config(&logger.AppLoggerConfig{
		Level:      logger.LoggerLevelInfo,
		FormatType: logger.LoggerFormatJSON,
	})

	server.RegisterService()
	server.StartWebServer("0.0.0.0", 8080)
}
