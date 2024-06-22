package testdata

import (
	_ "embed"
)

// FixtureStruct is a test struct
type FixtureStruct struct {
	A string `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" ctl:"text"`
	B string `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" ctl:"text"`
}

//go:embed example.go
var ExampleGo string

//go:embed mainex.go
var MainExGo string
