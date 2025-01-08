package main

import (
	"log"

	"github.com/HUSTSecLab/criticality_score/cmd/git-metrics-sync/internal/gmsync"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/spf13/pflag"
)

var flagConfigPath = pflag.String("config", "config.json", "path to the config file")
var batchSize = pflag.Int("batch", 1000, "batch size")

func main() {
	pflag.Parse()
	storage.BindDefaultConfigPath("config")

	log.Println("Starting synchronization...")
	gmsync.Run()
	log.Println("Synchronization complete.")
	gmsync.Union_repo(*batchSize)
}
