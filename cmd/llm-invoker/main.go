package main

import (
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/llm"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/spf13/pflag"
)

var flagConfigPath = pflag.String("config", "config.json", "path to the config file")
var choice = pflag.String("choice", "", "which function to run etc. home2git, identifyCountry, identifyIndustry")
var repo = pflag.String("repo", "", "input repo name")
var url = pflag.String("url", "", "input LLM url")
var batchSize = pflag.Int("batch", 1000, "batch size")
var outputCsv = pflag.String("output", "", "output csv file")

func main() {
	storage.BindDefaultConfigPath("config")

	pflag.Parse()
	if *choice == "home2git" {
		repolist := strings.Split(*repo, ",")
		llm.Home2git(storage.GetDefaultAppDatabaseContext(), *flagConfigPath, repolist, *url, *batchSize, *outputCsv)
	}
	if *choice == "identifyIndustry" {
		llm.IndustryID(storage.GetDefaultAppDatabaseContext(), *flagConfigPath, *url, *batchSize, *outputCsv)
	}
}
