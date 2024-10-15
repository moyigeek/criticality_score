/*
 * @Date: 2023-11-11 22:44:26
 * @LastEditTime: 2024-09-29 17:37:23
 * @Description: Collect Local Repo
 */
package main

import (
	// collector "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"os"
	"time"

	config "github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database"
	psql "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/database/psql"
	csv "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/file/csv"
	git "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	utils "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/workerpool"

	gogit "github.com/go-git/go-git/v5"

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
	urls := csv.GetCSVInput(path)
	var wg sync.WaitGroup
	utils.Info("%d urls in total", len(urls))
	wg.Add(len(urls))
	// var output [database.BATCH_SIZE]database.Metrics

	db := psql.InitDB()
	psql.CreateTable(db)

	for _, input := range urls {
		// for index , url := range urls {

		// time.Sleep(time.Second)
		workerpool.Go(func() {
			defer wg.Done()
			//fmt.Printf("Collecting %s\n", url[0])
			u := url.ParseURL(input[0])
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
				//output[index] = database.NewMetrics(
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
			//utils.Info("[*] %s Collected at %s", url,time.Now().String())
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
