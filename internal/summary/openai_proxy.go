package summary

import (
	"fmt"
	"strings"
	"sync"
)

type Client interface {
	GetCompletions(text string) (string, error)
	IsEnabled() bool
}

type OpenAIProxySummarizer struct {
	client Client
	model  string
	mu     sync.Mutex
}

func NewOpenAIProxySummarizer(client Client, model string) *OpenAIProxySummarizer {
	s := &OpenAIProxySummarizer{
		client: client,
		model:  model,
	}

	return s
}

func (s *OpenAIProxySummarizer) Summarize(text string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.client.IsEnabled() {
		return "", fmt.Errorf("openai api is disabled")
	}

	rawSummary, err := s.client.GetCompletions(text)
	if err != nil {
		return "", err
	}

	if strings.HasSuffix(rawSummary, ".") {
		return rawSummary, nil
	}

	sentences := strings.Split(rawSummary, ".")

	return strings.Join(sentences[:len(sentences)-1], ".") + ".", nil
}
