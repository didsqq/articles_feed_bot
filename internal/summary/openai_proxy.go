package summary

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type OpenAIProxySummarizer struct {
	client  *http.Client
	prompt  string
	model   string
	enabled bool
	mu      sync.Mutex
}

type Response struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewOpenAIProxySummarizer(apiKey, model, prompt string) *OpenAIProxySummarizer {
	s := &OpenAIProxySummarizer{
		client: &http.Client{},
		prompt: prompt,
		model:  model,
	}

	if apiKey != "" {
		s.enabled = true
	}

	log.Printf("openai proxy summarizer is enabled: %v", s.enabled)

	return s
}

func (s *OpenAIProxySummarizer) Summarize(text string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		return "", fmt.Errorf("openai proxy summarizer is disabled")
	}

	rawSummary, err := s.GetCompletions(text)
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(rawSummary, ".") {
		return rawSummary, nil
	}

	sentences := strings.Split(rawSummary, ".")

	return strings.Join(sentences[:len(sentences)-1], ".") + ".", nil
}

func (s *OpenAIProxySummarizer) GetCompletions(text string) (string, error) {
	url := "https://api.proxyapi.ru/openai/v1/chat/completions"
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{
	  "model": "gpt-4o-mini",
	  "messages": [{"role": "user", "content": "%s\n%s"}]
  }`, s.prompt, text))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		log.Printf("[ERROR] ошибка здесь http.NewRequestWithContext(ctx, method, url, payload)")
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer sk-fBHxLeqp5AG31bPXrwMrFE3OF0Vx3hj9")

	res, err := s.client.Do(req)
	if err != nil {
		log.Printf("[ERROR] ошибка здесь s.client.Do(req)")
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[ERROR] ошибка здесь  io.ReadAll(res.Body)")
		return "", err
	}
	var resp Response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Printf("[ERROR] ошибка здесь json.Unmarshal(body, &resp)")
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
