package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "seltabl",
	Short: "AI driven table parsing code generator for Go",
	Long: `
seltabl: A golang library for configurably parsing html sequences into stucts originally built for html tables, but can be used for any html sequence.

Command allows you to generate a golang struct from a html table given the data selectors for the table and the data selectors for the types of fields in the struct.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}
	return nil
}
