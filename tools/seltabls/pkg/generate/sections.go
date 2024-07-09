package generate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Field is a struct for a field
type Field struct {
	// Name is the name of the field.
	Name string `json:"name"`
	// Type is the type of the field.
	Type string `json:"type"`
	// Description is a description of the field.
	Description string `json:"description"`
	// HeaderSelector is the header selector for the field.
	HeaderSelector string `json:"header-selector"`
	// DataSelector is the data selector for the field.
	DataSelector string `json:"data-selector"`
	// ControlSelector is the control selector for the field.
	ControlSelector string `json:"control-selector"`
	// QuerySelector is the query selector for the field.
	QuerySelector string `json:"query-selector"`
	// MustBePresent is the must be present selector for the field.
	MustBePresent string `json:"must-be-present"`
}

// aggregateSections aggregates sections for the struct file.
func aggregateSections(
	ctx context.Context,
	s *StructFile,
	client *openai.Client,
	model string,
	sectCh chan string,
) (sec *Section, err error) {
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
		generation, _, err := Chat(
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
		*sec, err = decodeSection(generation)
		if err != nil {
			return nil, fmt.Errorf("failed to decode section: %w", err)
		}
	}
	return sec, nil
}

// decodeSection decodes a section from a string
func decodeSection(s string) (Section, error) {
	var section Section
	err := json.Unmarshal([]byte(s), &section)
	if err != nil {
		return Section{}, fmt.Errorf("failed to unmarshal section: %w", err)
	}
	return section, nil
}
