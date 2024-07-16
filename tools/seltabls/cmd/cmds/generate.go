package cmds

import (
	"context"
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
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

const (
	// baseURLDefault is the base url for the openai api of groq
	baseURLDefault = "https://api.groq.com/openai/v1"
	// baseURLDefault = "https://api.openai.com/v1"
	// baseTreeWidthDefault is the default tree width for the openai api of groq
	baseTreeWidthDefault = 3
	defaultNumSections   = 3
	// // baseFastModelDefault is the default fast model for the openai api of groq
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
			client := openai.NewClient(params.LLMKey)
			/* client := domain.CreateClient( */
			/*	params.BaseURI, */
			/*	params.LLMKey, */
			/* ) */
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
			log.Infof("calling mainGenerate")
			if err := mainGenerate(
				ctx,
				sels,
				client,
				htmlBody,
				doc,
				params,
			); err != nil {
				return err
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
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
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
}

var mut sync.Mutex

func runGenerate(
	ctx context.Context,
	selectors []master.Selector,
	client *openai.Client,
	htmlBody []byte,
	doc *goquery.Document,
	params GenerateCmdParams,
) error {
	mut.Lock()
	defer mut.Unlock()
	identifyCompletions, identifyHistories, err := domain.InvokeJSONN(
		ctx,
		client,
		params.FastModel,
		[]openai.ChatCompletionMessage{},
		domain.IdentifyPromptArgs{
			URL:     params.URL,
			Content: string(htmlBody),

			NumSections: params.NumSections,
		},
		domain.IdentifyResponse{},
		params.TreeWidth,
		string(htmlBody),
	)
	log.Infof("identifyCompletions: %+v", identifyCompletions)
	if err != nil ||
		len(identifyHistories) == 0 {
		return fmt.Errorf("failed to generate identify completions: %w", err)
	}
	var identified domain.IdentifyResponse
	identifyCompletion, identifyHistory, err := domain.InvokeJSON(
		ctx,
		client,
		params.SmartModel,
		[]openai.ChatCompletionMessage{},
		domain.IdentifyAggregateArgs{
			Schemas:   identifyCompletions,
			Content:   string(htmlBody),
			Selectors: selectors,
		},
		&identified,
		string(htmlBody),
	)
	log.Infof("identifyCompletion: %+v", identifyCompletion)
	if err != nil || len(identifyCompletion) == 0 ||
		len(identifyHistories) == 0 {
		return fmt.Errorf("failed to generate identify completions: %w", err)
	}
	eg, ctx := errgroup.WithContext(ctx)
	for _, section := range identified.Sections {
		eg.Go(func() error {
			s, err := domain.HTMLSel(doc, section.CSS)
			if err != nil {
				return fmt.Errorf("failed to get html: %w", err)
			}
			selectorOuts, selectorHistories, err := domain.InvokeJSONN(
				ctx,
				client,
				params.SmartModel,
				identifyHistory,
				domain.StructPromptArgs{
					URL:       params.URL,
					Content:   s,
					Selectors: domain.HTMLReduce(doc, selectors),
				},
				domain.FieldsResponse{},
				params.TreeWidth,
				string(htmlBody),
			)
			if err != nil || len(selectorOuts) == 0 ||
				len(selectorHistories) == 0 {
				return fmt.Errorf(
					"failed to generateN struct completions: %w",
					err,
				)
			}
			var structFile = domain.StructFilePromptArgs{}
			_, _, err = domain.InvokeJSON(
				ctx,
				client,
				params.SmartModel,
				[]openai.ChatCompletionMessage{},
				domain.StructAggregateArgs{
					Selectors: domain.HTMLReduce(doc, selectors),
					Content:   domain.HTMLReduct(doc, section.CSS),
					Schemas:   selectorOuts,
				},
				&structFile,
				string(htmlBody),
			)
			if err != nil {
				return fmt.Errorf("failed to generate struct aggregate: %w", err)
			}
			structFile.Name = section.Name
			structFile.URL = params.URL
			out, err := domain.NewPrompt(structFile)
			if err != nil {
				return fmt.Errorf("failed to create struct file: %w", err)
			}
			name := fmt.Sprintf("%s.go", structFile.Name)
			if structFile.Name == "" {
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
			log.Infof("Generated struct file: %s.go", structFile.Name)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to generate structs: %w", err)
	}
	return nil
}
