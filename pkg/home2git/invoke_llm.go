package home2git

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
)

type Response struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func InvokeModel(prompt string, url string)string{
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	data := map[string]interface{}{
		"model":  "modelscope.cn/Qwen/Qwen2.5-32B-Instruct-GGUF:Q8_0",
		"prompt": prompt,
		"stream": true,
	}
	jsonData, _ := json.Marshal(data)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "HTTP request failed"
	}
	defer response.Body.Close()

	reader := bufio.NewReader(response.Body)

	fullResponse := ""

	var chunk []byte
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return "Error reading response"
		}
		chunk = append(chunk, line...)
		var jsonResponse Response
		if err := json.Unmarshal(chunk, &jsonResponse); err != nil {
			continue
		}
		fullResponse += jsonResponse.Response
		if jsonResponse.Done {
			break
		}
		chunk = nil
	}
	return fullResponse

}