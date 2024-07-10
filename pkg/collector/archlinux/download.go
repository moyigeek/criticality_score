package archlinux

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	
	"golang.org/x/net/html"
)

// Fetches the HTML page and returns the response body
func fetchHTML(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to fetch %s: %s", url, resp.Status)
	}
	return resp, nil
}

// Parses the HTML and extracts folder names
func extractFolderNames(url string) ([]string, error) {
	resp, err := fetchHTML(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var folders []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.HasSuffix(a.Val, "/") && !strings.HasPrefix(a.Val, "../") {
					folders = append(folders, strings.TrimSuffix(a.Val, "/"))
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return folders, nil
}

// Downloads a file from the given URL and saves it to the specified path
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func DownloadFiles() {
	baseURL := "https://mirrors.hust.edu.cn/archlinux/"
	downloadDir := "./download"

	// Create the download directory if it doesn't exist
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		err = os.Mkdir(downloadDir, 0755)
		if err != nil {
			fmt.Printf("Error creating download directory: %v\n", err)
			return
		}
	}

	folders, err := extractFolderNames(baseURL)
	if err != nil {
		fmt.Printf("Error extracting folder names: %v\n", err)
		return
	}

	for _, folder := range folders {
		filesURL := fmt.Sprintf("%s%s/os/x86_64/%s.files.tar.gz", baseURL, folder, folder)
		filepath := fmt.Sprintf("%s/%s.files.tar.gz", downloadDir, folder)
		fmt.Printf("Downloading %s to %s\n", filesURL, filepath)
		err := downloadFile(filesURL, filepath)
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n", filesURL, err)
		} else {
			fmt.Printf("Downloaded %s successfully\n", filesURL)
		}
	}
}
