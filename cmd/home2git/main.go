package main

import (
	"flag"
	"github.com/HUSTSecLab/criticality_score/pkg/home2git"
)
var flagConfigPath = flag.String("config", "config.json", "path to the config file")
var choice = flag.String("choice", "", "which function to run etc. home2git, identifyCountry, identifyIndustry")
var outputCSV = flag.String("out", "output.csv", "output file")
var url = flag.String("url", "", "input LLM url")
func main() {
	flag.Parse()
	if *choice == "home2git" {
		home2git.Home2git(*flagConfigPath, *outputCSV, *url)
	}
	if  *choice == "identifyCountry" {
	}
}

