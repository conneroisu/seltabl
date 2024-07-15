package domain

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"context"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"golang.org/x/sync/errgroup"

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
	Sections []Section `json:"sections"     yaml:"sections"`
	// Name is the name of the package.
	Name string `json:"name" yaml:"name"`
}

// Section is a struct for a section in the html.
type Section struct {
	// Name is the name of the section.
	Name string `json:"name"        yaml:"name"`
	// Description is a description of the section.
	Description string `json:"description" yaml:"description"`
	// CSS is the css selector for the section.
	CSS string `json:"css"         yaml:"css"`
	// Fields is a list of fields in the section.
	Fields []Field `json:"fields"      yaml:"fields"`
}

// FieldsResponse is a struct for the fields response
type FieldsResponse struct {
	Fields []Field `json:"fields" yaml:"fields"`
}

func (f FieldsResponse) respond() string {
	return "fields_response"
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

// CreateClient creates a new client for the given api key.
func CreateClient(baseURL string, apiKey string) *openai.Client {
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = baseURL
	cfg.APIVersion = string(openai.APITypeOpenAI)
	cfg.APIType = openai.APITypeOpenAI
	return openai.NewClientWithConfig(cfg)
}

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
	log.Debugf("calling InvokeJSON in domain with prompt: %s", prompt.prompt())
	defer log.Debugf(
		"InvokeJSON completed in domain with prompt: %s",
		prompt.prompt(),
	)
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
				ctx, openai.ChatCompletionRequest{
					Model:    model,
					Messages: genHistory,
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: "json_object",
					},
				},
			)
			if err != nil {
				r, ok := err.(*openai.RequestError)
				if ok {
					log.Debugf("request error: %+v", err)
					log.Debugf(
						"request error status code: %+v",
						r.HTTPStatusCode,
					)
				}
				r2, ok := err.(*openai.APIError)
				if ok {
					log.Debugf("rate limit error: %+v", err)
					log.Debugf("rate limit error type: %+v", r2.Type)
					log.Debugf("rate limit error param: %+v", r2.Param)
					log.Debugf("rate limit error code: %+v", r2.Code)
				}
				return "", genHistory, err
			}
			genHistory = append(genHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completion.Choices[0].Message.Content})
			log.Debugf(
				"completion: %+v",
				completion.Choices[0].Message.Content,
			)
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
	log.Debugf("calling InvokeJSON in domain with prompt: %s", prompt.prompt())
	defer log.Debugf(
		"InvokeJSON completed in domain with prompt: %s",
		prompt.prompt(),
	)
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
				ctx, openai.ChatCompletionRequest{
					Model:    model,
					Messages: genHistory,
					ResponseFormat: &openai.ChatCompletionResponseFormat{
						Type: "json_object",
					},
				},
			)
			if err != nil {
				r, ok := err.(*openai.RequestError)
				if ok {
					log.Debugf("request error: %+v", err)
					log.Debugf(
						"request error status code: %+v",
						r.HTTPStatusCode,
					)
				}
				r2, ok := err.(*openai.APIError)
				if ok {
					log.Debugf("rate limit error: %+v", err)
					log.Debugf("rate limit error type: %+v", r2.Type)
					log.Debugf("rate limit error param: %+v", r2.Param)
					log.Debugf("rate limit error code: %+v", r2.Code)
				}
				return "", genHistory, err
			}
			genHistory = append(genHistory, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: completion.Choices[0].Message.Content})
			log.Debugf(
				"completion: %+v",
				completion.Choices[0].Message.Content,
			)
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
				log.Debugf(
					"calling Generate from GenerateN in domain with prompt: %s",
					prompt.prompt(),
				)
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

// InvokeJSONTxtN is a function for generating text using the OpenAI API multiple "N" times.
func InvokeJSONTxtN(
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
			log.Debugf(
				"calling InvokeTxt from InvokeTxtN in domain with prompt: %s",
				prompt.prompt(),
			)
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
	log.Debugf("calling InvokeTxt in domain with prompt: %s", prompt.prompt())
	defer log.Debugf(
		"InvokeTxt completed in domain with prompt: %s",
		prompt.prompt(),
	)
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

// DecodeJSON is a function for decoding json.
//
// It tries to fix the json if it fails.
func DecodeJSON(
	ctx context.Context,
	data []byte,
	v interface{},
	history []openai.ChatCompletionMessage,
	client *openai.Client,
	model string,
	htmlBody string,
) error {

	hCtx, cancel := context.WithTimeout(ctx, time.Second*12)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var err error
			selR, ok := v.(FieldsResponse)
			if ok {
				for _, field := range selR.Fields {
					err = field.Verify(
						htmlBody,
					)
				}
			}
			if err == nil {
				err = json.Unmarshal(data, &v)
				if err == nil {
					return nil
				}
			}
			out, hist, err := InvokePreTxt(
				ctx,
				client,
				model,
				history,
				IdentifyErrorArgs{Error: err},
			)
			if err != nil {
				return err
			}
			newHist := append(hist, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: out})
			out, hist, err = InvokeJSONSimple(
				ctx,
				client,
				model,
				newHist,
				DecodeErrorArgs{Error: err},
			)
			if err != nil {
				return DecodeJSON(
					hCtx,
					data,
					v,
					hist,
					client,
					model,
					htmlBody,
				)
			}
			err = json.Unmarshal([]byte(out), v)
			if err != nil {
				return DecodeJSON(
					hCtx,
					data,
					v,
					hist,
					client,
					model,
					htmlBody,
				)
			}
			return nil
		}
	}
}

// force type cast for Responder
var _ Responder = (*IdentifyResponse)(nil)
var _ Responder = (*FieldsResponse)(nil)

// Verify checks if the selectors are in the html
func (f *Field) Verify(htmlBody string) error {
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(htmlBody),
	)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}
	if f.DataSelector != "" {
		sel := doc.Find(f.DataSelector)
		if sel.Length() == 0 {
			return fmt.Errorf("failed to find selector: %s", f.DataSelector)
		}
	} else {
		return fmt.Errorf("no data found for selector %s with type %s in field %s with type %s", f.DataSelector, f.Type, f.Name, f.Type)
	}
	if f.ControlSelector != "" {
		sel := doc.Find(f.ControlSelector)
		if sel.Length() == 0 {
			return fmt.Errorf("failed to find selector: %s", f.ControlSelector)
		}
	} else {
		return fmt.Errorf("no control found for selector %s with type %s in field %s with type %s", f.ControlSelector, f.Type, f.Name, f.Type)
	}
	if f.QuerySelector != "" {
		sel := doc.Find(f.QuerySelector)
		if sel.Length() == 0 {
			return fmt.Errorf("failed to find selector: %s", f.QuerySelector)
		}
	} else {
		return fmt.Errorf("no query found for selector %s with type %s in field %s with type %s", f.QuerySelector, f.Type, f.Name, f.Type)
	}
	if f.HeaderSelector != "" {
		sel := doc.Find(f.HeaderSelector)
		if sel.Length() == 0 {
			return fmt.Errorf("failed to find selector: %s", f.HeaderSelector)
		}
	}
	mbp := f.MustBePresent
	docTxt := doc.Text()
	if !strings.Contains(docTxt, mbp) {
		return fmt.Errorf("must be present (%s) not found for field %s with type %s", mbp, f.Name, f.Type)
	}
	return nil
}
