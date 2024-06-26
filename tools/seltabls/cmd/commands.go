package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/cmd/cmds"
	"github.com/spf13/cobra"
)

// AddCommands adds the routes for the root command
func AddCommands(ctx context.Context, root *cobra.Command) error {
	root.AddCommand(cmds.NewVetCmd(ctx, os.Stdout))
	root.AddCommand(cmds.NewLSPCmd(ctx, os.Stdout, cmds.HandleMessage))
	root.AddCommand(cmds.NewCompletionCmd(ctx, os.Stdout))
	genCmd, err := cmds.NewGenerateCmd(ctx, os.Stdout, os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to create generate command: %w", err)
	}
	root.AddCommand(genCmd)
	return nil
}
