package internal

import (
	"context"

	"github.com/conneroisu/seltabl/tools/internal/config"
	"github.com/spf13/cobra"
)

// AddRoutes adds the routes to the root command
func AddRoutes(ctx context.Context, rootCmd *cobra.Command, cfg *config.Config) {
	rootCmd.AddCommand(NewGenerateCmd(ctx, cfg))
	rootCmd.AddCommand(NewCompletionCommand(ctx))
}
