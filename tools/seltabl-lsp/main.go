// Package main is the entry point for the command line tool
// a language server for the seltabl package called seltabl-lsp.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/cmd"
	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/analysis"
)

// main is the entry point for the command line tool, a
// language server for the seltabl package
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
func main() {
	rs := &cmd.Root{Writer: os.Stdout}
	rs.State = analysis.NewState(os.Getenv)
	rs.Logger = getLogger("./seltabl.log")
	if err := cmd.Execute(rs); err != nil {
		// log to logs.txt
		file, err := os.OpenFile(
			"logs.txt",
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0666,
		)
		if err != nil {
			rs.Logger.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		_, err = fmt.Println(err)
		if err != nil {
			_, err = file.WriteString(err.Error())
			if err != nil {
				rs.Logger.Println(err)
				os.Exit(1)
			}
		}
		rs.Logger.Println(err)
		os.Exit(1)
	}

}

// getLogger returns a logger that writes to a file
func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(
		fileName,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(logFile, "[seltabl-lsp]", log.LstdFlags)
}
