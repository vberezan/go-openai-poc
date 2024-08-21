package openai

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

func CallOpenAIChat(request CreateCompletionsRequest, apiKey string, apiUrl string) (string, error) {
	var client = resty.New()

	var resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey).
		SetBody(request).
		Post(apiUrl)

	if err != nil {
		return "", err
	}

	var openAIResponse CreateCompletionsResponse
	if err := json.Unmarshal(resp.Body(), &openAIResponse); err != nil {
		return "", err
	}

	if openAIResponse.Error != nil {
		var errorJSON, err = json.Marshal(openAIResponse.Error)
		if err != nil {
			return "", err
		}
		return string(errorJSON), nil
	} else {
		return openAIResponse.Choices[0].Message.Content, nil
	}
}
