package cmd

import (
	"github.com/conneroisu/seltabl/cmd/cmds"
	"github.com/spf13/cobra"
)

func AddRoutes(rootCmd *cobra.Command) {
	rootCmd.AddCommand(cmds.NewGenCmd())
}
