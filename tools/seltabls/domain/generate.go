package domain

import (
	"errors"
	"io"
	"os"

	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"

	"github.com/sashabaranov/go-openai"
)

// ConfigFile is a struct for a config file.
type ConfigFile struct {
	// Name is the name of the config file.
	Name string `yaml:"name"`
	// Description is the description of the config file.
	Description *string `yaml:"description,omitempty"`
	// URL is the url for the config file.
	URL string `yaml:"url"`
	// IgnoreElements is a list of elements to ignore when generating the
	// struct.
	IgnoreElements []string `yaml:"ignore-elements"`
	// Selectors is a list of selectors for the config file.
	Selectors []master.Selector `yaml:"selectors"`
	// HTMLBody is the html body for the config file.
	HTMLBody string `yaml:"html-body"`
	// NumberedHTMLBody is the numbered html body for the config file.
	NumberedHTMLBody string `yaml:"-"`
	// SmartModel is the model for the config file.
	SmartModel string `yaml:"model"`
	// FastModel is the model for the config file.
	FastModel string `yaml:"fast-model"`

	// Sections is a list of sections in the html.
	Sections []Section `json:"sections" yaml:"sections"`
}

// IdentifyResponse is a struct for the respond of an identify prompt.
//
// The identify prompt is used to describe the structure of a given
// html returning this struct in the form of json.
type IdentifyResponse struct {
	// Sections is a list of sections in the html.
	Sections []Section `json:"sections" yaml:"sections"`
	// PackageName is the package name for the identify response.
	PackageName string `json:"package-name" yaml:"package-name"`
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
	// Fields is a list of fields in the section.
	Fields []Field `json:"fields"      yaml:"fields"`
}

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

// TestFile is a struct for a test file
type TestFile struct {
	// Name is the name of the test file
	Name string `json:"name" yaml:"name"`
	// URL is the url for the test file
	URL string `json:"url"  yaml:"url"`
	// PackageName is the package name for the test file
	PackageName string `json:"-"    yaml:"package-name"`
}

// WriteFile writes the test file to the file system
func (t *TestFile) WriteFile(p []byte) (n int, err error) {
	err = os.WriteFile(
		fmt.Sprintf("%s_test.go", t.Name),
		p,
		0644,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to write test file: %w", err)
	}
	return len(p), nil
}

// StructFile is a struct for a struct file.
//
// It contains attributes relating to the name, url, and ignore elements of the
// struct file.
type StructFile struct {
	File os.File `json:"-"               yaml:"-"`
	// Name is the name of the struct file.
	Name string `json:"-"               yaml:"name"`
	// URL is the url for the struct file.
	URL string `json:"-"               yaml:"url"`
	// IgnoreElements is a list of elements to ignore when generating the
	// struct.
	IgnoreElements []string `json:"ignore-elements" yaml:"ignore-elements"`
	// Fields is a list of fields for the struct.
	Fields []Field `json:"fields"          yaml:"fields"`

	// TreeWidth is the width of the tree when generating the struct.
	TreeWidth int `json:"-" yaml:"tree-width"`
	// TreeDepth is the depth of the tree when generating the struct.
	TreeDepth int `json:"-" yaml:"tree-depth"`

	// ConfigFile is the config file for the struct file.
	ConfigFile ConfigFile `json:"-" yaml:"config-file"`
	// JSONValue is the json value for the struct yaml file.
	JSONValue string `json:"-" yaml:"json-value"`
	// HTMLContent is the html content for the struct file.
	HTMLContent string `json:"-" yaml:"html-content"`

	// Db is the database for the struct file.
	Db *data.Database[master.Queries] `json:"-" yaml:"-"`

	// Section is the section of the struct file.
	Section Section `json:"-" yaml:"section"`
}

// Chat is a struct for a chat by appending the prompt to the history.
func Chat(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt string,
) (out string, postHistory []openai.ChatCompletionMessage, err error) {
	log.Debugf("calling Chat in domain")
	defer log.Debugf("Chat completed in domain")
	for {
		select {
		case <-ctx.Done():
			return "", postHistory, ctx.Err()
		default:
			stream, err := client.CreateChatCompletionStream(
				ctx, openai.ChatCompletionRequest{
					Model: model,
					Messages: append(history, openai.ChatCompletionMessage{
						Role:    openai.ChatMessageRoleSystem,
						Content: prompt,
					}),
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: "json",
					},
					StreamOptions: &openai.StreamOptions{},
				},
			)
			defer stream.Close()

			if err != nil {
				return "", history, err
			}
			completion := ""
			for {
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					return "", history, fmt.Errorf("failed to receive streamed response: %w", err)
				}
				log.Debugf("response: %+v", response)
				completion += response.Choices[0].Delta.Content
			}
			history = append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completion})
			return completion, history, nil
		}
	}
}

// ChatPre is a function for chatting with the model by prepending the prompt to the history.
func ChatPre(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt string,
) (out string, postHistory []openai.ChatCompletionMessage, err error) {
	log.Debugf("calling ChatPre in domain")
	defer log.Debugf("ChatPre completed in domain")
	for {
		select {
		case <-ctx.Done():
			return "", postHistory, ctx.Err()
		default:
			stream, err := client.CreateChatCompletionStream(
				ctx, openai.ChatCompletionRequest{
					Model: model,
					Messages: append(
						[]openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: prompt}},
						history...,
					),
				},
			)
			if err != nil {
				return "", history, err
			}
			completion := ""
			for {
				response, err := stream.Recv()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					return "", history, fmt.Errorf("failed to receive streamed response: %w", err)
				}
				log.Debugf("response: %+v", response)
				completion += response.Choices[0].Delta.Content
			}
			log.Debugf("completion: %+v", completion)
			history = append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completion,
			})
			return completion, history, nil
		}
	}
}

// ChatN is a function for chatting with the model multiple "N" times.
func ChatN(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt string,
	n int,
) (outs []string, postHistories [][]openai.ChatCompletionMessage, err error) {
	log.Debugf("calling ChatN in domain")
	defer log.Debugf("ChatN completed in domain")
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			out, hist, err := Chat(
				ctx,
				client,
				model,
				history,
				prompt,
			)
			if err != nil {
				return nil, nil, err
			}
			outs = append(outs, out)
			postHistories = append(postHistories, hist)
			if len(outs) >= n {
				return outs, postHistories, nil
			}
			history = append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: out})
			return outs, postHistories, nil
		}
	}
}

// ChatPreN is a function for chatting with the model by prepending the prompt
// to the history multiple "N" times.
func ChatPreN(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt string,
	n int,
) (outs []string, postHistories [][]openai.ChatCompletionMessage, err error) {
	log.Debugf("calling ChatPreN in domain")
	defer log.Debugf("ChatPreN completed in domain")
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			out, hist, err := ChatPre(
				ctx,
				client,
				model,
				history,
				prompt,
			)
			if err != nil {
				return nil, nil, err
			}
			outs = append(outs, out)
			postHistories = append(postHistories, hist)
			if len(outs) >= n {
				return outs, postHistories, nil
			}
			history = append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: out})
			return outs, postHistories, nil
		}
	}
}

// CreateClient creates a new client for the given api key.
func CreateClient(baseURL string, apiKey string) *openai.Client {
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = baseURL
	cfg.APIVersion = string(openai.APITypeOpenAI)
	cfg.APIType = openai.APITypeOpenAI
	return openai.NewClientWithConfig(cfg)
}

// Generate is a function for generating text using the OpenAI API.
func Generate(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
) (out string, postHistory []openai.ChatCompletionMessage, err error) {
	log.Debugf("calling Generate in domain")
	defer log.Debugf("Generate completed in domain")
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
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: "json_object",
					},
				},
			)
			log.Debugf("completion: %+v", completion)
			if err != nil {
				return "", history, err
			}
			history = append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completion.Choices[0].Message.Content})
			return completion.Choices[0].Message.Content, history, nil
		}
	}
}

// GenerateTxt is a function for generating text using the OpenAI API.
func GenerateTxt(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
) (out string, postHistory []openai.ChatCompletionMessage, err error) {
	log.Debugf("calling GenerateTxt in domain")
	defer log.Debugf("GenerateTxt completed in domain")
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
			return completion.Choices[0].Message.Content, history, nil
		}
	}
}

// GenerateN is a function for generating text using the OpenAI API multiple "N" times.
func GenerateN(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
	n int,
) (outs []string, histories [][]openai.ChatCompletionMessage, err error) {
	log.Debugf("calling GenerateN in domain")
	defer log.Debugf("GenerateN completed in domain")
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			out, hist, err := Generate(
				ctx,
				client,
				model,
				history,
				prompt,
			)
			if err != nil {
				return nil, nil, err
			}
			outs = append(outs, out)
			histories = append(histories, hist)
			if len(outs) >= n {
				return outs, histories, nil
			}
			history = append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: out})
			return outs, histories, nil
		}
	}
}

// GeneratePre is a function for generating text using the OpenAI API by
// prepending the prompt to the history.
func GeneratePre(
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
						[]openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: prmpt}},
						history...,
					),
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: "json",
					},
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

// GeneratePreN is a function for generating text using the OpenAI API by
// prepending the prompt to the history multiple "N" times.
func GeneratePreN(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
	n int,
) (outs []string, postHistories [][]openai.ChatCompletionMessage, err error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			out, hist, err := GeneratePre(
				ctx,
				client,
				model,
				history,
				prompt,
			)
			if err != nil {
				return nil, nil, err
			}
			outs = append(outs, out)
			postHistories = append(postHistories, hist)
			if len(outs) >= n {
				return outs, postHistories, nil
			}
			history = append(history, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: out})
			return outs, postHistories, nil
		}
	}
}

// GeneratePreTxt is a function for generating text using the OpenAI API by
// prepending the prompt to the history.
func GeneratePreTxt(
	ctx context.Context,
	client *openai.Client,
	model string,
	history []openai.ChatCompletionMessage,
	prompt prompter,
) (out string, postHistory []openai.ChatCompletionMessage, err error) {
	log.Debugf("calling GeneratePreTxt in domain")
	defer log.Debugf("GeneratePreTxt completed in domain")
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
						[]openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: prmpt}},
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
