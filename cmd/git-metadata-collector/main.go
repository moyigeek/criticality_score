package main

import (
	"time"

	"github.com/HUSTSecLab/criticality_score/cmd/git-metadata-collector/internal/schedule"
	"github.com/HUSTSecLab/criticality_score/cmd/git-metadata-collector/internal/task"
	"github.com/HUSTSecLab/criticality_score/pkg/config"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/spf13/pflag"
)

var flagJobsCount = pflag.IntP("jobs", "j", 256, "jobs count")
var flagForceUpdateAll = pflag.Bool("force-update-all", false, "force update all repositories")

func main() {
	config.RegistCommonFlags(pflag.CommandLine)
	config.RegistGitStorageFlags(pflag.CommandLine)
	config.ParseFlags(pflag.CommandLine)

	// psql.CreateTable(db)
	gp := gopool.NewPool("collector", int32(*flagJobsCount), &gopool.Config{})
	cnt := 0

	for {
		t, err := schedule.GetTask()
		if err != nil {
			logger.Fatalf("Failed to get task: %s", err)
		}

		// begin sleep trick
		if cnt%10 == 0 {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}
		cnt++
		// end sleep trick

		gp.Go(func() {
			task.Collect(t)
		})
	}
}
