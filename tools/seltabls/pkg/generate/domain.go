package generate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/sashabaranov/go-openai"
)

// writeFile writes a file to the given path
func writeFile(name string, content string) error {
	f, err := os.Create(name + ".go")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// Chat is a struct for a chat
func Chat(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt string,
) (string, []openai.ChatCompletionMessage, error) {
	completion, err := client.CreateChatCompletion(
		ctx, openai.ChatCompletionRequest{
			Model: model,
			Messages: append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			}),
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: "json",
			},
		})
	if err != nil {
		return "", history, fmt.Errorf("failed to create chat completion: %w", err)
	}
	content := completion.Choices[0].Message.Content
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content})
	return content, history, nil
}

// GetURL gets the url and returns the body
func GetURL(url string) ([]byte, error) {
	cli := http.DefaultClient
	resp, err := cli.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get url: %s", resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return body, nil
}

// isURL checks if the string is a valid URL
func IsURL(s string) error {
	_, err := url.ParseRequestURI(s)
	return err
}

// Generatable is an interface for a generatable
type Generatable interface {
	Generate() error
}
