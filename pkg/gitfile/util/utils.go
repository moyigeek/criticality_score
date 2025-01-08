package util

import (
	"encoding/csv"
	"os"
	"path"

	"github.com/HUSTSecLab/criticality_score/pkg/gitfile/parser/url"
)

func GetCSVInput(path string) ([][]string, error) {
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

func Save2CSV(outputPath string, content [][]string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(content)

	return nil
}

func GetGitRepositoryPath(storagePath string, u *url.RepoURL) string {
	// join path
	return path.Join(storagePath, u.Resource, u.Pathname)
}
