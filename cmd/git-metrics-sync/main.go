package main

import (
	"flag"
	"log"

	"github.com/HUSTSecLab/criticality_score/cmd/git-metrics-sync/internal/gmsync"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")
var batchSize = flag.Int("batch", 1000, "batch size")

func main() {
	flag.Parse()
	storage.InitializeDefaultAppDatabase(*flagConfigPath)

	log.Println("Starting synchronization...")
	gmsync.Run()
	log.Println("Synchronization complete.")
	gmsync.Union_repo(*batchSize)
}
