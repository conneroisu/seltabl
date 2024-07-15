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
	// baseURLDefault = "https://api.groq.com/openai/v1"
	baseURLDefault = "https://api.openai.com/v1"
	// baseTreeWidthDefault is the default tree width for the openai api of groq
	baseTreeWidthDefault = 10
	// baseFastModelDefault is the default fast model for the openai api of groq
	// baseFastModelDefault = "llama3-8b-8192"
	// // baseSmartModelDefault is the default smart model for the openai api of groq
	// baseSmartModelDefault = "llama3-70b-8192"
	baseSmartModelDefault = "gpt-4o"
	baseFastModelDefault  = "gpt-4o"
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
				"https://stats.ncaa.org/team/2/stats/16540",
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
			if params.LLMKey == "" {
				params.LLMKey = os.Getenv("LLM_API_KEY")
				if params.LLMKey == "" {
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
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			println("Getting model, llmKey: ", params.LLMKey)
			client := openai.NewClient(params.LLMKey)
			/*                  client := domain.CreateClient( */
			/* params.BaseURI, */
			/* params.LLMKey, */
			/*    ) */
			a, err := client.GetModel(ctx, params.FastModel)
			if err != nil {
				return fmt.Errorf("failed to get model: %w", err)
			}
			log.Infof("model: %s", a.ID)
			log.Infof("fastModel: %s", params.FastModel)
			log.Infof("smartModel: %s", params.SmartModel)
			log.Debugf("RunE called for command: %s", cmd.Name())
			defer log.Debugf("RunE completed for command: %s", cmd.Name())
			state, err := analysis.NewState()
			if err != nil {
				return fmt.Errorf("failed to create state: %w", err)
			}
			sels, err := parsers.GetSelectors(
				ctx,
				&state.Database,
				params.URL,
				params.IgnoreElements,
			)
			if err != nil {
				return fmt.Errorf("failed to get selectors: %w", err)
			}
			log.Infof("Getting URL: %s", params.URL)
			htmlBody, err := domain.GetRuledURL(params.URL, params.IgnoreElements)
			if err != nil {
				return fmt.Errorf("failed to get url: %w", err)
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlBody)))
			if err != nil {
				return fmt.Errorf("failed to create document: %w", err)
			}
			log.Infof("calling mainGenerate")
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
	log.Infof("calling runGenerate")
	defer log.Infof("runGenerate completed")
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
	log.Infof("calling runGenerate")
	defer log.Infof("runGenerate completed")
	log.Infof("calling domain.GenerateN")
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
		betterErr, _, err := domain.GeneratePreTxt(
			ctx,
			client,
			params.SmartModel,
			identifyHistory,
			domain.IdentifyErrorArgs{Error: err},
		)
		if err != nil {
			return fmt.Errorf("failed to generate error: %w", err)
		}
		identifyCompletion, identifyHistory, err = domain.Generate(
			ctx,
			client,
			params.FastModel,
			append(identifyHistory, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleAssistant, Content: betterErr}),
			domain.IdentifyAggregateArgs{
				Schemas: identifyCompletions,
				Content: string(htmlBody),
			},
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
				Content:   domain.HTMLSel(doc, section.CSS),
				Selectors: domain.HTMLReduce(doc, selectors),
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
			domain.StructAggregateArgs{
				Selectors: domain.HTMLReduce(doc, selectors),
				Content:   domain.HTMLReduct(doc, section.CSS),
				Schemas:   selectorOuts,
			},
		)
		if err != nil || len(selectorOut) == 0 || len(selectorHistory) == 0 {
			return fmt.Errorf("failed to generate struct aggregate: %w", err)
		}
		var structFile domain.StructFilePromptArgs
		err = json.Unmarshal(
			[]byte(selectorOut),
			&structFile,
		)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct file: %w", err)
		}
		structFile.PackageName = params.Name
		structFile.Name = section.Name
		structFile.URL = params.URL
		for _, s := range selectors {
			structFile.IgnoreElements = append(structFile.IgnoreElements, s.Value)
		}
		out, err := domain.NewPrompt(structFile)
		if err != nil {
			return fmt.Errorf("failed to create struct file: %w", err)
		}
		err = os.WriteFile(
			fmt.Sprintf("%s.go", structFile.Name),
			[]byte(out),
			0644,
		)
		if err != nil {
			return fmt.Errorf("failed to write struct file: %w", err)
		}
	}
	return nil
}
