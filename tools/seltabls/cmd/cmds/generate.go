package cmds

import (
	"context"
	"io"

	"github.com/spf13/cobra"
)

// NewGenerateCmd returns the generate command
func NewGenerateCmd(_ context.Context, _ io.Writer, _ io.Reader) *cobra.Command {
	var url string
	cmd := &cobra.Command{
		Use:   "generate", // the name of the command
		Short: "Generates a new seltabl struct for a given url.",
		Long: `
Generates a new seltabl struct for a given url.

The command will create a new package in the current directory with the name "seltabl".
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.Help()
			return nil
		},
	}
	cmd.PersistentFlags().StringVarP(&url, "url", "u", "", "The url for which to generate a seltabl struct.")
	cmd.PersistentFlags().StringVarP(&url, "name", "n", "", "The name of the struct to generate.")
	cmd.PersistentFlags().StringVarP(&url, "llm-model", "m", "", "The name of the llm model to use for generating the struct.")
	cmd.PersistentFlags().StringVarP(&url, "llm-provider", "p", "", "The name of the llm provider to use for generating the struct.")
	registerCompletionFuncForGlobalFlags(cmd)
	return cmd
}

// registerCompletionFuncForGlobalFlags registers a completion function for the global flags
func registerCompletionFuncForGlobalFlags(cmd *cobra.Command) (err error) {
	err = cmd.RegisterFlagCompletionFunc(
		"url",
		func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{"https://github.com/conneroisu/seltabl/blob/main/testdata/ab_num_table.html"}, cobra.ShellCompDirectiveDefault
		},
	)
	return err
}
