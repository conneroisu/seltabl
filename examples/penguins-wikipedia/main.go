// Package main is the an example of how to use the seltabl package.
// for the seltabl package to scrape a html table from a given url.
// The table used in this example is from the wikipedia page for
// penguins.
package main

import (
	"fmt"
	"os"
)

// Scrapes from: https://en.wikipedia.org/wiki/List_of_penguins
func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// run runs the example
func run() error {
	fmt.Println("Hello, World from Example2!")
	return nil
}
