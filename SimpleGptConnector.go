package SimpleGptConnector

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Message struct {
	Role    string "json:\"role\""
	Content string "json:\"content\""
}

func CommunicateWithOpenAI(apiKey string, message Message) (string, error) {
	// 构造请求体
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       "davinci",
		"prompt":      message,
		"max_tokens":  10,
		"temperature": 0.5,
	})
	if err != nil {
		return "", err
	}

	// 发送 POST 请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析 JSON 响应体
	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}

	// 提取 OpenAI API 的回答内容
	choices := response["choices"].([]interface{})
	if len(choices) > 0 {
		text := choices[0].(map[string]interface{})["text"].(string)
		return text, nil
	} else {
		return "", nil
	}
}
