package main

import (
	"context"
	"fmt"

	"github.com/conneroisu/seltabl"
)

// ScrapeTeamCodes is a function for scraping team codes
func ScrapeTeamCodes(
	_ context.Context,
	content *string,
) (*[]NCAATeamCode, error) {
	teamCodes, err := seltabl.NewFromString[NCAATeamCode](*content)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode teams: %w",
			err,
		)
	}
	return &teamCodes, nil
}

// NCAATeamCode is a struct for a team code
type NCAATeamCode struct {
	ID   int    `json:"id"   seltabl:"ID"   hSel:"tr.grey_heading td:nth-child(1)" dSel:"tr:not(.grey_heading):not(.heading) td:nth-child(1)" cSel:"$text" must:"NCAA Codes"`
	Name string `json:"name" seltabl:"Name" hSel:"tr.grey_heading td:nth-child(2)" dSel:"tr:not(.grey_heading):not(.heading) td:nth-child(2)" cSel:"$text"`
}
