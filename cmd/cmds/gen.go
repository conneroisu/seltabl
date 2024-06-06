package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewGenCmd returns a new cobra command for the gen subcommand
func NewGenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gen",
		Short: "Generates a scraper utilizing seltabl",
		Long: `
Subcommand to generate a scraper with a given file name.
		
Usage: 
		
	$ seltabl gen
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("gen called")
			return nil
		},
	}
}
