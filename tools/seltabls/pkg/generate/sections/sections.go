package sections

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/sashabaranov/go-openai"
)

// aggregateSections aggregates sections for the struct file.
func aggregateSections(
	ctx context.Context,
	s *domain.StructFile,
	client *openai.Client,
	model string,
	sectCh chan string,
) (sec *domain.Section, err error) {
	log.Debugf("aggregateSections called with s: %v", s)
	defer log.Debugf("aggregateSections returned with sec: %v", sec)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		aggregatePrompt, err := NewAggregatePrompt(
			s.URL,
			s.HTMLContent,
			s.ConfigFile.Selectors,
			[]string{string(<-sectCh)},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create struct prompt: %w", err)
		}
		generation, _, err := domain.Chat(
			ctx,
			client,
			model,
			[]openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: aggregatePrompt,
				}},
			aggregatePrompt,
		)
		sc, err := decodeSection(generation)
		if err != nil {
			return nil, fmt.Errorf("failed to decode section: %w", err)
		}
		sec = &sc
	}
	return sec, nil
}

// decodeSection decodes a section from a string
func decodeSection(s string) (domain.Section, error) {
	log.Debugf("decodeSection called with s: %s", s)
	var section domain.Section
	err := json.Unmarshal([]byte(s), &section)
	if err != nil {
		return domain.Section{}, fmt.Errorf("failed to unmarshal section: %w", err)
	}
	return section, nil
}
