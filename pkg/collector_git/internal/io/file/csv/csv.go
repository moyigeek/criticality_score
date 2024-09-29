/*
 * @Date: 2024-09-07 19:55:19
 * @LastEditTime: 2024-09-27 22:00:17
 * @Description:
 */
package csv

import (
	"encoding/csv"
	"os"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/internal/utils"
)

func GetCSVInput(path string) [][]string {
	if path == "" {
		path = config.INPUT_CSV_PATH
	}
	file, err := os.Open(path)
	utils.CheckIfError(err)
	defer file.Close()
	// for _, item := range record {
	// 	fmt.Println(item)
	//}
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	urls, err := reader.ReadAll()
	utils.CheckIfError(err)
	return urls
}

func SaveToCSV(content [][]string) {
	file, err := os.Create(config.OUTPUT_CSV_PATH)
	utils.CheckIfError(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(content)
}
