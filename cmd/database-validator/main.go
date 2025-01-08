package main

import (
	"flag"
	"log"

	"github.com/HUSTSecLab/criticality_score/cmd/database-validator/internal/checkvalid"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")
var flagOutputFile = flag.String("output", "output.csv", "path to the output file")
var flagCheckCloneValid = flag.Bool("checkCloneValid", false, "check clone valid")
var flagMaxThreads = flag.Int("maxThreads", 10, "max threads")

func main() {
	flag.Parse()
	storage.InitializeDefaultAppDatabase(*flagConfigPath)
	db, err := storage.GetDefaultAppDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	invalidLinks := checkvalid.CheckVaild(db, *flagCheckCloneValid, *flagMaxThreads)
	checkvalid.WriteCsv(invalidLinks, *flagOutputFile)
	log.Println("checkvalid finished")
}
