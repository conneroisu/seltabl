package cmd

import (
	"context"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/cmd/cmds"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/server"
	"github.com/spf13/cobra"
)

// AddCommands adds the routes for the root command.
func AddCommands(ctx context.Context, root *cobra.Command) error {
	root.AddCommand(cmds.NewVetCmd(
		ctx,
		os.Stdout,
		os.Stdin,
	))
	root.AddCommand(cmds.NewLSPCmd(
		ctx,
		os.Stdout,
		os.Stdin,
		server.HandleMessage,
	))
	root.AddCommand(cmds.NewCompletionCmd(
		ctx,
		os.Stdout,
		os.Stdin,
	))
	return nil
}
