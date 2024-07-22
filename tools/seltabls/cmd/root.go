package cmd

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/log"
	"github.com/conneroisu/seltabl/tools/seltabls/data"
	"github.com/conneroisu/seltabl/tools/seltabls/data/master"
	"github.com/mitchellh/go-homedir"
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
	configPath, err := CreateConfigDir("~/.config/seltabls/")
	if err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	logPath := path.Join(configPath, "state.log")
	logFile, err := os.OpenFile(
		logPath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}
	log.SetOutput(logFile)
	log.SetLevel(log.DebugLevel)
	log.Infof("Starting seltabls")
	log.SetReportCaller(true)
	log.SetPrefix("seltabls")
	log.SetReportTimestamp(false)
	db, err := data.NewDb(
		ctx,
		master.New,
		&data.Config{
			Schema:   master.MasterSchema,
			URI:      "sqlite://uri.sqlite",
			FileName: path.Join(configPath, "uri.sqlite"),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	err = AddCommands(ctx, cmd, db)
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

// CreateConfigDir creates a new config directory and returns the path.
func CreateConfigDir(dirPath string) (string, error) {
	path, err := homedir.Expand(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to expand home directory: %w", err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		if os.IsExist(err) {
			return path, nil
		}
		return "", fmt.Errorf(
			"failed to create or find config directory: %w",
			err,
		)
	}
	return path, nil
}
