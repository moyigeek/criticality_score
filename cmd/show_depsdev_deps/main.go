package main

import (
	"flag"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_depsdev"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")

func main() {
	flag.Parse()
	storage.InitializeDatabase(*flagConfigPath)

	collector_depsdev.Run(*flagConfigPath)
}
