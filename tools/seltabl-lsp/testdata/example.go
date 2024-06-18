package testdata

// @url: https://stats.ncaa.org/game_upload/team_codes
type TeamCode struct {
	ID   int    `json:"id"   seltabl:"ID"   hSel:"tr.grey_heading td:nth-child(1)" dSel:"tr:not(.grey_heading):not(.heading) td:nth-child(1)" ctl:"text" must:"NCAA Codes"`
	Name string `json:"name" seltabl:"Name" hSel:"tr.grey_heading td:nth-child(2)" dSel:"tr:not(.grey_heading):not(.heading) td:nth-child(2)" ctl:"text"`
}
