package testdata

import (
	_ "embed"
)

// FixtureStruct is a test struct
type FixtureStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
}

//go:embed ab_num_table.html
var FixtureABNumTable string

// FixtureABNumTableResult is the expected result of parsing the ab_num_table
var FixtureABNumTableResult = []FixtureStruct{
	{A: "1", B: "2"},
	{A: "3", B: "4"},
	{A: "5", B: "6"},
	{A: "7", B: "8"},
}
