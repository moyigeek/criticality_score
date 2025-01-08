package main

import (
	"log"

	"github.com/HUSTSecLab/criticality_score/cmd/database-validator/internal/checkvalid"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/spf13/pflag"
)

var flagConfigPath = pflag.String("config", "config.json", "path to the config file")
var flagOutputFile = pflag.String("output", "output.csv", "path to the output file")
var flagCheckCloneValid = pflag.Bool("checkCloneValid", false, "check clone valid")
var flagMaxThreads = pflag.Int("maxThreads", 10, "max threads")

func main() {
	pflag.Parse()
	storage.BindDefaultConfigPath(*flagConfigPath)
	db, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()
	invalidLinks := checkvalid.CheckVaild(db, *flagCheckCloneValid, *flagMaxThreads)
	checkvalid.WriteCsv(invalidLinks, *flagOutputFile)
	log.Println("checkvalid finished")
}
