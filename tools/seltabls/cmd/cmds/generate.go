package cmds

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/domain"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/generate"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/llm"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/parsers"
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

// NewGenerateCmd returns the generate command.
func NewGenerateCmd(
	ctx context.Context,
	w io.Writer,
	r io.Reader,
) *cobra.Command {
	var url, name, fastModel, smartModel, llmKey, baseURI string
	var ignoreElements []string
	var treeWidth, treeDepth int
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
			var logLevel string
			log.Debugf("PreRunE called for command: %s", cmd.Name())
			defer log.Debugf("PreRunE completed for command: %s", cmd.Name())
			cmd.SetOutput(w)
			log.SetOutput(w)
			cmd.SetIn(r)
			cmd.SetErr(w)
			cmd.SetContext(ctx)
			cmd.PersistentFlags().StringVarP(
				&logLevel,
				"log-level",
				"v",
				"debug",
				`The log verbosity level to use.`,
			)
			lev, err := log.ParseLevel(logLevel)
			if err != nil {
				return fmt.Errorf("failed to parse log level: %w", err)
			}
			log.SetLevel(lev)
			cmd.PersistentFlags().StringVarP(
				&url,
				"url",
				"u",
				"",
				"The url for which to generate a seltabl struct go file, test file, and config file.",
			)
			cmd.PersistentFlags().StringVarP(
				&name,
				"name",
				"n",
				"",
				"The name of the struct to generate.",
			)
			cmd.PersistentFlags().StringVarP(&fastModel,
				"fast-model",
				"f",
				baseFastModelDefault,
				fmt.Sprintf(
					"The name of the fast model to use for generating the struct. Defaults to %s.",
					baseFastModelDefault,
				),
			)
			cmd.PersistentFlags().StringVarP(&smartModel,
				"smart-model",
				"s",
				baseSmartModelDefault,
				fmt.Sprintf(
					"The name of the smart model to use for generating the struct. Defaults to %s.",
					baseSmartModelDefault,
				),
			)
			cmd.PersistentFlags().StringVarP(
				&llmKey,
				"llm-key",
				"k",
				"",
				"The key for the llm model to use for generating the struct.",
			)
			cmd.PersistentFlags().StringVarP(
				&baseURI,
				"base-uri",
				"b",
				baseURLDefault,
				fmt.Sprintf(
					"The base uri for the openai api of groq. Defaults to %s.",
					baseURLDefault,
				),
			)
			cmd.PersistentFlags().StringArrayVarP(
				&ignoreElements,
				"ignore-elements",
				"i",
				baseIgnoreElementsDefault,
				fmt.Sprintf(
					"The elements to ignore when generating the struct. Defaults to %s.",
					baseIgnoreElementsDefault,
				),
			)
			cmd.PersistentFlags().IntVarP(
				&treeWidth,
				"tree-width",
				"w",
				baseTreeWidthDefault,
				fmt.Sprintf(
					"The width of the tree when generating the struct. Defaults to %d.",
					baseTreeWidthDefault,
				),
			)
			cmd.PersistentFlags().IntVarP(
				&treeDepth,
				"tree-depth",
				"d",
				baseTreeDepthDefault,
				fmt.Sprintf(
					"The depth of the tree when generating the struct. Defaults to %d.",
					baseTreeDepthDefault,
				),
			)

			if llmKey == "" {
				llmKey := os.Getenv("LLM_API_KEY")
				if llmKey == "" {
					return fmt.Errorf("LLM_API_KEY is not set")
				}
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			log.Debugf("RunE called for command: %s", cmd.Name())
			defer log.Debugf("RunE completed for command: %s", cmd.Name())
			if url == "" {
				input := huh.NewInput().
					Title("Enter the url for which to generate a seltabl struct:").
					Prompt("?").
					Validate(domain.IsURL).
					Value(&url)
				input.Run()
			}
			client := llm.CreateClient(
				baseURI,
				llmKey,
			)
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
				url,
				ignores,
			)
			if err != nil {
				return fmt.Errorf("failed to get selectors: %w", err)
			}
			log.Infof("Getting URL: %s", url)
			htmlBody, err := generate.GetURL(url, ignores)
			if err != nil {
				return fmt.Errorf("failed to get url: %w", err)
			}
			err = generate.Suite(
				ctx,
				treeWidth,
				treeDepth,
				client,
				fastModel,
				smartModel,
				name,
				url,
				string(htmlBody),
				ignores,
				sels,
			)
			if err != nil {
				return fmt.Errorf("failed to generate suite: %w", err)
			}
			return nil
		},
	}
}
