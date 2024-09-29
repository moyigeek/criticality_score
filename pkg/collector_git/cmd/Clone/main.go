/*
* @Date: 2023-11-11 22:44:26
 * @LastEditTime: 2024-09-29 17:08:07
* @Description: Just Clone
*/
package main

import (
	"os"
	"sync"
	"time"

	collector "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/file/csv"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	utils "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/workerpool"
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
	wg.Add(len(urls))

	for _, input := range urls {
		// 启动太快会被 Github 挡
		time.Sleep(3 * time.Second)
		workerpool.Go(func() {
			defer wg.Done()
			// fmt.Printf("[*] Collecting %s\n", url[0])
			u := url.ParseURL(input[0])
			r, err := collector.Collect(&u)
			utils.HandleErr(err, u.URL)
			if err != nil {
				r = nil
			}
			if r == nil {
				utils.Warning("[*] Cloning %s Failed at %s", input, time.Now().String())
			} else {
				utils.Info("[*] %s Cloned at %s", input, time.Now().String())
			}
		})
	}

	wg.Wait()
}
