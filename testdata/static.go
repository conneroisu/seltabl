package testdata

import (
	_ "embed"
)

type FixtureStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
}

//go:embed ab_num_table.html
var FixtureABNumTable string

type NumberedStruct struct {
	Header1 string `json:"Header 1" seltabl:"Header 1" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	Header2 string `json:"Header 2" seltabl:"Header 2" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
	Header3 string `json:"Header 3" seltabl:"Header 3" hSel:"tr:nth-child(1) td:nth-child(3)" dSel:"tr td:nth-child(3)" cSel:"$text"`
}

//go:embed numbered_table.html
var NumberedTable string

type SuperNovaStruct struct {
	Supernova string `json:"Supernova" seltabl:"Supernova" hSel:"tr:nth-child(1) th:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	Year      string `json:"Year" seltabl:"Year" hSel:"tr:nth-child(1) th:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
	Type      string `json:"Type" seltabl:"Type" hSel:"tr:nth-child(1) th:nth-child(3)" dSel:"tr td:nth-child(3)" cSel:"$text"`
	Distance  string `seltabl:"Distance" hSel:"tr:nth-child(1) th:nth-child(4)" dSel:"tr td:nth-child(4)" json:"Distance" `
	Notes     string `json:"Notes" seltabl:"Notes" hSel:"tr:nth-child(1) th:nth-child(5)" dSel:"tr td:nth-child(5)"`
}

//go:embed supernova.html
var SuperNovaTable string
