// Package identify provides the identify functionality for the seltabl package.
package identify

import (
	"context"
	"encoding/json"
	"fmt"

	// Embedded for the identify template
	_ "embed"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/sync/errgroup"
)

// decodeIdentify decodes the identify response.
func decodeIdentify(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	out string,
) (result domain.IdentifyResponse, err error) {
	log.Debugf("DecodeIdentify called with out: %s", out)
	defer log.Debugf("DecodeIdentify called with out: %s", out)
	var id domain.IdentifyResponse
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
	generation, history, err := domain.Chat(
		ctx,
		client,
		model,
		history,
		prompt,
	)
	if err != nil {
		return id, fmt.Errorf("failed to chat with llm provider: %w", err)
	}
	return decodeIdentify(ctx, client, model, history, generation)
}

// generateIdentity generates the identity for the struct file.
func generateIdentity(
	ctx context.Context,
	s *domain.StructFile,
	client *openai.Client,
) (identity domain.IdentifyResponse, err error) {
	log.Debugf("generateIdentity called with s: %v", s)
	defer log.Debugf("generateIdentity called with s: %v", s)
	eg := errgroup.Group{}
	outCh := make(chan string)
	for range make([]int, s.TreeWidth) {
		eg.Go(func() error {
			var identifyPrompt string
			identifyPrompt, err = NewIdentifyPrompt(
				s.URL,
				s.HTMLContent,
			)
			if err != nil {
				return fmt.Errorf("failed to create identify prompt: %w", err)
			}
			out, _, err := domain.Chat(
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
	s *domain.StructFile,
	client *openai.Client,
	outCh chan string,
) (identity domain.IdentifyResponse, err error) {
	log.Debugf("aggregateIdentity called with outCh: %v", outCh)
	defer log.Debugf("aggregateIdentity called with outCh: %v", outCh)
	eg := errgroup.Group{}
	for range make([]int, s.TreeDepth) {
		eg.Go(func() error {
			channelLength := len(outCh)
			if channelLength >= 3 {
				aggPrompt, err := NewIdentifyAggregatePrompt(
					s.HTMLContent,
					[]string{<-outCh, <-outCh, <-outCh},
				)
				if err != nil {
					return fmt.Errorf(
						"failed to create struct prompt: %w",
						err,
					)
				}
				out, history, err := domain.Chat(
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
				identity, err = decodeIdentify(
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
