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
	var ctx context.Context
	var err error
	ctx = context.Background()
	err = cmd.Execute(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute command: %s\n", err)
		println("failed to execute command: " + err.Error())
		panic(err)
	}
}

// TableStruct is a test struct.
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
type TableStruct struct {
	A string `json:"a" hSel:"html head meta[name=csrf-param]"     dSel:"tr td:nth-child(1)" ctl:"text"`
	B string `json:"b" hSel:"html"                                dSel:"tr:nth-child(1)"    ctl:"text"`
	C string `json:"c" hSel:"html > head > meta[name=csrf-token]" dSel:"tr td:nth-child(1)" ctl:"text"`
	D string `json:"d" hSel:"html > body > div.contentArea > table > tbody > tr.row_odd"                 dSel:"html"               ctl:"text"`
}
