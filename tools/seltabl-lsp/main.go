// Package main is the entry point for the command line tool
// a language server for the seltabl package called seltabl-lsp.
package main

import (
	"os"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/cmd"
)

// main is the entry point for the command line tool, a
// language server for the seltabl package
func main() {
	rs := &cmd.Root{Writer: os.Stdout}
	if err := cmd.Execute(rs); err != nil {
		rs.Logger.Println(err)
	}
}

// TableStruct is a test struct
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
type TableStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" ctl:"text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" ctl:"text"`
}
