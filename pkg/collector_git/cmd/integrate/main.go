/*
 * @Date: 2023-11-11 22:44:26
 * @LastEditTime: 2024-11-27 21:40:17
 * @Description:
 */
package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"
	psql "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database/psql"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/logger"
	git "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/workerpool"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")
var flagStoragePath = flag.String("storage", "./storage", "path to git storage location")
var flagJobsCount = flag.Int("jobs", 256, "jobs count")
var flagForceUpdateAll = flag.Bool("force-update-all", false, "force update all repositories")

func getUrls() ([]string, error) {
	conn, err := storage.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}

	var sqlStatement string

	if *flagForceUpdateAll {
		sqlStatement = `SELECT git_link from git_metrics`
	} else {
		sqlStatement = `SELECT git_link from git_metrics where need_update = true`
	}

	rows, err := conn.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	var ret []string
	for rows.Next() {
		var link string
		rows.Scan(&link)
		ret = append(ret, link)
	}
	return ret, nil
}

func main() {
	flag.Parse()
	storage.InitializeDatabase(*flagConfigPath)
	config.SetStoragetPath(*flagStoragePath)

	urls, err := getUrls()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	logger.Infof("%d urls in total", len(urls))
	wg.Add(len(urls))

	db, err := psql.InitDBFromStorageConfig()
	if err != nil {
		logger.Fatal("Connecting Database Failed")
	}
	// psql.CreateTable(db)
	workerpool.SetCap(int32(*flagJobsCount))

	for index, input := range urls {
		if index%10 == 0 {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}

		workerpool.Go(func() {
			defer wg.Done()
			u := url.ParseURL(input)
			r, err := collector.Collect(&u)
			if err != nil {
				logger.Panicf("Collecting %s Failed", u.URL)
			}
			logger.Infof("[*] %s Collected", input)

			repo, err := git.ParseGitRepo(r)
			if err != nil {
				logger.Panicf("Parsing %s Failed", input)
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
