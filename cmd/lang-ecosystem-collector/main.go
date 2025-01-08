package main

import (
	"github.com/HUSTSecLab/criticality_score/pkg/depsdev"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/spf13/pflag"
)

var (
	flagConfigPath    = pflag.String("config", "config.json", "path to the config file")
	flagBatchSize     = pflag.Int("batch", 100, "batch size")
	workerCount       = pflag.Int("workers", 10, "number of workers")
	calculatePageRank = pflag.Bool("pagerank", false, "calculate page rank")
)

func main() {
	pflag.Parse()
	storage.BindDefaultConfigPath("config")

	depsdev.Depsdev(*flagConfigPath, *flagBatchSize, *workerCount, *calculatePageRank)
}
