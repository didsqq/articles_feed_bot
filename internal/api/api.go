package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type OpenAIClient struct {
	client  *http.Client
	apiKey  string
	prompt  string
	Enabled bool
}

func NewOpenAIClient(apiKey, prompt string) *OpenAIClient {
	c := &OpenAIClient{
		client: &http.Client{},
		apiKey: apiKey,
		prompt: prompt,
	}

	if apiKey != "" {
		c.Enabled = true
	}

	log.Printf("openai client is enabled: %v", c.Enabled)

	return c
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ResponseB struct {
	Balance float64 `json:"balance"`
}

func (c *OpenAIClient) GetCompletions(text string) (string, error) {
	url := "https://api.proxyapi.ru/openai/v1/chat/completions"
	method := "POST"
	query := strings.NewReader(fmt.Sprintf(`{
	  "model": "gpt-4o-mini",
	  "messages": [{"role": "user", "content": "%s\n%s"}]
  	}`, c.prompt, text))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, query)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) GetBalance() (float64, error) {
	url := "https://api.proxyapi.ru/proxyapi/balance"
	method := "GET"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var balance ResponseB
	err = json.Unmarshal(body, &balance)
	if err != nil {
		return 0, err
	}

	return balance.Balance, nil
}

func (c *OpenAIClient) IsEnabled() bool {
	return c.Enabled
}
