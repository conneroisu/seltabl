package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

const (
	// baseURLDefault is the base url for the openai api of groq
	baseURLDefault = "https://api.groq.com/openai/v1"
	// baseTreeWidthDefault is the default tree width for the openai api of groq
	baseTreeWidthDefault = 10
	// baseTreeDepthDefault is the default tree depth for the openai api of groq
	baseTreeDepthDefault = 10
	// baseFastModelDefault is the default fast model for the openai api of groq
	baseFastModelDefault = "llama3-8b-8192"
	// baseSmartModelDefault is the default smart model for the openai api of groq
	baseSmartModelDefault = "llama3-70b-8192"
)

// baseIgnoreElementsDefault is the default ignore elements for the openai api of groq
var baseIgnoreElementsDefault = []string{
	"script",
	"meta",
	"style",
	"link",
	"img",
	"footer",
	"header",
}

// GenerateCmdParams is a struct for the generate command parameters
type GenerateCmdParams struct {
	URL            string   `json:"url"`            // required
	Name           string   `json:"name"`           // required
	FastModel      string   `json:"fastModel"`      // required
	SmartModel     string   `json:"smartModel"`     // required
	LLMKey         string   `json:"llmKey"`         // required
	BaseURI        string   `json:"baseUri"`        // required
	IgnoreElements []string `json:"ignoreElements"` // required
	TreeWidth      int      `json:"treeWidth"`      // required
}

// NewGenerateCmd returns the generate command.
func NewGenerateCmd(
	ctx context.Context,
	w io.Writer,
	r io.Reader,
) *cobra.Command {
	var params GenerateCmdParams
	var client *openai.Client
	return &cobra.Command{
		Use:   "generate", // the name of the command
		Short: "Generates a new seltabl struct for a given url with test coverage.",
		Long: `
Generates a new seltabl struct for a given url.

The command will create a new package in the current directory with the name given.
Additionally, it will generate a test file and configuration with the name "{name}_test.go" and "{name}_seltabl.yaml" respectively in the current directory.
So the output fo the command:

.
└── name
    ├── subitem.go
    ├── subitem_test.go
    └── subitem_seltabl.yaml	

`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			log.Debugf("PreRunE called for command: %s", cmd.Name())
			defer log.Debugf("PreRunE completed for command: %s", cmd.Name())
			cmd.SetOutput(w)
			cmd.SetIn(r)
			cmd.SetErr(w)
			cmd.SetContext(ctx)
			cmd.PersistentFlags().StringVarP(
				&params.URL,
				"url",
				"u",
				"",
				"The url for which to generate a seltabl struct go file, test file, and config file.",
			)
			cmd.PersistentFlags().StringVarP(
				&params.Name,
				"name",
				"n",
				"",
				"The name of the struct to generate.",
			)
			cmd.PersistentFlags().StringVarP(
				&params.FastModel,
				"fast-model",
				"f",
				baseFastModelDefault,
				fmt.Sprintf(
					"The name of the fast model to use for generating the struct. Defaults to %s.",
					baseFastModelDefault,
				),
			)
			cmd.PersistentFlags().StringVarP(
				&params.SmartModel,
				"smart-model",
				"s",
				baseSmartModelDefault,
				fmt.Sprintf(
					"The name of the smart model to use for generating the struct. Defaults to %s.",
					baseSmartModelDefault,
				),
			)
			cmd.PersistentFlags().StringVarP(
				&params.LLMKey,
				"llm-key",
				"k",
				"",
				"The key for the llm model to use for generating the struct.",
			)
			cmd.PersistentFlags().StringVarP(
				&params.BaseURI,
				"base-uri",
				"b",
				baseURLDefault,
				fmt.Sprintf(
					"The base uri for the openai api of groq. Defaults to %s.",
					baseURLDefault,
				),
			)
			cmd.PersistentFlags().StringArrayVarP(
				&params.IgnoreElements,
				"ignore-elements",
				"i",
				baseIgnoreElementsDefault,
				fmt.Sprintf(
					"The elements to ignore when generating the struct. Defaults to %s.",
					baseIgnoreElementsDefault,
				),
			)
			cmd.PersistentFlags().IntVarP(
				&params.TreeWidth,
				"tree-width",
				"w",
				baseTreeWidthDefault,
				fmt.Sprintf(
					"The width of the tree when generating the struct. Defaults to %d.",
					baseTreeWidthDefault,
				),
			)
			cmd.PersistentFlags().IntVarP(
				&params.TreeWidth,
				"tree-depth",
				"d",
				baseTreeDepthDefault,
				fmt.Sprintf(
					"The depth of the tree when generating the struct. Defaults to %d.",
					baseTreeDepthDefault,
				),
			)
			if params.LLMKey == "" {
				llmKey := os.Getenv("LLM_API_KEY")
				if llmKey == "" {
					return fmt.Errorf("LLM_API_KEY is not set")
				}
			}

			if params.URL == "" {
				input := huh.NewInput().
					Title("Enter the url for which to generate a seltabl struct:").
					Prompt("?").
					Validate(domain.IsURL).
					Value(&params.URL)
				input.Run()
			}
			client = domain.CreateClient(
				params.BaseURI,
				params.LLMKey,
			)
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			log.Debugf("RunE called for command: %s", cmd.Name())
			defer log.Debugf("RunE completed for command: %s", cmd.Name())
			state, err := analysis.NewState()
			if err != nil {
				return fmt.Errorf("failed to create state: %w", err)
			}
			ignores := []string{
				"script",
				"style",
				"link",
				"img",
				"footer",
				"header",
			}
			sels, err := parsers.GetSelectors(
				ctx,
				&state.Database,
				params.URL,
				ignores,
			)
			if err != nil {
				return fmt.Errorf("failed to get selectors: %w", err)
			}
			log.Infof("Getting URL: %s", params.URL)
			htmlBody, err := domain.GetRuledURL(params.URL, ignores)
			if err != nil {
				return fmt.Errorf("failed to get url: %w", err)
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlBody)))
			if err != nil {
				return fmt.Errorf("failed to create document: %w", err)
			}
			if err := mainGenerate(
				ctx,
				sels,
				client,
				htmlBody,
				doc,
				params,
			); err != nil {
				return fmt.Errorf("failed to generate suite: %w", err)
			}

			return nil
		},
	}
}

func mainGenerate(
	ctx context.Context,
	selectors []master.Selector,
	client *openai.Client,
	htmlBody []byte,
	doc *goquery.Document,
	params GenerateCmdParams,
) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context cancelled: %w", ctx.Err())
	default:
		return runGenerate(
			ctx,
			selectors,
			client,
			htmlBody,
			doc,
			params,
		)
	}
}

func runGenerate(
	ctx context.Context,
	selectors []master.Selector,
	client *openai.Client,
	htmlBody []byte,
	doc *goquery.Document,
	params GenerateCmdParams,
) error {
	identifyCompletions, identifyHistories, err := domain.GenerateN(
		ctx,
		client,
		params.FastModel,
		[]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: string(htmlBody),
			},
		},
		domain.IdentifyPromptArgs{
			URL:     params.URL,
			Content: string(htmlBody),
		},
		params.TreeWidth,
	)
	if err != nil || len(identifyCompletions) == 0 || len(identifyHistories) == 0 {
		return fmt.Errorf("failed to generate identify completions: %w", err)
	}
	identifyCompletion, identifyHistory, err := domain.Generate(
		ctx,
		client,
		params.FastModel,
		[]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: string(htmlBody),
			},
		},
		domain.IdentifyAggregateArgs{
			Schemas: identifyCompletions,
			Content: string(htmlBody),
		},
	)
	if err != nil || len(identifyCompletion) == 0 || len(identifyHistories) == 0 {
		return fmt.Errorf("failed to generate identify completions: %w", err)
	}
retry:
	var identified domain.IdentifyResponse
	err = json.Unmarshal(
		[]byte(identifyCompletion),
		&identified,
	)
	if err != nil {
		identifyCompletion, identifyHistory, err = domain.GeneratePre(
			ctx,
			client,
			params.SmartModel,
			identifyHistory,
			domain.IdentifyErrorArgs{Error: err},
		)
		goto retry
	}
	for _, section := range identified.Sections {
		selectorOuts, selectorHistories, err := domain.GenerateN(
			ctx,
			client,
			params.SmartModel,
			identifyHistory,
			domain.StructPromptArgs{
				URL:       params.URL,
				Content:   domain.HtmlSel(doc, section.CSS),
				Selectors: selectors,
			},
			params.TreeWidth,
		)
		if err != nil || len(selectorOuts) == 0 || len(selectorHistories) == 0 {
			return fmt.Errorf("failed to generate struct completions: %w", err)
		}
		selectorOut, selectorHistory, err := domain.Generate(
			ctx,
			client,
			params.SmartModel,
			[]openai.ChatCompletionMessage{},
			domain.StructAggregateArgs{Selectors: selectors, Content: string(htmlBody), Schemas: selectorOuts},
		)
		if err != nil || len(selectorOut) == 0 || len(selectorHistory) == 0 {
			return fmt.Errorf("failed to generate struct aggregate: %w", err)
		}
		var structFile domain.StructFilePromptArgs
		err = json.Unmarshal(
			[]byte(selectorOut),
			&structFile,
		)
	}
	return nil
}
