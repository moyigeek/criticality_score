/*
 * @Date: 2024-09-06 21:09:14
 * @LastEditTime: 2024-11-27 21:18:04
 * @Description: The Cli for collector
 */
package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	collector "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/logger"
	git "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/git"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/workerpool"

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

			output := make([][]string, 0)

			for _, path := range paths {
				workerpool.Go(func() {
					defer wg.Done()
					fmt.Printf("Collecting %s\n", path)
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
					repo, err := git.ParseGitRepo(r)
					if err != nil {
						logger.Panicf("Parsing %s Failed", path)
					}
					output = append(output, []string{
						repo.URL,
						repo.Name,
						repo.Owner,
						repo.Source,
						repo.License,
						fmt.Sprintf("%s", repo.Languages),
						fmt.Sprintf("%s", repo.Ecosystems),
						repo.Metrics.CreatedSince.String(),
						repo.Metrics.UpdatedSince.String(),
						fmt.Sprintf("%d", repo.Metrics.ContributorCount),
						fmt.Sprintf("%d", repo.Metrics.OrgCount),
						fmt.Sprintf("%f", repo.Metrics.CommitFrequency),
					})
					logger.Infof("%s Collected", repo.Name)
				})
			}

			wg.Wait()
			for _, o := range output {
				fmt.Printf("Repo URL: %s\n", o[0])
				fmt.Printf("Repo Name: %s   ", o[1])
				fmt.Printf("Owner: %s   ", o[2])
				fmt.Printf("Source: %s\n", o[3])
				fmt.Printf("License: %s\n", o[4])
				fmt.Printf("Languages: %s\n", o[5])
				fmt.Printf("Ecosystems: %s\n", o[6])
				fmt.Printf("Created Since: %s\n", o[7])
				fmt.Printf("Updated Since: %s\n", o[8])
				fmt.Printf("Contributor Count: %s   ", o[9])
				fmt.Printf("Org Count: %s   ", o[10])
				fmt.Printf("Commit Frequency: %s\n\n", o[11])
			}
			return nil
		}}
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
