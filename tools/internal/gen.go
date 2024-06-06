package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/conneroisu/seltabl/tools/internal/config"
	"github.com/conneroisu/seltabl/tools/internal/editor"
	"github.com/spf13/cobra"
)

// NewGenerateCmd returns a new cobra command for the gen subcommand
func NewGenerateCmd(_ context.Context, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "generate"
	cmd.Short = "Generates a scraper utilizing seltabl"
	cmd.Long = `
Subcommand to generate a scraper with a given file name.
		
Usage: 
		
	$ seltabl generate
`
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := cmd.RunE(cmd, args); err != nil {
			fmt.Println(fmt.Sprintf("failed to run generate command: %v", err))
			os.Exit(1)
		}
	}

	cmd.RunE = func(_ *cobra.Command, _ []string) error {
		testFilePath := fmt.Sprintf("%s_test.go", cfg.PackageName)
		var err error
		if err = huh.NewInput().
			Title("What is the package name? (the name of the generated file and struct)").
			Prompt("?").
			Validate(isValidPackageName).
			Value(&cfg.PackageName).Run(); err != nil {
			return fmt.Errorf("failed to run form: %w", err)
		}
		if err = huh.NewInput().
			Title("What is the full URL?").
			Prompt("?").
			Validate(isValidURL).Value(&cfg.URL).Run(); err != nil {
			return fmt.Errorf("failed to run form: %w", err)
		}
		var fields []string
		body, err := get(cfg.URL)
		if err != nil {
			return fmt.Errorf("failed to get url: %w", err)
		}
		options, err := GetFieldOptions(cfg, body)
		if err != nil {
			return fmt.Errorf("failed to get field options: %w", err)
		}
		selForm := huh.NewMultiSelect[string]().
			Options(
				options...,
			).
			Title("Toppings").
			Limit(-1).
			Value(&fields)
		err = selForm.Run()
		if err != nil {
			return fmt.Errorf("failed to run form: %w", err)
		}
		file, err := os.Create(cfg.InputFilePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
		err = json.NewEncoder(file).Encode(fields)
		if err != nil {
			return fmt.Errorf("failed to encode json: %w", err)
		}
		outputTestFile, err := os.Create(testFilePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer outputTestFile.Close()
		app, args := editor.GetEditor()
		editCmd := exec.Command(app, strings.Join(args, " "), cfg.InputFilePath)
		editCmd.Stdout = os.Stdout
		editCmd.Stderr = os.Stderr
		err = editCmd.Run()
		if err != nil {
			return fmt.Errorf("failed to run nvim: %w", err)
		}
		return nil
	}
	cmd.PersistentFlags().StringVar(&cfg.OutputFilePath, "output", "scraper.go", "output file name")
	cmd.PersistentFlags().StringVar(&cfg.InputFilePath, "input", "input.json", "input file name")
	cmd.PersistentFlags().StringVar(&cfg.PackageName, "package", "seltabl", "package name")
	cmd.PersistentFlags().StringVar(&cfg.StructName, "struct", "Scraper", "struct name")
	cmd.PersistentFlags().StringVar(&cfg.URL, "url", "", "url")
	return cmd
}
