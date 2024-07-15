package main

import (
	"fmt"
	"os"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_depsdev"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: show_depdev_deps <input_file_name> <output_file_name>")
		return
	}

	inputName := os.Args[1]
	outputName := os.Args[2]

	collector_depsdev.Run(inputName, outputName)
}
