// this tool is used to collect git metadata in storage path, but not clone the repository.
package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/HUSTSecLab/criticality_score/pkg/gitfile/collector"
	git "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/url"
	gitUtil "github.com/HUSTSecLab/criticality_score/pkg/gitfile/util"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var flagConfigPath = pflag.StringP("config", "c", "config.json", "path to the config file")
var flagStoragePath = pflag.StringP("storage", "s", "./storage", "path to git storage location")
var flagJobsCount = pflag.IntP("jobs", "j", 256, "jobs count")
var flagForceUpdateAll = pflag.Bool("force-update-all", false, "force update all repositories")
var flagDisableUpdateInfo = pflag.Bool("disable-update-info", false, "disable update meta, like language and license")
var flagDisableUpdateLog = pflag.Bool("disable-update-log", false, "disable update log, like commit frequency")

func getUrls() ([]string, error) {
	conn, err := storage.GetDefaultAppDatabaseConnection()
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
	pflag.Usage = func() {
		fmt.Println("This tool is used to collect git metadata in storage path, but not clone the repository.")
		pflag.PrintDefaults()
	}

	const viperStorageKey = "storage"

	pflag.Parse()
	viper.BindPFlag(viperStorageKey, pflag.Lookup("storage"))
	viper.BindEnv(viperStorageKey, "STORAGE_PATH")

	pflag.Parse()
	storage.InitializeDefaultAppDatabase(*flagConfigPath)

	urls, err := getUrls()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	logger.Infof("%d urls in total", len(urls))
	wg.Add(len(urls))

	db, err := storage.GetDefaultAppDatabaseConnection()
	if err != nil {
		logger.Fatal("Connecting Database Failed")
	}

	gopool.SetCap(int32(*flagJobsCount))

	for _, input := range urls {

		gopool.Go(func() {
			defer wg.Done()
			u := url.ParseURL(input)

			path := gitUtil.GetGitRepositoryPath(viper.GetString(viperStorageKey), &u)
			r, err := collector.Open(path)

			if err != nil || r == nil {
				logger.Errorf("Open %s failed: %s", u.URL, err)
				return
			}

			result := git.NewRepo()
			err = result.WalkRepo(r)

			if err != nil {
				logger.Errorf("WalkRepo %s failed: %s", input, err)
				return
			}

			sqlResult, err := db.Exec(`UPDATE git_metrics SET
				ecosystem = $1,
				license = $2,
				language = $3
				WHERE git_link = $4`,
				result.Ecosystems,
				result.License,
				result.Languages,
				input)

			if err != nil {
				logger.Errorf("Update database for %s failed: %v", input, err)
				return
			}

			rowAffected, err := sqlResult.RowsAffected()

			if err != nil {
				logger.Errorf("Get RowsAffected for %s Failed: %v", input, err)
				return
			}

			if rowAffected == 0 {
				logger.Warnf("Update %s failed: row affected = 0", input)
				return
			}

			logger.Infof("Success: %s", input)

		})
	}
	wg.Wait()
}
