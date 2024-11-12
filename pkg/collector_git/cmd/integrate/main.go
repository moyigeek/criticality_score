/*
 * @Date: 2023-11-11 22:44:26
 * @LastEditTime: 2024-09-29 17:37:23
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
	git "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	utils "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"
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
	utils.Info("%d urls in total", len(urls))
	wg.Add(len(urls))
	// var output [database.BATCH_SIZE]database.Metrics

	db := psql.InitDBFromStorageConfig()
	// psql.CreateTable(db)
	workerpool.SetCap(int32(*flagJobsCount))

	for _, input := range urls {
		// for index , url := range urls {

		// time.Sleep(time.Second)
		workerpool.Go(func() {
			defer wg.Done()
			// fmt.Printf("Collecting %s\n", url[0])
			u := url.ParseURL(input)

			r, err := collector.Collect(&u)
			utils.HandleErr(err, u.URL)
			if err != nil {
				r = nil
			}
			if r == nil {
				utils.Warning("[*] Cloning %s Failed at %s", input, time.Now().String())
			} else {
				utils.Info("[*] %s Cloned at %s", input, time.Now().String())
				repo := git.ParseGitRepo(r)
				if repo == nil {
					utils.Warning("[!] %s Collect Failed", input)
					return
				}
				// output[index] = database.NewMetrics(
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
			}
			// utils.Info("[*] %s Collected at %s", url,time.Now().String())
		})
	}
	wg.Wait()
	/*
		 if err := psql.BatchInsertMetrics(db,output) ; err != nil {
			 utils.Warning("Insert Failed")
			 utils.CheckIfError(err)
		 }
		 utils.Info("%d metrics inserted!",len(output))
	*/
}
