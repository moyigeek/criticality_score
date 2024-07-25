package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/HUSTSecLab/criticality_score/pkg/home2git"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <input_csv_file> <output_csv_file>")
		return
	}

	inputCSV := os.Args[1]
	outputCSV := os.Args[2]

	// 打开输入文件
	inFile, err := os.Open(inputCSV)
	if err != nil {
		fmt.Printf("Failed to open input file %s: %v\n", inputCSV, err)
		return
	}
	defer inFile.Close()

	// 创建CSV读取器
	reader := csv.NewReader(inFile)

	// 打开输出文件
	outFile, err := os.Create(outputCSV)
	if err != nil {
		fmt.Printf("Failed to create output file %s: %v\n", outputCSV, err)
		return
	}
	defer outFile.Close()

	// 创建CSV写入器
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	// 写入标题行到输出文件
	writer.Write([]string{"PackageName", "Homepage", "GitHub Repository"})

	// 读取每一行输入
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Failed to read from input file %s: %v\n", inputCSV, err)
			continue
		}

		if len(record) < 2 {
			continue // 跳过无效的行
		}

		packageName := record[0]
		homepageURL := record[1]
		links := []string{} // 假设从其他地方获取链接列表

		// 调用 ProcessHomepage 函数处理主页 URL
		githubURL := home2git.ProcessHomepage(packageName, links, homepageURL)

		// 将结果写入到输出文件
		writer.Write([]string{packageName, homepageURL, githubURL})
	}

	fmt.Println("Processing complete. Results saved to", outputCSV)
}
