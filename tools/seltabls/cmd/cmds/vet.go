package cmds

import (
	"context"
	"io"

	"github.com/spf13/cobra"
)

// NewVetCmd returns the vet command
func NewVetCmd(ctx context.Context, w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vet",
		Short: "Evaluate code for common errors or invalid selectors",
		Long: `
Similar to go vet, but for seltabl.
Evaluate code for common errors or invalid selectors.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// eg, ctx := errgroup.WithContext(ctx)

			return nil
		},
	}
	return cmd
}
