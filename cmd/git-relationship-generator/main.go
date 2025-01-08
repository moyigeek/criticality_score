package main

import (
	"log"

	"github.com/HUSTSecLab/criticality_score/cmd/git-relationship-generator/internal/pkgdep2git"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/spf13/pflag"
)

var flagConfigPath = pflag.String("config", "config.json", "path to the config file")
var batchSize = pflag.Int("batch", 1000, "batch size for updating scores")

func main() {
	pflag.Parse()
	storage.BindDefaultConfigPath("config")
	db, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	depMap := pkgdep2git.FetchAlldep(db)
	gitdepMap := pkgdep2git.GenGitDep(db, depMap)
	err = pkgdep2git.BatchUpdate(db, *batchSize, gitdepMap)
	if err != nil {
		log.Fatalf("Error updating database: %v", err)
	}
}
