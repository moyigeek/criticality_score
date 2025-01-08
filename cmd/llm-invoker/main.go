package main

import (
	"flag"
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/llm-invoker"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")
var choice = flag.String("choice", "", "which function to run etc. home2git, identifyCountry, identifyIndustry")
var repo = flag.String("repo", "", "input repo name")
var url = flag.String("url", "", "input LLM url")
var batchSize = flag.Int("batch", 1000, "batch size")
var outputCsv = flag.String("output", "", "output csv file")

func main() {
	flag.Parse()
	if *choice == "home2git" {
		repolist := strings.Split(*repo, ",")
		llm.Home2git(*flagConfigPath, repolist, *url, *batchSize, *outputCsv)
	}
	if *choice == "identifyIndustry" {
		llm.IndustryID(*flagConfigPath, *url, *batchSize, *outputCsv)
	}
}
