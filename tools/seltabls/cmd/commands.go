package cmd

import (
	"context"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/cmd/cmds"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/server"
	"github.com/spf13/cobra"
)

// AddCommands adds the routes for the root command.
func AddCommands(
	ctx context.Context,
	root *cobra.Command,
	db *data.Database[master.Queries],
) error {
	root.AddCommand(cmds.NewVetCmd(
		ctx,
		os.Stdout,
		os.Stdin,
		db,
	))
	root.AddCommand(cmds.NewLSPCmd(
		ctx,
		os.Stdout,
		os.Stdin,
		server.HandleMessage,
		db,
	))
	root.AddCommand(cmds.NewCompletionCmd(
		ctx,
		os.Stdout,
		os.Stdin,
	))
	return nil
}
