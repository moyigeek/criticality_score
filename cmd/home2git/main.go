package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"flag"
	"net/http"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/home2git"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
)
var flagConfigPath = flag.String("config", "config.json", "path to the config file")
func check(githubURL string)string{
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(githubURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return githubURL
	}
	return ""
}

func fetchAllLinks() ([][]string, error) {
	db, err := storage.GetDatabaseConnection()
	rows, err := db.Query("SELECT package, homepage FROM gentoo_packages union select package, homepage from homebrew_packages union select package, homepage from nix_packages")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links [][]string
	for rows.Next() {
		var packageName, link string
		if err := rows.Scan(&packageName, &link); err != nil {
			return nil, err
		}
		links = append(links, []string{link, packageName})
	}
	fmt.Println("Fetched", len(links), "links")
	return links, nil
}
func main() {
	flag.Parse()
	err := storage.InitializeDatabase(*flagConfigPath)
	if err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <output_csv_file>")
		return
	}
	outputCSV := os.Args[1]
	var PackageList [][]string
	PackageList, _ = fetchAllLinks()

	outFile, err := os.Create(outputCSV)
	if err != nil {
		fmt.Printf("Failed to create output file %s: %v\n", outputCSV, err)
		return
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)

	if err := writer.Write([]string{"PackageName", "Homepage", "GitHub Repository"}); err != nil {
		fmt.Printf("Error writing header: %v\n", err)
		return
	}
	writer.Flush()

	for i := 0; i < len(PackageList); i++ {
		homepageURL := PackageList[i][0]
		packageName := PackageList[i][1]

		htmlContent, err := home2git.DownloadHTML(homepageURL)
		if err != nil {
			continue
		}

		links, _ := home2git.FindLinksInHTML(homepageURL, htmlContent, 1)
		githubURL := home2git.ProcessHomepage(packageName, links, homepageURL)
		res := check(githubURL)
		if res == "" {
			continue
		}

		if err := writer.Write([]string{packageName, homepageURL, res}); err != nil {
			fmt.Printf("Error writing row: %v\n", err)
			continue
		}
		writer.Flush()
	}

	fmt.Println("Processing complete. Results saved to", outputCSV)
}

