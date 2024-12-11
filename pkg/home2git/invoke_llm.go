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

func InvokeModel(prompt string, attempts int, url string)string{
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	data := map[string]interface{}{
		"model":  "llama3.3:70b",
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

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		var jsonResponse Response
		if err := json.Unmarshal(line, &jsonResponse); err != nil {
			return "JSON decoding error"
		}
		fullResponse += jsonResponse.Response
		if jsonResponse.Done {
			break
		}
	}
	return fullResponse

}