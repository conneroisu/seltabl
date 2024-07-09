package generate

import (
	"context"
	"encoding/json"
	"fmt"

	// Embedded for the identify template
	_ "embed"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

// IdentifyResponse is a struct for the respond of an identify prompt.
//
// The identify prompt is used to describe the structure of a given
// html returning this struct in the form of json.
type IdentifyResponse struct {
	// Sections is a list of sections in the html.
	Sections []Section `json:"sections" yaml:"sections"`
}

// Section is a struct for a section in the html.
type Section struct {
	// Name is the name of the section.
	Name string `json:"name"        yaml:"name"`
	// Description is a description of the section.
	Description string `json:"description" yaml:"description"`
	// CSS is the css selector for the section.
	CSS string `json:"css"         yaml:"css"`
	// Start is the start of the section in the html.
	Start int `json:"start"       yaml:"start"`
	// End is the end of the section in the html.
	End int `json:"end"         yaml:"end"`
}

// DecodeIdentify decodes the identify response.
func DecodeIdentify(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	out string,
) (result IdentifyResponse, err error) {
	var id IdentifyResponse
	err = json.Unmarshal([]byte(out), &id)
	if err == nil {
		return id, nil
	}
	prompt, err := NewIdentifyErrorPrompt(err)
	if err != nil {
		return id, fmt.Errorf(
			"failed to create identify error prompt: %w",
			err,
		)
	}
	history = append(history, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	})
	generation, history, err := Chat(
		ctx,
		client,
		model,
		history,
		prompt,
	)
	if err != nil {
		return id, fmt.Errorf("failed to chat with llm provider: %w", err)
	}
	return DecodeIdentify(ctx, client, model, history, generation)
}

// generateIdentity generates the identity for the struct file.
func generateIdentity(
	ctx context.Context,
	s *StructFile,
	client *openai.Client,
) (identity IdentifyResponse, err error) {
	eg := errgroup.Group{}
	outCh := make(chan string)
	for range make([]int, s.TreeWidth) {
		eg.Go(func() error {
			identifyPrompt, err := NewIdentifyPrompt(
				s.URL,
				s.HTMLContent,
			)
			out, _, err := Chat(
				ctx,
				client,
				s.ConfigFile.FastModel,
				[]openai.ChatCompletionMessage{
					{
						Role:    openai.ChatMessageRoleUser,
						Content: identifyPrompt,
					},
				},
				identifyPrompt,
			)
			if err != nil {
				return fmt.Errorf("failed to create identify prompt: %w", err)
			}
			outCh <- out
			return nil
		})
		eg.Go(func() error {
			identity, err = aggregateIdentity(ctx, s, client, outCh)
			if err != nil {
				return fmt.Errorf("failed to aggregate identity: %w", err)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return identity, fmt.Errorf("failed to generate identity: %w", err)
	}
	return identity, nil
}

// aggregateIdentity aggregates the identify response.
func aggregateIdentity(
	ctx context.Context,
	s *StructFile,
	client *openai.Client,
	outCh chan string,
) (identity IdentifyResponse, err error) {
	eg := errgroup.Group{}
	for range make([]int, s.TreeDepth) {
		eg.Go(func() error {
			channelLength := len(outCh)
			if channelLength >= 3 {
				aggPrompt, err := NewAggregatePrompt(
					s.URL,
					s.HTMLContent,
					s.ConfigFile.Selectors,
					[]string{<-outCh, <-outCh, <-outCh},
				)
				if err != nil {
					return fmt.Errorf(
						"failed to create struct prompt: %w",
						err,
					)
				}
				out, history, err := Chat(
					ctx,
					client,
					s.ConfigFile.SmartModel,
					[]openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: aggPrompt,
						},
					},
					aggPrompt,
				)
				if err != nil {
					return fmt.Errorf(
						"failed to chat with llm provider: %w",
						err,
					)
				}
				identity, err = DecodeIdentify(
					ctx,
					client,
					s.ConfigFile.FastModel,
					history,
					out,
				)
				if err != nil {
					return fmt.Errorf("failed to decode identify: %w", err)
				}
				return nil
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return identity, fmt.Errorf("failed to aggregate identity: %w", err)
	}
	return identity, nil
}
