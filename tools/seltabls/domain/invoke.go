package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

// NewClient creates a new client for the given api key.
func NewClient(apiKey string) *anthropic.Client {
	return anthropic.NewClient(apiKey)
}

// Invoke is a function for generating json using the OpenAI API with the given
// prompt.
func Invoke(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	history []anthropic.Message,
	prompt prompter,
) (out string, postHistory []anthropic.Message, err error) {
	for {
		select {
		case <-ctx.Done():
			return "", postHistory, ctx.Err()
		default:
			return invoke(ctx, client, model, history, prompt)
		}
	}
}

// invoke is a function for generating text using the OpenAI API by
func invoke(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	history []anthropic.Message,
	prompt prompter,
) (out string, postHistory []anthropic.Message, err error) {
	prmpt, err := NewPrompt(prompt)
	if err != nil {
		return "", history, err
	}
	genHistory := append(history, anthropic.Message{
		Role: anthropic.RoleUser,
		Content: []anthropic.MessageContent{
			anthropic.NewTextMessageContent(prmpt),
		},
	})
	completion, err := client.CreateMessages(
		ctx,
		anthropic.MessagesRequest{
			Model:     model,
			Messages:  genHistory,
			MaxTokens: 4_096,
		},
	)
	if err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf(
				"Messages error, type: %s, message: %s",
				e.Type,
				e.Message,
			)
			return "", history, fmt.Errorf(
				"messages error, type: %s, message: %s",
				e.Type,
				e.Message,
			)
		}
		log.Errorf("Messages error: %v", err)
		return "", genHistory, err
	}
	genHistory = append(
		genHistory,
		anthropic.Message{
			Role:    anthropic.RoleUser,
			Content: completion.Content,
		})
	if len(completion.Content) == 0 {
		return "", genHistory, fmt.Errorf("no choices found")
	}
	if len(*completion.Content[0].Text) == 0 {
		return "", genHistory, fmt.Errorf("no content found")
	}
	return *completion.Content[0].Text, genHistory, nil
}

// InvokeN is a function for generating json using the OpenAI API multiple "N" times.
func InvokeN(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	prompt prompter,
	n int,
) (outs []string, histories [][]anthropic.Message, err error) {
	log.Debugf("Current invocation: %s", prompt.prompt())
	outs = make([]string, n)
	histories = make([][]anthropic.Message, n)
	var eg *errgroup.Group
	var hCtx context.Context
	eg, hCtx = errgroup.WithContext(ctx)
	for idx := 0; idx < n; idx++ {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			eg.Go(func() error {
				out, hist, err := Invoke(
					hCtx,
					client,
					model,
					histories[idx],
					prompt,
				)
				if err != nil {
					return err
				}
				log.Debugf("History: %v", hist)
				histories[idx] = hist
				outs[idx] = out
				return nil
			})
		}
	}
	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}
	return outs, histories, nil
}

// ChatUnmarshal ensures the given data is a valid JSON object.
func ChatUnmarshal(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	data []byte,
	v interface{},
) error {
	tryID := func() error {
		errMsg, _, err := Invoke(
			ctx,
			client,
			model,
			[]anthropic.Message{},
			IdentifyErrorArgs{
				History: history,
				Error:   err,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to generate identify completion: %w", err)
		}
		history = append(
			history,
			anthropic.Message{
				Role:    openai.ChatMessageRoleUser,
				Content: []anthropic.MessageContent{anthropic.NewTextMessageContent(errMsg)},
			})
		identifyCompletion, history, err = domain.Invoke(
			ctx,
			client,
			params.FastModel,
			history,
			domain.IdentifyAggregateArgs{
				Schemas:   identifyCompletions,
				Content:   string(htmlBody),
				Selectors: selectors,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to generate identify completion: %w", err)
		}
		ext, err := extractJSON(identifyCompletion)
		if err != nil {
			return fmt.Errorf("failed to extract JSON from identifyCompletion: %w", err)
		}
		err = json.Unmarshal([]byte(ext), &identified)
		if err != nil {
			return fmt.Errorf("failed to unmarshal identifyCompletion: %w", err)
		}
		return err
	}
	retryLimit := 3
	for i := 0; i < retryLimit; i++ {
		if err := tryID(); err == nil {
			break
		}
		log.Debugf("Failed to extract JSON from identifyCompletion: %s", identifyCompletion)
		log.Debugf("Retrying")
	}
}
