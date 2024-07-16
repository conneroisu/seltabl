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
			genHistory := append(
				history,
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
					Messages:  genHistory,
					MaxTokens: 1000,
				},
			)
			if err != nil {
				var e *anthropic.APIError
				if errors.As(err, &e) {
					fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
				} else {
					fmt.Printf("Messages error: %v\n", err)
				}
				return "", genHistory, err
			}
			genHistory = append(
				genHistory,
				anthropic.Message{
					Role:    openai.ChatMessageRoleAssistant,
					Content: generation.Content,
				})
			if len(*generation.Content[0].Text) == 0 {
				return "", genHistory, fmt.Errorf("no choices found")
			}
			if len(*generation.Content[0].Text) == 0 {
				return "", genHistory, fmt.Errorf("no content found")
			}
			return *generation.Content[0].Text, genHistory, nil
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
					MaxTokens: 1000,
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

// InvokeJSONN is a function for generating json using the OpenAI API multiple "N" times.
func InvokeJSONN(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	history []anthropic.Message,
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
				out, hist, err := InvokeJSON(
					hCtx,
					client,
					model,
					history,
					prompt,
					output,
					htmlBody,
				)
				if err != nil {
					return err
				}
				outs[idx] = out
				histories[idx] = hist
				return nil
			})
		}
	}
	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}
	return outs, histories, nil
}

// InvokeTxtN is a function for generating text using the OpenAI API multiple "N" times.
func InvokeTxtN(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	history []anthropic.Message,
	prompt prompter,
	n int,
) (outs []string, histories [][]anthropic.Message, err error) {
	outs = make([]string, n)
	histories = make([][]anthropic.Message, n)
	for idx := 0; idx < n; idx++ {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			out, hist, err := InvokeTxt(
				ctx,
				client,
				model,
				history,
				prompt,
			)
			if err != nil {
				return nil, nil, err
			}
			outs[idx] = out
			histories[idx] = hist
			if len(outs) >= n {
				return outs, histories, nil
			}
		}
	}
	return outs, histories, nil
}

// InvokeTxt is a function for generating text using the OpenAI API.
func InvokeTxt(
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
			prmpt, err := NewPrompt(prompt)
			if err != nil {
				return "", history, err
			}
			completion, err := client.CreateMessages(
				ctx, anthropic.MessagesRequest{
					Model:     model,
					MaxTokens: 1000,
					Messages: append(history, anthropic.Message{
						Role:    openai.ChatMessageRoleUser,
						Content: []anthropic.MessageContent{anthropic.NewTextMessageContent(prmpt)},
					}),
				},
			)
			if err != nil {
				return "", history, err
			}
			history = append(
				history,
				anthropic.Message{
					Role:    openai.ChatMessageRoleAssistant,
					Content: completion.Content,
				})
			return *completion.Content[0].Text, history, nil
		}
	}
}

// InvokePreTxt is a function for generating text using the OpenAI API by
// prepending the prompt to the history.
func InvokePreTxt(
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
			prmpt, err := NewPrompt(prompt)
			if err != nil {
				return "", history, err
			}
			generation, err := client.CreateMessages(
				ctx, anthropic.MessagesRequest{
					Model:     model,
					MaxTokens: 1000,
					Messages: append(
						[]anthropic.Message{{
							Role:    openai.ChatMessageRoleUser,
							Content: []anthropic.MessageContent{anthropic.NewTextMessageContent(prmpt)},
						}},
						history...,
					),
				},
			)
			if err != nil {
				return "", history, err
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
