package main

import (
	_ "embed"
)

//go:embed base_prompt.md
var basePrompt string

//go:embed output.tmpl
var output string

//go:embed output_test.tmpl
var outputTest string
