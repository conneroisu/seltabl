//go:build ignore
// +build ignore

package testdata

import (
	"fmt"
	"os"

	"github.com/conneroisu/seltabl/tools/cmd"
)

// main is the entry point for the command line tool, a language server for the seltabl package.
func main() {
	if err := cmd.Execute(); err != nil {
		// log to logs.txt
		file, err := os.OpenFile(
			"logs.txt",
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0666,
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		_, err = fmt.Println(err)
		if err != nil {
			file.WriteString(err.Error())
		}
		fmt.Fprintf(file, "%s\n", err)
		os.Exit(1)
	}
}

// TableStruct is a struct for a table
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
type TableStruct struct {
	A string `json"ctl" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"" ctl:"text"`
	B string `json"ctl" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"" ctl:"text"`
}
