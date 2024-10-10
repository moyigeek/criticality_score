package main

import (
	"fmt"
	"os"

	"github.com/HUSTSecLab/criticality_score/pkg/collector/archlinux"
	"github.com/HUSTSecLab/criticality_score/pkg/collector/debian"
	"github.com/HUSTSecLab/criticality_score/pkg/collector/nix"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: main <archlinux|debian|nix> [gendot <output.dot>]")
		return
	}

	switch os.Args[1] {
	case "archlinux":
		if len(os.Args) == 4 && os.Args[2] == "gendot" {
			archlinux.Archlinux(os.Args[3])
		} else if len(os.Args) == 2 {
			archlinux.Archlinux("")
		} else {
			fmt.Println("Usage: main archlinux [gendot <output.dot>]")
		}
	case "debian":
		if len(os.Args) == 4 && os.Args[2] == "gendot" {
			debian.Debian(os.Args[3])
		} else if len(os.Args) == 2 {
			debian.Debian("")
		} else {
			fmt.Println("Usage: main debian [gendot <output.dot>]")
		}
	case "nix":
		if len(os.Args) == 2 {
			nix.Nix()
		} else {
			fmt.Println("Usage: main debian [gendot <output.dot>]")
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println("Usage: main <archlinux|debian|nix> [gendot <output.dot>]")
	}
}
