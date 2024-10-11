package main

import (
	"flag"
	"log"

	"github.com/HUSTSecLab/criticality_score/pkg/gitmetricsync"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")

func main() {
	flag.Parse()
	storage.InitializeDatabase(*flagConfigPath)

	log.Println("Starting synchronization...")
	gitmetricsync.Run()
	log.Println("Synchronization complete.")
}
