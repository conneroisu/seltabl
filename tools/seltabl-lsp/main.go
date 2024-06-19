// Package main is the entry point for the command line tool
// a language server for the seltabl package called seltabl-lsp.
package main

import (
	"os"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/cmd"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/internal/config"
)

// main is the entry point for the command line tool, a
// language server for the seltabl package
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
func main() {
	cfg, err := config.CreateConfigDir()
	if err != nil {
		panic(err)
	}
	rs := &cmd.Root{Writer: os.Stdout, Config: cfg}
	if err := cmd.Execute(rs); err != nil {
		rs.Logger.Println(err)
	}
}
