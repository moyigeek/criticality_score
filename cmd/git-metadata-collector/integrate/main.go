package main

import (
	"log"
	"sync"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/config"
	"github.com/HUSTSecLab/criticality_score/pkg/gitfile/collector"
	git "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/repository"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/sqlutil"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/lib/pq"
	"github.com/spf13/pflag"
)

var flagJobsCount = pflag.IntP("jobs", "j", 256, "jobs count")
var flagForceUpdateAll = pflag.Bool("force-update-all", false, "force update all repositories")

func getUrls() ([]string, error) {
	conn, err := storage.GetDefaultAppDatabaseContext().GetDatabaseConnection()
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
	config.RegistCommonFlags(pflag.CommandLine)
	config.RegistGitStorageFlags(pflag.CommandLine)
	config.ParseFlags(pflag.CommandLine)

	urls, err := getUrls()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	logger.Infof("%d urls in total", len(urls))
	wg.Add(len(urls))

	ctx := storage.GetDefaultAppDatabaseContext()
	gmr := repository.NewGitMetricsRepository(ctx)

	if err != nil {
		logger.Fatal("Connecting Database Failed")
	}
	// psql.CreateTable(db)
	gopool.SetCap(int32(*flagJobsCount))

	for index, input := range urls {
		if index%10 == 0 {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}

		gopool.Go(func() {
			defer wg.Done()
			u := url.ParseURL(input)
			r, err := collector.Collect(&u, config.GetGitStoragePath())
			if err != nil {
				logger.Errorf("Collecting %s Failed", u.URL)
				return
			}
			logger.Infof("[*] %s Collected", input)

			repo, err := git.ParseRepo(r)
			if err != nil {
				logger.Errorf("Parsing %s Failed", input)
				return
			}

			err = gmr.InsertOrUpdate(&repository.GitMetric{
				GitLink:          sqlutil.ToData(input),
				CreatedSince:     sqlutil.ToNullable(repo.CreatedSince),
				UpdatedSince:     sqlutil.ToNullable(repo.UpdatedSince),
				ContributorCount: sqlutil.ToNullable(repo.ContributorCount),
				CommitFrequency:  sqlutil.ToNullable(repo.CommitFrequency),
				OrgCount:         sqlutil.ToNullable(repo.OrgCount),
				License:          sqlutil.ToNullable(pq.StringArray(repo.Licenses)),
				Language:         sqlutil.ToNullable(pq.StringArray(repo.Languages)),
			})

			if err != nil {
				logger.Errorf("Inserting %s Failed", input)
				return
			}
		})
	}
	wg.Wait()
}
