package main

import (
	"github.com/HUSTSecLab/criticality_score/pkg/config"
	"github.com/HUSTSecLab/criticality_score/pkg/depsdev"
	"github.com/spf13/pflag"
)

var (
	flagBatchSize     = pflag.Int("batch", 100, "batch size")
	workerCount       = pflag.Int("workers", 50, "number of workers")
	calculatePageRank = pflag.Bool("pagerank", false, "calculate page rank")
)

func main() {
	config.RegistCommonFlags(pflag.CommandLine)
	config.ParseFlags(pflag.CommandLine)

	depsdev.Depsdev(*flagBatchSize, *workerCount, *calculatePageRank)
}
