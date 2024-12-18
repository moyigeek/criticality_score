package main

import (
	"flag"
	"strings"
	"github.com/HUSTSecLab/criticality_score/pkg/home2git"
)
var flagConfigPath = flag.String("config", "config.json", "path to the config file")
var choice = flag.String("choice", "", "which function to run etc. home2git, identifyCountry, identifyIndustry")
var repo = flag.String("repo", "", "input repo name")
var url = flag.String("url", "", "input LLM url")
var batchSize = flag.Int("batch", 1000, "batch size")

func main() {
	flag.Parse()
	if *choice == "home2git" {
		repolist := strings.Split(*repo, ",")
		home2git.Home2git(*flagConfigPath, repolist, *url, *batchSize)
	}
	if  *choice == "identifyIndustry" {
		home2git.IndustryID(*flagConfigPath, *url, *batchSize)
	}
}

