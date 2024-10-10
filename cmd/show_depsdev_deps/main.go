package main

import (
	"fmt"
	"os"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_depsdev"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: show_depdev_deps config.json")
		return
	}

	config := os.Args[1]

	collector_depsdev.Run(config)
}
