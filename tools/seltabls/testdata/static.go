package testdata

import (
	_ "embed"
)

//go:embed example.go
var ExampleGo string

//go:embed mainex.go
var MainExGo string
