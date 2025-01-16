package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type OpenAIClient struct {
	client  *http.Client
	apiKey  string
	prompt  string
	model   string
	enabled bool
}

func NewOpenAIClient(apiKey, prompt, model string) *OpenAIClient {
	c := &OpenAIClient{
		client: &http.Client{},
		apiKey: apiKey,
		prompt: prompt,
		model:  model,
	}

	if apiKey != "" {
		c.enabled = true
	}

	log.Printf("openai client is enabled: %v", c.enabled)

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

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (c *OpenAIClient) GetCompletions(text string) (string, error) {
	url := "https://api.proxyapi.ru/openai/v1/chat/completions"
	method := "POST"

	requestBody := RequestBody{
		Model: c.model,
		Messages: []Message{
			{
				Role:    "user",
				Content: c.prompt + "\n" + text,
			},
		},
	}

	query, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(query))
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
	return c.enabled
}
