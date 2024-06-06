package main

import (
	_ "embed"
)

//go:embed base_prompt.md
var basePrompt string

//go:embed output.tmpl
var output string
