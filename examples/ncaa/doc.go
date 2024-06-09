// Package main is the an example of how to use the seltabl package.
// for the seltabl package
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

var url = "https://stats.ncaa.org/game_upload/team_codes"

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func run(ctx context.Context) error {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	str := string(body)
	codes, err := ScrapeTeamCodes(ctx, &str)
	if err != nil {
		panic(err)
	}
	for _, code := range *codes {
		fmt.Printf("%+v\n", code)
	}
	return nil
}
