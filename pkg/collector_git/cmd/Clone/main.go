/*
* @Date: 2023-11-11 22:44:26
 * @LastEditTime: 2024-11-27 21:12:49
* @Description: Just Clone
*/
package main

import (
	"log"
	"os"
	"sync"
	"time"

	collector "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/collector"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/file/csv"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/logger"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/workerpool"
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
		log.Fatalf("Failed to read %s", path)
	}
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for index, input := range urls {
		if index%10 == 0 {
			time.Sleep(5 * time.Second)
		} else {
			time.Sleep(2 * time.Second)
		}

		workerpool.Go(func() {
			defer wg.Done()
			// fmt.Printf("[*] Collecting %s\n", url[0])
			u := url.ParseURL(input[0])
			_, err := collector.Collect(&u)
			if err != nil {
				logger.Panicf("Cloning %s Failed", input)
			} else {
				logger.Infof("%s Cloned", input)
			}
		})
	}

	wg.Wait()
}
