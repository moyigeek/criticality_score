/*
 * @Date: 2023-11-11 22:44:26
 * @LastEditTime: 2024-11-27 21:16:28
 * @Description: Collect Local Repo
 */
package main

import (
	// collector "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"os"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"
	psql "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database/psql"
	csv "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/file/csv"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/logger"
	git "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/workerpool"

	//"fmt"
	"sync"
)

func main() {
	var path string
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else {
		path = ""
	}
	urls, err := csv.GetCSVInput(path)
	if err != nil {
		logger.Fatalf("Failed to read %s", path)
	}
	var wg sync.WaitGroup
	logger.Infof("%d urls in total", len(urls))
	wg.Add(len(urls))
	// var output [database.BATCH_SIZE]database.Metrics

	db, err := psql.InitDB()
	if err != nil {
		logger.Fatal("Failed to connect database")
	}
	psql.CreateTable(db)

	for index, input := range urls {
		if index%10 == 0 {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}

		workerpool.Go(func() {
			defer wg.Done()
			u := url.ParseURL(input[0])
			r, err := collector.Collect(&u)
			if err != nil {
				logger.Panicf("Collecting %s Failed", input)
			}

			repo, err := git.ParseGitRepo(r)
			if err != nil {
				logger.Panicf("[!] Paring %s Failed", input)
			}

			output := database.NewGitMetrics(
				repo.Name,
				repo.Owner,
				repo.Source,
				repo.URL,
				repo.Ecosystems,
				repo.Metrics.CreatedSince,
				repo.Metrics.UpdatedSince,
				repo.Metrics.ContributorCount,
				repo.Metrics.OrgCount,
				repo.Metrics.CommitFrequency,
				false,
			)
			psql.InsertTable(db, &output)

		})
	}
	wg.Wait()
}
