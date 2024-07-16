package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

// NewClient creates a new client for the given api key.
func NewClient(apiKey string) *anthropic.Client {
	return anthropic.NewClient(apiKey)
}

// InvokeJSONSimple is a function for generating json using the OpenAI API with the given prompt.
//
// It does not decode the json.
func InvokeJSONSimple(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	history []anthropic.Message,
	prompt prompter,
) (out string, postHistory []anthropic.Message, err error) {
	hCtx, cancel := context.WithTimeout(ctx, time.Second*12)
	hHistory := make([]anthropic.Message, len(history))
	copy(hHistory, history)
	defer cancel()
	for {
		select {
		case <-hCtx.Done():
			return "", postHistory, hCtx.Err()
		default:
			prmpt, err := NewPrompt(prompt)
			if err != nil {
				return "", history, err
			}
			hHistory := append(
				hHistory,
				anthropic.Message{
					Role: openai.ChatMessageRoleUser,
					Content: []anthropic.MessageContent{
						anthropic.NewTextMessageContent(prmpt),
					},
				})
			generation, err := client.CreateMessages(
				hCtx,
				anthropic.MessagesRequest{
					Model:     model,
					Messages:  hHistory,
					MaxTokens: 10_000,
				},
			)
			if err != nil {
				var e *anthropic.APIError
				if errors.As(err, &e) {
					fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
					return "", hHistory, fmt.Errorf("messages error, type: %s, message: %s", e.Type, e.Message)
				}
				fmt.Printf("Messages error: %v\n", err)
				return "", hHistory, fmt.Errorf("messages error: %v", err)
			}
			hHistory = append(
				hHistory,
				anthropic.Message{
					Role:    openai.ChatMessageRoleAssistant,
					Content: generation.Content,
				})
			if len(*generation.Content[0].Text) == 0 {
				return "", hHistory, fmt.Errorf("no choices found")
			}
			if len(*generation.Content[0].Text) == 0 {
				return "", hHistory, fmt.Errorf("no content found")
			}
			return *generation.Content[0].Text, hHistory, nil
		}
	}
}

// InvokeJSON is a function for generating json using the OpenAI API with the given prompt.
func InvokeJSON(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	history []anthropic.Message,
	prompt prompter,
	output interface{},
	htmlBody string,
) (out string, postHistory []anthropic.Message, err error) {
	for {
		select {
		case <-ctx.Done():
			return "", postHistory, ctx.Err()
		default:
			prmpt, err := NewPrompt(prompt)
			if err != nil {
				return "", history, err
			}
			genHistory := append(history, anthropic.Message{
				Role:    openai.ChatMessageRoleUser,
				Content: []anthropic.MessageContent{anthropic.NewTextMessageContent(prmpt)},
			})
			completion, err := client.CreateMessages(
				ctx,
				anthropic.MessagesRequest{
					Model:     model,
					Messages:  genHistory,
					MaxTokens: 10_000,
				},
			)
			if err != nil {
				var ok bool
				_, ok = err.(*openai.RequestError)
				if ok {
					log.Debugf("request error: %+v", err)
				}
				_, ok = err.(*openai.APIError)
				if ok {
					log.Debugf("rate limit error: %+v", err)
				}
				return "", genHistory, err
			}
			genHistory = append(
				genHistory,
				anthropic.Message{
					Role:    openai.ChatMessageRoleAssistant,
					Content: completion.Content,
				})
			err = DecodeJSON(
				ctx,
				[]byte(*completion.Content[0].Text),
				output,
				genHistory,
				client,
				model,
				htmlBody,
			)
			if err != nil {
				return "", genHistory, err
			}
			if len(completion.Content) == 0 {
				return "", genHistory, fmt.Errorf("no choices found")
			}
			if len(*completion.Content[0].Text) == 0 {
				return "", genHistory, fmt.Errorf("no content found")
			}
			return *completion.Content[0].Text, genHistory, nil
		}
	}
}

// InvokeN is a function for generating json using the OpenAI API multiple "N" times.
func InvokeN(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	_ [][]anthropic.Message,
	prompt prompter,
	output responder,
	n int,
	htmlBody string,
) (outs []string, histories [][]anthropic.Message, err error) {
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
				var out string
				out, histories[idx], err = InvokeJSON(
					hCtx,
					client,
					model,
					histories[idx],
					prompt,
					output,
					htmlBody,
				)
				if err != nil {
					return err
				}
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

// InvokePre is a function for generating text using the OpenAI API by
// prepending the prompt to the history.
func InvokePre(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	history []anthropic.Message,
	prompt prompter,
) (out string, postHistory []anthropic.Message, err error) {
	for {
		select {
		case <-ctx.Done():
			return ``, postHistory, ctx.Err()
		default:
			prmpt, err := NewPrompt(prompt)
			if err != nil {
				return ``, history, err
			}
			generation, err := client.CreateMessages(
				ctx, anthropic.MessagesRequest{
					Model:     model,
					MaxTokens: 10_000,
					Messages: append(
						[]anthropic.Message{{
							Role:    openai.ChatMessageRoleUser,
							Content: []anthropic.MessageContent{anthropic.NewTextMessageContent(prmpt)},
						},
							{
								Role:    openai.ChatMessageRoleAssistant,
								Content: []anthropic.MessageContent{anthropic.NewTextMessageContent(`Okay, I am ready for the history of the conversation.`)},
							}},
						history...,
					),
				},
			)
			if err != nil {
				return ``, history, err
			}
			history = append(
				history,
				anthropic.Message{
					Role:    openai.ChatMessageRoleAssistant,
					Content: generation.Content,
				})
			return *generation.Content[0].Text, history, nil
		}
	}
}
