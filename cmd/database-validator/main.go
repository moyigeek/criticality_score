package main

import (
	"flag"
	"log"

	checkvalid "github.com/HUSTSecLab/criticality_score/pkg/database-validator"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")
var outputFile = flag.String("output", "output.csv", "path to the output file")
var checkCloneValid = flag.Bool("checkCloneValid", false, "check clone valid")
var maxThreads = flag.Int("maxThreads", 10, "max threads")

func main() {
	flag.Parse()
	storage.InitializeDatabase(*flagConfigPath)
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	invalidLinks := checkvalid.CheckVaild(db, *checkCloneValid, *maxThreads)
	checkvalid.WriteCsv(invalidLinks, *outputFile)
	log.Println("checkvalid finished")
}
