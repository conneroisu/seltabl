// Package main is the entry point for the command line tool
// a language server for the seltabl package called seltabls.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/cmd"
)

// main is the entry point for the command line tool, a
// language server for the seltabl package
func main() {
	ctx := context.Background()
	if err := cmd.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute command: %s\n", err)
		println("failed to execute command: " + err.Error())
		panic(err)
	}
}

// TableStruct is a test struct
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
type TableStruct struct {
	A string `hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" ctl:"text"`
	B string `hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" ctl:"text"`
	C string `hSel:"html body div[class=contentArea] table" dSel:"tr td:nth-child(1)" ctl:"text"`
	D string `hSel:"html head meta[name=csrf-token]" dSel:"tr td:nth-child(2)" ctl:"text"`
}
