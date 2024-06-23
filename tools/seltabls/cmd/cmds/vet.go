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
		Short: "Vet the project",
		Long: `
Vet the project.
`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
	return cmd
}
