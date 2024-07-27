package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "embed"

	"github.com/stretchr/testify/assert"
)

//go:embed codes.html
var TeamCodesHTML string

// TestScrapeTeamCodes tests the ScrapeTeamCodes function
func TestScrapeTeamCodes(t *testing.T) {
	t.Run("Scrape Team Codes Static Inputs", func(t *testing.T) {
		t.Parallel()
		args := []struct {
			name    string
			content string
			wantErr bool
		}{
			{
				name:    "static input",
				content: TeamCodesHTML,
				wantErr: false,
			},
		}
		for _, arg := range args {
			arg := arg
			t.Run(arg.name, func(t *testing.T) {
				t.Parallel()
				ctx := context.Background()
				ctx, cancel := context.WithCancel(ctx)
				defer cancel()
				start := time.Now()
				teamResults, err := ScrapeTeamCodes(ctx, &arg.content)
				t.Logf("time: %s\n", time.Since(start))
				stop := time.Now()
				fmt.Printf("time: %s\n", stop.Sub(start))
				assert.Equal(t, arg.wantErr, err != nil)
				if arg.wantErr {
					return
				}
				assert.NoError(t, err)
				assert.NotNil(t, teamResults)
				assert.NotEmpty(t, teamResults)
				for _, teamCode := range *teamResults {
					assert.NotEmpty(t, teamCode.ID)
					assert.NotEmpty(t, teamCode.Name)
				}
			})
		}
	})
}
