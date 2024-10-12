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

	config "github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"
	psql "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database/psql"
	git "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	utils "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/workerpool"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"

	gogit "github.com/go-git/go-git/v5"
)

var flagConfigPath = flag.String("config", "config.json", "path to the config file")

func getUrls() ([]string, error) {
	conn, err := storage.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query("SELECT git_link from git_metrics")
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

	urls, err := getUrls()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	utils.Info("%d urls in total", len(urls))
	wg.Add(len(urls))
	// var output [database.BATCH_SIZE]database.Metrics

	db := psql.InitDBFromStorageConfig()
	psql.CreateTable(db)

	for _, input := range urls {
		// for index , url := range urls {

		// time.Sleep(time.Second)
		workerpool.Go(func() {
			defer wg.Done()
			// fmt.Printf("Collecting %s\n", url[0])
			u := url.ParseURL(input)
			r, err := gogit.PlainOpen(config.STORAGE_PATH + u.Pathname)
			if err != nil {
				// r = collector.Collect(&u)
				r = nil
			}
			if r == nil {
				utils.Warning("[*] Collect %s Failed at %s", input, time.Now().String())
			} else {
				repo := git.ParseGitRepo(r)
				if repo == nil {
					utils.Warning("[!] %s Collect Failed", input)
					return
				}
				// output[index] = database.NewMetrics(
				output := database.NewMetrics(
					repo.Name,
					repo.Owner,
					repo.Source,
					repo.URL,
					repo.Metrics.CreatedSince,
					repo.Metrics.UpdatedSince,
					repo.Metrics.ContributorCount,
					repo.Metrics.OrgCount,
					repo.Metrics.CommitFrequency,
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
