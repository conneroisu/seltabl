// Package main is the entry point for the command line tool
// a language server for the seltabl package called seltabl-lsp.
package main

import (
	"fmt"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/cmd"
)

// main is the entry point for the command line tool, a
// language server for the seltabl package
func main() {
	if err := cmd.Execute(); err != nil {
		out := os.Stderr
		fmt.Fprintf(out, "failed to execute command: %s\n", err)
		fmt.Fprintf(out, "exiting...\n")
		os.Exit(1)
	}
}

// TableStruct is a test struct
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
type TableStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" ctl:"text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(2)" ctl:"text"`
	C string `json:"c" seltabl:"c" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(3)" ctl:"text"`
	D string `json:"d" seltabl:"d" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(4)" ctl:"text"`
}
