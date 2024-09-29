/*
 * @Date: 2023-11-11 22:44:26
 * @LastEditTime: 2024-09-29 17:17:04
 * @Description: Just Count downloaded repos
 */
package main

import (
	"fmt"
	"os"

	config "github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	csv "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/io/file/csv"
	url "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/parser/url"
	utils "github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"

	"github.com/go-git/go-git/v5"
	//"fmt"
)

func main() {
	var path string
	if len(os.Args) == 2 {
		path = os.Args[1]
	} else {
		path = ""
	}
	inputs := csv.GetCSVInput(path)
	var count int
	for _, input := range inputs {
		u := url.ParseURL(input[0])
		path = config.STORAGE_PATH + u.Pathname
		//fmt.Println(path)
		_, err := git.PlainOpen(path)
		if err == git.ErrRepositoryNotExists {
			utils.Info("%s Not Collected", input[0])
			continue
		}
		count++
	}
	fmt.Println(count)
}
