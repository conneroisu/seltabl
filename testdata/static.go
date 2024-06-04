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

var FixtureABNumTableResult = []FixtureStruct{
	{A: "1", B: "2"},
	{A: "3", B: "4"},
	{A: "5", B: "6"},
	{A: "7", B: "8"},
}

type NumberedStruct struct {
	Header1 string `json:"Header 1" seltabl:"Header 1" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	Header2 string `json:"Header 2" seltabl:"Header 2" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
	Header3 string `json:"Header 3" seltabl:"Header 3" hSel:"tr:nth-child(1) td:nth-child(3)" dSel:"tr td:nth-child(3)" cSel:"$text"`
}

//go:embed numbered_table.html
var NumberedTable string

var NumberedTableResult = []NumberedStruct{
	{Header1: "Row 1, Cell 1", Header2: "Row 1, Cell 2", Header3: "Row 1, Cell 3"},
	{Header1: "Row 2, Cell 1", Header2: "Row 2, Cell 2", Header3: "Row 2, Cell 3"},
	{Header1: "Row 3, Cell 1", Header2: "Row 3, Cell 2", Header3: "Row 3, Cell 3"},
}

type SuperNovaStruct struct {
	Supernova string `json:"Supernova" seltabl:"Supernova" hSel:"tr:nth-child(1) th:nth-child(1)" dSel:"tr td:nth-child(1)" cSel:"$text"`
	Year      string `json:"Year" seltabl:"Year" hSel:"tr:nth-child(1) th:nth-child(2)" dSel:"tr td:nth-child(2)" cSel:"$text"`
	Type      string `json:"Type" seltabl:"Type" hSel:"tr:nth-child(1) th:nth-child(3)" dSel:"tr td:nth-child(3)" cSel:"$text"`
	Distance  string `seltabl:"Distance" hSel:"tr:nth-child(1) th:nth-child(4)" dSel:"tr td:nth-child(4)" json:"Distance" `
	Notes     string `json:"Notes" seltabl:"Notes" hSel:"tr:nth-child(1) th:nth-child(5)" dSel:"tr td:nth-child(5)"`
}

//go:embed supernova.html
var SuperNovaTable string

// SuperNovaTableResult is the expected result of parsing the supernova table
var SuperNovaTableResult = []SuperNovaStruct{
	{Supernova: "SN 1006", Year: "1006", Type: "Type Ia", Distance: "7,200", Notes: "Brightest recorded supernova in history"},
	{Supernova: "SN 1054 (Crab Nebula)", Year: "1054", Type: "Type II", Distance: "6,500", Notes: "Formed the Crab Nebula and pulsar"},
	{Supernova: "SN 1572 (Tycho's Supernova)", Year: "1572", Type: "Type Ia", Distance: "8,000-10,000", Notes: "Observed by Tycho Brahe"},
	{Supernova: "SN 1604 (Kepler's Supernova)", Year: "1604", Type: "Type Ia", Distance: "20,000", Notes: "Last observed supernova in the Milky Way"},
	{Supernova: "SN 1987A", Year: "1987", Type: "Type II", Distance: "168,000", Notes: "Closest observed supernova since 1604"},
	{Supernova: "SN 1993J", Year: "1993", Type: "Type IIb", Distance: "11,000,000", Notes: "In the galaxy M81"},
}
