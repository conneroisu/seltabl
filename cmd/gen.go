package main

import "github.com/dave/jennifer/jen"

func gen(packageName string, structName string) {
	file := jen.NewFile(packageName)
	file.PackageComment(`
		//go:build ignore
		// +build ignore
	`)
	file.ImportName("github.com/stretchr/testify/assert", "assert")
	file.ImportName("github.com/conneroisu/seltabl", "seltabl")
	file.ImportName("strings", "strings")
}
