package llm

import "github.com/sashabaranov/go-openai"

// CreateClient creates a new client for the given api key.
func CreateClient(baseURL string, apiKey string) *openai.Client {
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = baseURL
	cfg.APIVersion = string(openai.APITypeOpenAI)
	cfg.APIType = openai.APITypeOpenAI
	return openai.NewClientWithConfig(cfg)
}
