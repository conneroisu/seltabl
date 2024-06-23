package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command
func NewRootCmd(ctx context.Context) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "seltabls", // the name of the command
		Short: "A command line tool for parsing html tables into structs",
		Long: `
CLI and Language Server for the seltabl package.

Language server provides completions, hovers, and code actions for seltabl defined structs.
	
CLI provides a command line tool for verifying, linting, and reporting on seltabl defined structs.
`,
	}
	err := AddRoutes(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to add routes: %w", err)
	}
	return cmd, nil
}

// Execute runs the root command
func Execute(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cmd, err := NewRootCmd(ctx)
	if err != nil {
		return fmt.Errorf("failed to create root command: %w", err)
	}
	if err := cmd.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("failed to execute root command: %w", err)
	}
	return nil
}
