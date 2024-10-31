package main

import (
	"flag"
	"fmt"

	"github.com/HUSTSecLab/criticality_score/pkg/collector/archlinux"
	"github.com/HUSTSecLab/criticality_score/pkg/collector/debian"
	"github.com/HUSTSecLab/criticality_score/pkg/collector/nix"
	"github.com/HUSTSecLab/criticality_score/pkg/collector/homebrew"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var (
	flagConfigPath = flag.String("config", "config.json", "path to the config file")
	flagType       = flag.String("type", "", "type of the distribution")
	flagGenDot     = flag.String("gendot", "", "output dot file")
)

func main() {
	flag.Parse()
	storage.InitializeDatabase(*flagConfigPath)

	switch *flagType {
	case "archlinux":
		archlinux.Archlinux(*flagGenDot)
	case "debian":
		debian.Debian(*flagGenDot)
	case "nix":
		if *flagGenDot == "" {
			fmt.Errorf("Nix not support gendot")
		}
		nix.Nix()
	case "homebrew":
		homebrew.Homebrew(*flagGenDot)
	}	
}
