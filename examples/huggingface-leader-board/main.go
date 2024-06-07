// Package main shows how to use the seltabl package to scrape a table from a given url.
// The table used in this example is from the huggingface llm leader board.
package main

import (
	"fmt"
	"os"
)

// main scrapes from: https://huggingface.co/spaces/HuggingFaceH4/LLM-Leaderboard
func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// run runs the example
func run() error {
	fmt.Println("Hello, World from llm leader board!")
	return nil
}
