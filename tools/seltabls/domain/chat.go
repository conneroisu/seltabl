package domain

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

var client *openai.Client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

// ChatNParams is a struct for the chat n params
type ChatNParams struct {
	Histories [][]openai.ChatCompletionMessage
	Prompts   []string
}

// Len returns the length of the slice.
func (p ChatNParams) Len() int {
	return len(p.Prompts)
}

// ChatJSON is a struct for a chat that returns json.
func ChatJSON(
	ctx context.Context,
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
		return "", history, err
	}
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: completion.Choices[0].Message.Content})
	return completion.Choices[0].Message.Content, history, nil
}

// ChatTeacher is a method to get prompting help from a teacher who sees the
// whole conversation and is able to explain the prompt or error.
func ChatTeacher(
	ctx context.Context,
	model string,
	studentHistory []openai.ChatCompletionMessage,
	prompt string,
) (out string, history []openai.ChatCompletionMessage, err error) {
	completion, err := client.CreateChatCompletion(
		ctx, openai.ChatCompletionRequest{
			Model: model,
			Messages: append([]openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			}, studentHistory...),
		},
	)
	if err != nil {
		return "", history, err
	}
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: completion.Choices[0].Message.Content})
	return completion.Choices[0].Message.Content, history, nil
}

// Chat is a struct for a chat that returns a string.
func Chat(
	ctx context.Context,
	model string,
	history []openai.ChatCompletionMessage,
	prompt string,
) (string, []openai.ChatCompletionMessage, error) {
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: prompt,
	})
	completion, err := client.CreateChatCompletion(
		ctx, openai.ChatCompletionRequest{
			Model:    model,
			Messages: history,
		})
	if err != nil {
		return "", history, err
	}
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: completion.Choices[0].Message.Content})
	return completion.Choices[0].Message.Content, history, nil
}

// ChatN is a generic number method to get prompting help from a teacher who sees the
func ChatN(
	ctx context.Context,
	model string,
	params ChatNParams,
) ([]string, [][]openai.ChatCompletionMessage, error) {
	var outs []string
	var histories [][]openai.ChatCompletionMessage = make([][]openai.ChatCompletionMessage, len(params.Prompts))
	var err error
	eg, _ := errgroup.WithContext(ctx)
	for i := range params.Len() {
		eg.Go(func() error {
			var out string
			var hist []openai.ChatCompletionMessage
			out, hist, err = Chat(
				ctx,
				model,
				params.Histories[i],
				params.Prompts[i],
			)
			if err != nil {
				return fmt.Errorf("failed to chat with llm provider: %w", err)
			}
			histories[i] = hist
			outs = append(outs, out)
			return nil
		})
	}
	return outs, histories, nil
}

// ChatNTeacher is a generic number method to get prompting help from a teacher who sees the
// whole conversation and is able to explain the prompt or error.
func ChatNTeacher(
	ctx context.Context,
	model string,
	params ChatNParams,
) ([]string, [][]openai.ChatCompletionMessage, error) {
	var outs []string
	var histories [][]openai.ChatCompletionMessage = make([][]openai.ChatCompletionMessage, len(params.Prompts))
	var err error
	eg, _ := errgroup.WithContext(ctx)
	for i := range params.Len() {
		eg.Go(func() error {
			var out string
			var hist []openai.ChatCompletionMessage
			out, hist, err = ChatTeacher(
				ctx,
				model,
				params.Histories[i],
				params.Prompts[i],
			)
			if err != nil {
				return fmt.Errorf("failed to chat with llm provider: %w", err)
			}
			histories[i] = hist
			outs = append(outs, out)
			return nil
		})
	}
	return outs, histories, nil
}

// ChatNonce is a generic number method to get prompting help from a teacher who sees the
// whole conversation and is able to explain the prompt or error.
func ChatNonce(
	ctx context.Context,
	model string,
	histories [][]openai.ChatCompletionMessage,
	prompt string,
) ([]string, [][]openai.ChatCompletionMessage, error) {
	var outs []string
	var err error
	// make an array of prompts with the same length as the histories
	prompts := make([]string, len(histories))
	for i := range histories {
		prompts[i] = prompt
	}
	// call ChatN with the same length as the histories
	outs, histories, err = ChatN(ctx, model, ChatNParams{Histories: histories, Prompts: prompts})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to chat with llm provider: %w", err)
	}
	return outs, histories, nil
}
