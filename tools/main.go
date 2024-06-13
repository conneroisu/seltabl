// Package main is the entry point for the command line tool
// a language server for the seltabl package
package main

import (
	"fmt"
	"os"

	"github.com/conneroisu/seltabl/tools/cmd"
)

// main is the entry point for the command line tool
func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
