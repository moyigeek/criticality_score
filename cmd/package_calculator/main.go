package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/package_calculator"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")
var flagRepoName = flag.String("repo", "", "name of the repository")

func main() {
	flag.Parse()

	if *flagRepoName == "" {
		log.Fatal("Repository name must be provided")
	}

	storage.InitializeDatabase(*flagConfigPath)
	db, err := storage.GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	query := fmt.Sprintf("SELECT frompackage, topackage FROM %s_relationships", *flagRepoName)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error querying database: %v", err)
	}
	defer rows.Close()

	if err := package_calculator.CalculatePackages(rows); err != nil {
		log.Fatalf("Error calculating packages: %v", err)
	}
}