package generate

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

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
		return "", history, fmt.Errorf(
			"failed to create chat completion: %w",
			err,
		)
	}
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: completion.Choices[0].Message.Content})
	return completion.Choices[0].Message.Content, history, nil
}
