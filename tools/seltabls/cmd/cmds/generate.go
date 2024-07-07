package cmds

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/generate"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/llm"
	"github.com/spf13/cobra"
)

const content = `package main`

// baseURL is the base url for the openai api of groq
const baseURL = "https://api.groq.com/openai/v1"

// NewGenerateCmd returns the generate command
func NewGenerateCmd(
	ctx context.Context,
	w io.Writer,
	r io.Reader,
) *cobra.Command {
	var url string
	var name string
	var llmModel string
	var llmKey string
	var ignoreElements []string
	cmd := &cobra.Command{
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
				&url,
				"url",
				"u",
				"",
				"The url for which to generate a seltabl struct.",
			)
			cmd.PersistentFlags().StringVarP(
				&name,
				"name",
				"n",
				"",
				"The name of the struct to generate.",
			)
			cmd.PersistentFlags().StringVarP(
				&llmModel,
				"llm-model",
				"m",
				"llama3-70b-8192",
				"The name of the llm model to use for generating the struct.",
			)
			cmd.PersistentFlags().StringVarP(
				&llmKey,
				"llm-key",
				"k",
				"",
				"The key for the llm model to use for generating the struct.",
			)
			cmd.PersistentFlags().StringArrayVarP(
				&ignoreElements,
				"ignore-elements",
				"i",
				[]string{"script", "meta", "style", "link", "img", "footer", "header"},
				"The elements to ignore when generating the struct.",
			)
			if llmKey == "" {
				llmKey := os.Getenv("LLM_API_KEY")
				if llmKey == "" {
					return fmt.Errorf("LLM_API_KEY is not set")
				}
			}
			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			if url == "" {
				input := huh.NewInput().
					Title("Enter the url for which to generate a seltabl struct:").
					Prompt("?").
					Validate(generate.IsURL).
					Value(&url)
				input.Run()
			}
			client := llm.CreateClient(
				url,
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
			sels, err := analysis.GetSelectors(
				ctx,
				&state,
				url,
				ignores,
			)
			if err != nil {
				return fmt.Errorf("failed to get selectors: %w", err)
			}
			htmlBody, err := generate.GetURL(url)
			if err != nil {
				return fmt.Errorf("failed to get url: %w", err)
			}
			err = generate.Suite(
				ctx,
				client,
				name,
				url,
				ignores,
				string(htmlBody),
				sels,
			)
			if err != nil {
				return fmt.Errorf("failed to generate suite: %w", err)
			}
			return nil
		},
	}
	return cmd
}

// verify verifies the generated struct
func verify(
	ctx context.Context,
	cmd *cobra.Command,
	name string,
) error {
	err := NewVetCmd(ctx, os.Stdout).RunE(cmd, []string{name + ".go"})
	if err == nil {
		fmt.Printf("Generated %s\n", name)
		return nil
	}
	return fmt.Errorf("failed to vet generated struct: %w", err)
}
