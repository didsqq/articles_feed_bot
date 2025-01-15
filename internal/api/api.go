package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenAIClient struct {
	client *http.Client
	apiKey string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		client: &http.Client{},
		apiKey: apiKey,
	}
}

type Response struct {
	Balance float64 `json:"balance"`
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

	var balance Response
	err = json.Unmarshal(body, &balance)
	if err != nil {
		return 0, err
	}

	return balance.Balance, nil
}
