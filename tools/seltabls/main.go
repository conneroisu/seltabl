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
	err := cmd.Execute(ctx)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"failed to execute command: %s\n",
			err,
		)
		println("failed to execute command: " + err.Error())
		panic(err)
	}
}

// TableStruct is a test struct.
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
type TableStruct struct {
	A string `hSel:"html>body>div.footer>div>span>a[href]" dSel:"html>body>div.footer>div>span>a[href]"    ctl:"$text"`
	B string `hSel:"html>body>div.footer>div>span>a[href]" dSel:"html>body>div>table>tbody>tr.row_odd"     ctl:"$text"`
	C string `hSel:"html>body>div>table>tbody>tr.row_even" dSel:"html>body>div>table>tbody>tr.row_even>td" ctl:"$text"`
	D string `hSel:"html>body>div.footer>div>span>a[href]" dSel:"html>body>div.footer>div>span>a[href]"    ctl:"$text"`
	E string `hSel:"html>body>div.footer>div>span>a[href]" dSel:"html>body>div.footer>div>span>a[href]"    ctl:"$text"`
}
