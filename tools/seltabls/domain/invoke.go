package domain

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

// InvokeJSONSimple is a function for generating json using the OpenAI API with the given prompt.
//
// It does not decode the json.
func InvokeJSONSimple(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
) (out string, postHistory []openai.ChatCompletionMessage, err error) {
	for {
		select {
		case <-ctx.Done():
			return "", postHistory, ctx.Err()
		default:
			prmpt, err := NewPrompt(prompt)
			if err != nil {
				return "", history, err
			}
			genHistory := append(
				history,
				openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: prmpt,
				})
			completion, err := client.CreateChatCompletion(
				ctx, openai.ChatCompletionRequest{
					Model:    model,
					Messages: genHistory,
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: "json_object",
					},
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
			genHistory = append(genHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completion.Choices[0].Message.Content})
			if len(completion.Choices) == 0 {
				return "", genHistory, fmt.Errorf("no choices found")
			}
			if len(completion.Choices[0].Message.Content) == 0 {
				return "", genHistory, fmt.Errorf("no content found")
			}
			return completion.Choices[0].Message.Content, genHistory, nil
		}
	}
}

// InvokeJSON is a function for generating json using the OpenAI API with the given prompt.
func InvokeJSON(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
	output interface{},
	htmlBody string,
) (out string, postHistory []openai.ChatCompletionMessage, err error) {
	for {
		select {
		case <-ctx.Done():
			return "", postHistory, ctx.Err()
		default:
			prmpt, err := NewPrompt(prompt)
			if err != nil {
				return "", history, err
			}
			genHistory := append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: prmpt,
			})
			completion, err := client.CreateChatCompletion(
				ctx,
				openai.ChatCompletionRequest{
					Model:    model,
					Messages: genHistory,
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: "json_object",
					},
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
				openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: completion.Choices[0].Message.Content,
				})
			err = DecodeJSON(
				ctx,
				[]byte(completion.Choices[0].Message.Content),
				output,
				genHistory,
				client,
				model,
				htmlBody,
			)
			if err != nil {
				return "", genHistory, err
			}
			if len(completion.Choices) == 0 {
				return "", genHistory, fmt.Errorf("no choices found")
			}
			if len(completion.Choices[0].Message.Content) == 0 {
				return "", genHistory, fmt.Errorf("no content found")
			}
			return completion.Choices[0].Message.Content, genHistory, nil
		}
	}
}

// InvokeJSONN is a function for generating json using the OpenAI API multiple "N" times.
func InvokeJSONN(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
	output Responder,
	n int,
	htmlBody string,
) (outs []string, histories [][]openai.ChatCompletionMessage, err error) {
	outs = make([]string, n)
	histories = make([][]openai.ChatCompletionMessage, n)
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
				history = append(history, openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: out})
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
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
	n int,
) (outs []string, histories [][]openai.ChatCompletionMessage, err error) {
	outs = make([]string, n)
	histories = make([][]openai.ChatCompletionMessage, n)
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
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
) (out string, postHistory []openai.ChatCompletionMessage, err error) {
	for {
		select {
		case <-ctx.Done():
			return "", postHistory, ctx.Err()
		default:
			prmpt, err := NewPrompt(prompt)
			if err != nil {
				return "", history, err
			}
			completion, err := client.CreateChatCompletion(
				ctx, openai.ChatCompletionRequest{
					Model: model,
					Messages: append(history, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleSystem,
						Content: prmpt,
					}),
				},
			)
			if err != nil {
				return "", history, err
			}
			history = append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completion.Choices[0].Message.Content})
			log.Debugf(
				"completion: %+v",
				completion.Choices[0].Message.Content,
			)
			return completion.Choices[0].Message.Content, history, nil
		}
	}
}

// InvokePreTxt is a function for generating text using the OpenAI API by
// prepending the prompt to the history.
func InvokePreTxt(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
) (out string, postHistory []openai.ChatCompletionMessage, err error) {
	for {
		select {
		case <-ctx.Done():
			return "", postHistory, ctx.Err()
		default:
			prmpt, err := NewPrompt(prompt)
			if err != nil {
				return "", history, err
			}
			completion, err := client.CreateChatCompletion(
				ctx, openai.ChatCompletionRequest{
					Model: model,
					Messages: append(
						[]openai.ChatCompletionMessage{{
							Role:    openai.ChatMessageRoleSystem,
							Content: prmpt,
						}},
						history...,
					),
				},
			)
			if err != nil {
				return "", history, err
			}
			history = append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completion.Choices[0].Message.Content,
			})
			return completion.Choices[0].Message.Content, history, nil
		}
	}
}
