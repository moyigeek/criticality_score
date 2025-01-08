package main

import (
	"flag"

	"github.com/HUSTSecLab/criticality_score/cmd/apiserver/internal/server"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")

func main() {
	flag.Parse()
	storage.InitializeDefaultAppDatabase(*flagConfigPath)

	server.RegisterService()
	server.StartWebServer("0.0.0.0", 8080)
}
