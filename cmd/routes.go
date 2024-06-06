package main

import (
	"github.com/spf13/cobra"
)

// AddRoutes adds the routes to the root command
func AddRoutes(rootCmd *cobra.Command, cfg *SeltablConfig) {
	rootCmd.AddCommand(NewGenerateCmd(cfg))
}
