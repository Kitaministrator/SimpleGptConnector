package SimpleGptConnector

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CustomProxy struct {
	ReplacedDomain string
}

type OpenaiConfig struct {
	ApiDomain string
	ApiKey    string
}

var customProxy CustomProxy

var openaiConfig OpenaiConfig

func CommunicateWithOpenAI(apiKey string, messages Message) (string, error) {

	// Debug
	defer printDebug()

	// Replace the domain of OpenAI API
	setApiConfig()

	// construct the request url
	apiAddr := openaiConfig.ApiDomain + "/v1/chat/completions"

	// construct the request body
	requestBody, _ := json.Marshal(map[string]interface{}{
		"model":       "gpt-3.5-turbo",
		"messages":    []Message{messages},
		"temperature": 0.7, //float64
		"max_tokens":  100, //int64
		// "top_p":             1,      //float64
		// "n":                 1,      //int64
		// "stream":            false,  //bool
		// "stop":              nil, //string or array
		// "presence_penalty":  0,      //float64
		// "frequency_penalty": 0,      //float64
		// "logit_bias":        nil, //map[string]interface{}
		// "user":              nil, //string
	})

	// construct the request
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiAddr, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// send out the request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// parse the response body
	var response map[string]interface{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}

	// extract the open ai message context
	// choices := response["choices"].([]interface{})
	// if len(choices) > 0 {
	// 	text := choices[0].(map[string]interface{})["text"].(string)
	// 	return text, nil
	// } else {
	// 	return "", nil
	// }
	choices, ok := response["choices"].([]interface{})
	if ok && len(choices) > 0 {

		/// Debug
		log.Println("len(choices) > 0")
		for k, v := range choices {
			log.Printf("k: %v, v: %v", k, v)
		}
		/// End of debug block

		text := choices[0].(map[string]interface{})["text"].(string)
		return text, nil
	} else {
		log.Println("len(choices) = 0")
		return "", nil
	}
}

func setApiConfig() {

	customProxy.ReplacedDomain = os.Getenv("REPLACED_DOMAIN")

	if customProxy.ReplacedDomain == "" {
		openaiConfig.ApiDomain = "https://api.openai.com"
	} else {
		openaiConfig.ApiDomain = customProxy.ReplacedDomain
	}

}

func printDebug() {
	log.Println("========  Start to print Debug Info  ========")
	log.Println("customProxy.ReplacedDomain: ", customProxy.ReplacedDomain)
	log.Println("openaiConfig.ApiDomain: ", openaiConfig.ApiDomain)
	log.Println("========  End of printing Debug Info  ========")
}
