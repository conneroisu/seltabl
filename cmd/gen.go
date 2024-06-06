package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"unicode"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// NewGenerateCmd returns a new cobra command for the gen subcommand
func NewGenerateCmd(ctf *SeltablConfig) *cobra.Command {
	cmd := &cobra.Command{}
	cmd.Use = "generate"
	cmd.Short = "Generates a scraper utilizing seltabl"
	cmd.Long = `
Subcommand to generate a scraper with a given file name.
		
Usage: 
		
	$ seltabl generate
`
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		fmt.Println("gen called")
		fmt.Println(fmt.Sprintf("output file: %s", ctf.OutputFile))
		fmt.Println(fmt.Sprintf("input file: %s", ctf.InputFile))
		form := huh.NewInput().
			Title("What is the package name?").
			Prompt("?").
			Validate(isValidPackageName).
			Value(&ctf.PackageName)
		err := form.Run()
		if err != nil {
			return fmt.Errorf("failed to run form: %w", err)
		}
		form = huh.NewInput().
			Title("What is the url?").
			Prompt("?").
			Validate(isValidURL).
			Value(&ctf.URL)
		err = form.Run()
		if err != nil {
			return fmt.Errorf("failed to run form: %w", err)
		}
		form = huh.NewInput().
			Title("What is the struct name?").
			Prompt("?").
			Validate(isValidPackageName).
			Value(&ctf.StructName)
		err = form.Run()
		if err != nil {
			return fmt.Errorf("failed to run form: %w", err)
		}
		var fields []string
		selForm := huh.NewMultiSelect[string]().
			Options(
				getFieldOptions(ctf)...,
			).
			Title("Toppings").
			Limit(-1).
			Value(&fields)
		err = selForm.Run()
		if err != nil {
			return fmt.Errorf("failed to run form: %w", err)
		}
		file, err := os.Create(ctf.InputFile)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
		err = json.NewEncoder(file).Encode(fields)
		if err != nil {
			return fmt.Errorf("failed to encode json: %w", err)
		}
		// create the output_test.tmpl file
		outputTestFile, err := os.Create(ctf.PackageName + "_test.go")
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer outputTestFile.Close()
		comp := OutputTestDoc(*ctf)
		if err != nil {
			return fmt.Errorf("failed to render output test file: %w", err)
		}
		comp.Render(ctx, outputTestFile)
		red, err := os.ReadFile(ctf.PackageName + "_test.go")
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		outputTestFile.Truncate(0)
		// remove all \&#34; from the output replace with \"
		redStr := string(red)
		redStr = strings.ReplaceAll(redStr, "&#34;", "\"")
		redStr = strings.ReplaceAll(redStr, "&amp;", "&")
		split := strings.Split(redStr, "\n")
		splittee := []string{}
		prefixHit := false
		for i, line := range split {
			if !prefixHit {
				if strings.HasPrefix(line, "package ") {
					prefixHit = true
				}
			}
			if i == 0 {
				continue
			}
			splittee = append(splittee, line)
		}
		red = []byte(strings.Join(splittee, "\n"))
		outputTestFile.WriteString(redStr)

		comd := exec.Command("nvim", ctf.InputFile)
		comd.Stdout = os.Stdout
		comd.Stderr = os.Stderr
		err = comd.Run()
		if err != nil {
			return fmt.Errorf("failed to run nvim: %w", err)
		}
		return nil
	}
	cmd.PersistentFlags().StringVar(&ctf.OutputFile, "output", "scraper.go", "output file name")
	cmd.PersistentFlags().StringVar(&ctf.InputFile, "input", "input.json", "input file name")
	return cmd
}

// isValidPackageName validates the package name
func isValidPackageName(val string) error {
	if val == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	// Check if the first character is a letter
	if !unicode.IsLetter(rune(val[0])) {
		return fmt.Errorf("package name must start with a letter")
	}
	// Check the remaining characters
	for _, r := range val[1:] {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("package name can only contain letters and digits")
		}
	}
	return nil
}

// isValidURL validates the url
func isValidURL(val string) error {
	if val == "" {
		return fmt.Errorf("url cannot be empty")
	}
	_, err := url.ParseRequestURI(val)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	return nil
}
