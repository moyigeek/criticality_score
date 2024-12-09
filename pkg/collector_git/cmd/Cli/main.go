/*
 * @Date: 2024-09-06 21:09:14
 * @LastEditTime: 2024-12-09 19:31:36
 * @Description: The Cli for collector
 */
package main

import (
	"os"
	"strings"
	"sync"

	collector "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/logger"
	git "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/bytedance/gopkg/util/gopool"
	gogit "github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "collector_git",
		Usage: "Collect Git-based Repository Metrics",
		Action: func(c *cli.Context) error {
			paths := []string{}
			for i := 0; i < c.NArg(); i++ {
				paths = append(paths, c.Args().Get(i))
			}

			var wg sync.WaitGroup
			wg.Add(len(paths))

			repos := make([]*git.Repo, 0)

			for _, path := range paths {
				gopool.Go(func() {
					defer wg.Done()
					logger.Infof("Collecting %s\n", path)

					r := &gogit.Repository{}
					var err error

					if strings.Contains(path, "://") {
						u := url.ParseURL(path)
						r, err = collector.EzCollect(&u)
						if err != nil {
							logger.Panicf("Collecting %s Failed", u.URL)
						}
					} else {
						r, err = collector.Open(path)
						if err != nil {
							logger.Panicf("Opening %s Failed", path)
						}
					}

					repo, err := git.ParseRepo(r)
					if err != nil {
						logger.Panicf("Parsing %s Failed", path)
					}

					repos = append(repos, repo)
					logger.Infof("%s Collected", repo.Name)
				})
			}

			wg.Wait()
			for _, repo := range repos {
				repo.Show()
			}
			return nil
		}}
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
