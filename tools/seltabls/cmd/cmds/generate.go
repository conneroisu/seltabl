package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
	"github.com/google/uuid"
	"github.com/liushuangls/go-anthropic/v2"
	"github.com/ollama/ollama/api"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	// baseURLDefault is the base url for the openai api of groq
	// baseURLDefault = "https://api.groq.com/openai/v1"
	baseURLDefault = "https://api.openai.com/v1"
	// baseTreeWidthDefault is the default tree width for the openai api of groq
	baseTreeWidthDefault = 3
	defaultNumSections   = 3
	// // baseFastModelDefault is the default fast model for the openai api of groq
	// baseFastModelDefault = "llama3-8b-8192"
	// // baseSmartModelDefault is the default smart model for the openai api of groq
	// baseSmartModelDefault = "llama3-70b-8192"
	baseSmartModelDefault = "claude-3-5-sonnet-20240620"
	baseFastModelDefault  = "claude-3-haiku-20240307"
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
	FastModel      string   `json:"fastModel"`      // required
	SmartModel     string   `json:"smartModel"`     // required
	LLMKey         string   `json:"llmKey"`         // required
	BaseURI        string   `json:"baseUri"`        // required
	IgnoreElements []string `json:"ignoreElements"` // required
	TreeWidth      int      `json:"treeWidth"`      // required
	NumSections    int      `json:"numSections"`    // required
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
			cmd.PersistentFlags().IntVarP(
				&params.NumSections,
				"num-sections",
				"z",
				defaultNumSections,
				"The number of sections to generate.",
			)
			cmd.PersistentFlags().StringVarP(
				&params.FastModel,
				"fast-model",
				"f",
				baseFastModelDefault,
				"The name of the fast model to use for generating the struct. ",
			)
			cmd.PersistentFlags().StringVarP(
				&params.SmartModel,
				"smart-model",
				"s",
				baseSmartModelDefault,
				"The name of the smart model to use for generating the struct.",
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
				"The elements to ignore when generating the struct.",
			)
			cmd.PersistentFlags().IntVarP(
				&params.TreeWidth,
				"tree-width",
				"w",
				baseTreeWidthDefault,
				"The width of the tree when generating the struct.",
			)
			if params.LLMKey == "" {
				params.LLMKey = os.Getenv("LLM_API_KEY")
				if params.LLMKey == "" {
					return fmt.Errorf("LLM_API_KEY is not set")
				}
			}
			log.SetLevel(log.DebugLevel)
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
		RunE: func(_ *cobra.Command, _ []string) error {
			client := anthropic.NewClient(params.LLMKey)
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
			htmlBody, err := domain.GetRuledURL(
				params.URL,
				params.IgnoreElements,
			)
			if err != nil {
				return fmt.Errorf("failed to get url: %w", err)
			}
			doc, err := goquery.NewDocumentFromReader(
				strings.NewReader(string(htmlBody)),
			)
			if err != nil {
				return fmt.Errorf("failed to create document: %w", err)
			}
			print(doc, sels, client)
			log.Infof("calling mainGenerate")
			select {
			case <-ctx.Done():
				err = ctx.Err()
			default:
				err = runGenerate(
					ctx,
					sels,
					client,
					htmlBody,
					doc,
					params,
				)
				break
			}
			if err != nil {
				return err
			}
			return nil
		},
	}
}

var mut sync.Mutex

// TODO: need to rethink this with channels and goroutines
// TODO: each node in the tree has treeWidth.
func runGenerate(
	ctx context.Context,
	selectors []master.Selector,
	client *anthropic.Client,
	htmlBody []byte,
	doc *goquery.Document,
	params GenerateCmdParams,
) error {
	cli, err := api.ClientFromEnvironment()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	mut.Lock()
	defer mut.Unlock()
	log.Debugf("Generating Sections")
	identifyCompletions, _, err := domain.InvokeN(
		ctx,
		client,
		params.FastModel,
		domain.IdentifyArgs{
			URL:         params.URL,
			Content:     string(htmlBody),
			NumSections: params.NumSections,
			Selectors:   selectors,
		},
		params.TreeWidth,
	)
	if err != nil {
		return fmt.Errorf("failed to generate identify completions: %w", err)
	}
	log.Debugf("Generated (%d) Smart Sections: %s", len(identifyCompletions), strings.Join(identifyCompletions, "\n\n===\n\n"))
	log.Debugf("Aggregating Sections")
	var identifyCompletion string
	identifyCompletion, _, err = domain.Invoke(
		ctx,
		client,
		params.SmartModel,
		[]anthropic.Message{},
		domain.IdentifyAggregateArgs{
			Schemas:   identifyCompletions,
			Content:   string(htmlBody),
			Selectors: selectors,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to generate identify completion: %w", err)
	}
	resp, _ := json.MarshalIndent(identifyCompletion, "", "  ")
	log.Debugf("Generated Sections")
	log.Debugf("Generated Smart Sections: %s", string(resp))
	eg, ctx := errgroup.WithContext(ctx)
	var identified domain.IdentifyResponse
	err = domain.ChatUnmarshal(ctx, cli, []byte(identifyCompletion), &identified)
	if err != nil {
		log.Debugf("Failed to extract JSON from identifyCompletion: %s", identifyCompletion)
		return fmt.Errorf("failed to extract JSON from identifyCompletion: %w", err)
	}
	resp, _ = json.MarshalIndent(identified, "", "  ")
	log.Debugf("Unmarshaled Smart Sections: %s", string(resp))
	for _, section := range identified.Sections {
		var sectionSelectors = domain.HTMLReduce(doc, selectors)
		s, err := domain.HTMLSel(doc, section.CSS)
		if err != nil {
			return fmt.Errorf("failed to get html: %w", err)
		}
		log.Debugf("Generating Fast Selectors")
		selectorOuts, _, err := domain.InvokeN(
			ctx,
			client,
			params.FastModel,
			domain.StructPromptArgs{
				URL:       params.URL,
				Content:   s,
				Selectors: sectionSelectors,
			},
			params.TreeWidth,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to generateN struct completions: %w",
				err,
			)
		}
		for idx, selectorOut := range selectorOuts {
			var structFile = domain.StructFilePromptArgs{}
			err = domain.ChatUnmarshal(ctx, cli, []byte(selectorOut), &structFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal selectorOut: %w", err)
			}
			log.Debugf("Unmarshaled Fast Selectors (%d): %s", idx, structFile.Name)
		}
		log.Debugf("Generated (%d) Fast Selectors: %s", len(selectorOuts), strings.Join(selectorOuts, "\n\n===\n\n"))
		var structFile = domain.StructFilePromptArgs{}
		log.Debugf("Generated Selectors")
		log.Debugf("Aggregating Selectors")
		var smartStructAggregateResponse string
		smartStructAggregateResponse, _, err = domain.Invoke(
			ctx,
			client,
			params.SmartModel,
			[]anthropic.Message{},
			domain.StructAggregateArgs{
				Selectors: sectionSelectors,
				Content:   domain.HTMLReduct(doc, section.CSS),
				Schemas:   selectorOuts,
			},
		)
		if err != nil {
			return fmt.Errorf(
				"failed to generate struct aggregate: %w",
				err,
			)
		}
		err = domain.ChatUnmarshal(
			ctx,
			cli,
			[]byte(smartStructAggregateResponse),
			&structFile,
		)
		if err != nil {
			return fmt.Errorf("failed to unmarshal struct file: %w", err)
		}
		structFile.URL = params.URL
		structFile.IgnoreElements = params.IgnoreElements
		structFile.Fields = section.Fields
		out, err := domain.NewPrompt(structFile)
		if err != nil {
			return fmt.Errorf("failed to create struct file: %w", err)
		}
		name := fmt.Sprintf("%s.go", section.Name)
		if section.Name == "" {
			name = fmt.Sprintf("%s.go", uuid.New().String())
		}
		err = os.WriteFile(
			name,
			[]byte(out),
			0644,
		)
		if err != nil {
			return fmt.Errorf("failed to write struct file: %w", err)
		}
		log.Infof("Generated struct file: %s.go", name)
	}
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to generate structs: %w", err)
	}
	return nil
}
