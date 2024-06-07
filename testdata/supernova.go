package testdata

import (
	_ "embed"
)

// SuperNovaStruct is a test struct
type SuperNovaStruct struct {
	Supernova string `json:"Supernova" seltabl:"Supernova" hSel:"tr:nth-child(1) th:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	Year      string `json:"Year"      seltabl:"Year"      hSel:"tr:nth-child(1) th:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
	Type      string `json:"Type"      seltabl:"Type"      hSel:"tr:nth-child(1) th:nth-child(3)" dSel:"tr td:nth-child(3)" cSel:"$text"`
	Distance  string `json:"Distance"  seltabl:"Distance"  hSel:"tr:nth-child(1) th:nth-child(4)" dSel:"tr td:nth-child(4)"`
	Notes     string `json:"Notes"     seltabl:"Notes"     hSel:"tr:nth-child(1) th:nth-child(5)" dSel:"tr td:nth-child(5)"`
}

//go:embed supernova.html
var SuperNovaTable string

// SuperNovaTableResult is the expected result of parsing the supernova table
var SuperNovaTableResult = []SuperNovaStruct{
	{
		Supernova: "SN 1006",
		Year:      "1006",
		Type:      "Type Ia",
		Distance:  "7,200",
		Notes:     "Brightest recorded supernova in history",
	},
	{
		Supernova: "SN 1054 (Crab Nebula)",
		Year:      "1054",
		Type:      "Type II",
		Distance:  "6,500",
		Notes:     "Formed the Crab Nebula and pulsar",
	},
	{
		Supernova: "SN 1572 (Tycho's Supernova)",
		Year:      "1572",
		Type:      "Type Ia",
		Distance:  "8,000-10,000",
		Notes:     "Observed by Tycho Brahe",
	},
	{
		Supernova: "SN 1604 (Kepler's Supernova)",
		Year:      "1604",
		Type:      "Type Ia",
		Distance:  "20,000",
		Notes:     "Last observed supernova in the Milky Way",
	},
	{
		Supernova: "SN 1987A",
		Year:      "1987",
		Type:      "Type II",
		Distance:  "168,000",
		Notes:     "Closest observed supernova since 1604",
	},
	{
		Supernova: "SN 1993J",
		Year:      "1993",
		Type:      "Type IIb",
		Distance:  "11,000,000",
		Notes:     "In the galaxy M81",
	},
}
