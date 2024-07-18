package domain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/ollama/ollama/api"
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
	hCtx, cancel := context.WithTimeout(ctx, time.Second*90)
	defer cancel()
	for {
		select {
		case <-hCtx.Done():
			return "", history, hCtx.Err()
		default:
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
				hCtx,
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
					if e.Type == anthropic.ErrTypeRateLimit {
						log.Debugf("Rate limit error: %s", e.Message)
						time.Sleep(time.Second * 45)
						return invoke(hCtx, client, model, history, prompt)
					}
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
	}
}

// InvokeResponse is a struct for the Invoke response.
type InvokeResponse struct {
	Out  string
	Hist []anthropic.Message
}

// InvokeNResponse is a struct for the InvokeN response.
type InvokeNResponse struct {
	Outs []string
	Hist [][]anthropic.Message
}

// InvokeN is a function for generating json using the OpenAI API multiple "N" times.
func InvokeN(
	ctx context.Context,
	client *anthropic.Client,
	model string,
	prompt prompter,
	n int,
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

var stream = false

// ChatUnmarshal ensures the given data is a valid JSON object
func ChatUnmarshal(
	ctx context.Context,
	client *api.Client,
	data []byte,
	v interface{},
) error {
	var err error
	var prmpt string
	var hCtx context.Context
	hCtx, cancel := context.WithTimeout(ctx, time.Second*70)
	defer cancel()
	for {
		select {
		case <-hCtx.Done():
			return hCtx.Err()
		default:
			err = json.Unmarshal(data, &v)
			if err == nil {
				return nil
			}
			prmpt, err = NewPrompt(FixJSONArgs{JSON: string(data)})
			if err != nil {
				return fmt.Errorf("failed to create fix json prompt: %w", err)
			}
			respFunc := func(resp api.GenerateResponse) error {
				err = json.Unmarshal([]byte(resp.Response), &v)
				if err == nil {
					log.Debugf("unmarshaled: %s", resp.Response)
					return nil
				}
				log.Debugf("retrying failed to unmarshal: %s\n response: %s", err, resp.Response)
				ChatUnmarshal(hCtx, client, data, v)
				return nil
			}
			err = client.Generate(hCtx, &api.GenerateRequest{
				Format: "json",
				Model:  "llama3",
				Stream: &stream,
				Prompt: prmpt,
			}, respFunc)
			if err != nil {
				return fmt.Errorf("failed to generate fix json: %w", err)
			}
			return nil
		}
	}
}

// InvokeC invokes the client.
func InvokeC(
	ctx context.Context,
	client *api.Client,
	model string,
	historyInput []api.Message,
	prompt prompter,
	v interface{},
) (string, []api.Message, error) {
	var out string
	var err error
	hCtx, cancel := context.WithTimeout(ctx, time.Second*120)
	defer cancel()
	for {
		select {
		case <-hCtx.Done():
			return "", nil, hCtx.Err()
		default:
			var invokationPrompt string
			invokationPrompt, err = NewPrompt(prompt)
			if err != nil {
				return "", nil, fmt.Errorf(
					"failed to create (%s) prompt: %w",
					prompt.prompt(),
					err,
				)
			}
			historyInput = append(historyInput, api.Message{
				Role:    "system",
				Content: invokationPrompt,
			})
			var respFunc api.GenerateResponseFunc
			respFunc = func(resp api.GenerateResponse) error {
				if v == nil {
					out = resp.Response
				}
				err = ChatUnmarshal(
					hCtx,
					client,
					[]byte(resp.Response),
					v,
				)
				if err == nil {
					return nil
				}
				return fmt.Errorf(
					"failed to unmarshal response (%s): %w",
					resp.Response,
					err,
				)
			}
			err = client.Generate(
				hCtx,
				&api.GenerateRequest{
					Format: "json",
					Stream: &stream,
					Model:  model,
					Prompt: invokationPrompt,
				},
				respFunc,
			)
			historyInput = append(historyInput, api.Message{
				Role:    "assistant",
				Content: out,
			})
			return out, historyInput, err
		}
	}
}
