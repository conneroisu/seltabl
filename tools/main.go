// Package main is the entry point for the command line tool
// a language server for the seltabl package
package main

import (
	"fmt"
	"os"

	"github.com/conneroisu/seltabl/tools/cmd"
	"github.com/conneroisu/seltabl/tools/pkg/analysis"
)

// main is the entry point for the command line tool, a
// language server for the seltabl package
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
func main() {
	rs := &cmd.Root{Writer: os.Stdout}
	rs.State = analysis.NewState(rs)
	if err := cmd.Execute(rs); err != nil {
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
