package cmd

import (
	"context"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/cmd/cmds"
	"github.com/conneroisu/seltabl/tools/seltabls/pkg/server"
	"github.com/sourcegraph/conc"
	"github.com/spf13/cobra"
)

// AddCommands adds the routes for the root command
func AddCommands(ctx context.Context, root *cobra.Command) error {
	wg := conc.WaitGroup{}
	wg.Go(func() {
		root.AddCommand(cmds.NewVetCmd(
			ctx,
			os.Stdout,
			os.Stdin,
		))
	})
	wg.Go(func() {
		root.AddCommand(cmds.NewLSPCmd(
			ctx,
			os.Stdout,
			os.Stdin,
			server.HandleMessage,
		))
	})
	wg.Go(func() {
		root.AddCommand(cmds.NewCompletionCmd(
			ctx,
			os.Stdout,
			os.Stdin,
		))
	})
	wg.Go(func() {
		root.AddCommand(cmds.NewGenerateCmd(
			ctx,
			os.Stdout,
			os.Stdin,
		))
	})
	wg.Wait()
	return nil
}
