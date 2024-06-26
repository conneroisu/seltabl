package cmds

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/analysis"
	"github.com/spf13/cobra"
)

// NewVetCmd returns the vet command
func NewVetCmd(ctx context.Context, w io.Writer) *cobra.Command {
	var verbose bool
	cmd := &cobra.Command{
		Use:   "vet",
		Short: "Evaluate code for common errors or invalid selectors",
		Long: `
Similar to go vet, but for seltabl.
Evaluate code for common errors or invalid selectors.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(w)
			cmd.SetContext(ctx)
			ctx := context.Background()
			files, err := filepath.Glob(args[0])

			state, err := analysis.NewState()
			if err != nil {
				return fmt.Errorf("failed to create state: %w", err)
			}
			for _, file := range files {
				if verbose {
					fmt.Printf("file: %s\n", file)
				}
				f, err := os.Open(file)
				if err != nil {
					return fmt.Errorf("failed to open file: %w", err)
				}
				// if it is not a go file then skip
				if filepath.Ext(file) != ".go" {
					continue
				}
				bdy, err := io.ReadAll(f)
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
				a := string(bdy)
				diags, err := state.OpenDocument(ctx, file, &a)
				if err != nil {
					return fmt.Errorf("failed to open document: %w", err)
				}
				for _, diag := range diags {
					fmt.Printf("Message %s\n", diag.Message)
					fmt.Printf("Range %+v\n", diag.Range)
					fmt.Printf("Severity %s\n", diag.Severity)
					fmt.Printf("Source %s\n", diag.Source)
				}
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVarP(
		&verbose,
		"verbose",
		"v",
		false,
		"The name of the struct to generate.",
	)
	return cmd
}
