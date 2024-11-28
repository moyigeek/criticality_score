/*
 * @Date: 2024-09-07 19:55:19
 * @LastEditTime: 2024-11-27 20:28:18
 * @Description:
 */
package csv

import (
	"encoding/csv"
	"os"

	"github.com/HUSTSecLab/criticality_score/pkg/collector_git/config"
)

func GetCSVInput(path string) ([][]string, error) {
	if path == "" {
		path = config.INPUT_CSV_PATH
	}

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	urls, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	return urls, nil
}

func SaveToCSV(content [][]string) error {
	file, err := os.Create(config.OUTPUT_CSV_PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(content)

	return nil
}
