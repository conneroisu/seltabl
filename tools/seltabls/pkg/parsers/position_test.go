package parsers

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
	"github.com/stretchr/testify/assert"
)

// TestFindStructNode tests the findStructNode function with a valid struct node
// as input.
func TestFindStructNode(t *testing.T) {
	source := `package main

type MyStruct struct {
	Field1 string ` + "`json:\"field1\"`" + `
	Field2 int    ` + "`json:\"field2\"`" + `
}`
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "example.go", source, parser.AllErrors)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	structNode := FindStructNodes(node)
	if structNode == nil {
		t.Error("Expected to find struct node, but did not find any")
	}
}

// TestIsPositionInStructTag tests the isPositionInStructTag function.
func TestIsPositionInStructTag(t *testing.T) {
	source := `package main

type MyStruct struct {
	Field1 string ` + "`json:\"field1\"`" + `
	Field2 int    ` + "`json:\"field2\"`" + `
}`
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "example.go", source, parser.AllErrors)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	structNodes := FindStructNodes(node)
	if structNodes == nil {
		t.Fatalf("Expected to find struct node, but did not find any")
	}

	// Define test cases
	testCases := []struct {
		name     string
		position lsp.Position
		expected bool
	}{
		{"Position in Field1 tag", lsp.Position{Line: 4, Character: 20}, true},
		{"Position in Field2 tag", lsp.Position{Line: 5, Character: 20}, true},
		{"Position outside tags", lsp.Position{Line: 3, Character: 1}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			failed := false
			for _, structNode := range structNodes {
				result := IsPositionInTag(structNode, tc.position, fset)
				if result != tc.expected {
					failed = true
					t.Errorf("expected %v, got %v", tc.expected, result)

				}
			}
			if failed {
				t.Fail()
			}
		})
	}
}

// TestIsPositionInStructTagValue tests the PositionInStructTagValue function.
func TestIsPositionInStructTagValue(t *testing.T) {
	source := `package main

	// @url: https://example.com
type MyStruct struct {
	Field1 string ` + "`json:\"field1\"`" + `
	Field2 int    ` + "`json:\"field2\"`" + `
}`
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "example.go", source, parser.Trace)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	structNodes := FindStructNodes(node)
	if structNodes == nil {
		t.Fatalf("Expected to find struct node, but did not find any")
	}

	// Define test cases
	testCases := []struct {
		name     string
		position lsp.Position
		expected bool
	}{
		{
			"Position in Field1 tag value",
			lsp.Position{Line: 5, Character: 22},
			true,
		},
		{
			"Position in Field2 tag value",
			lsp.Position{Line: 6, Character: 22},
			true,
		},
		{
			"Position out Field1 tag value",
			lsp.Position{Line: 5, Character: 21},
			false,
		},
		{
			"Position out Field2 tag value",
			lsp.Position{Line: 6, Character: 21},
			false,
		},
		{
			"Position outside tag values",
			lsp.Position{Line: 4, Character: 1},
			false,
		},
		{"Position out of range", lsp.Position{Line: 7, Character: 99}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			failed := false
			for _, structNode := range structNodes {
				_, is := PositionInStructTagValue(
					structNode,
					tc.position,
					fset,
				)
				assert.Equal(t, tc.expected, is)
			}
			if failed {
				t.Fail()
			}
		})
	}
}

func TestIsPositionInStructTagValue2(t *testing.T) {
	source := `// Package main is the entry point for the command line tool
// a language server for the seltabl package called seltabls.
package main

import (
	"os"

	"github.com/conneroisu/seltabl/tools/seltabls/cmd"
)

// main is the entry point for the command line tool, a
// language server for the seltabl package
func main() {
	rs := &cmd.Root{Writer: os.Stdout}
	if err := cmd.Execute(rs); err != nil {
		rs.Logger.Println(err)
	}
}

// TableStruct is a test struct
// @url: https://stats.ncaa.org/game_upload/team_codes
// @ignore-elements: script, style, link, img, footer, header
type TableStruct struct {
	A string ` + "`" + `json:"a" seltabl:"a" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" ctl:"text"` + "`" + `
	B string ` + "`" + `json:"b" seltabl:"b" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" ctl:"text"` + "`" + `
}
`
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "main.go", source, parser.Trace)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	structNodes := FindStructNodes(node)
	if structNodes == nil {
		t.Fatalf("Expected to find struct node, but did not find any")
	}

	// Define test cases
	testCases := []struct {
		name     string
		position lsp.Position
		expected bool
	}{
		{
			"Position out Field1 tag value",
			lsp.Position{Line: 5, Character: 22},
			false,
		},
		{
			"Position out Field2 tag value",
			lsp.Position{Line: 6, Character: 22},
			false,
		},
		{
			"Position out Field1 tag value",
			lsp.Position{Line: 5, Character: 0},
			false,
		},
		{
			"Position out Field2 tag value",
			lsp.Position{Line: 6, Character: 0},
			false,
		},
		{
			"Position in Field1 tag value",
			lsp.Position{Line: 24, Character: 18},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			failed := false
			for _, structNode := range structNodes {
				_, is := PositionInStructTagValue(
					structNode,
					tc.position,
					fset,
				)
				if is != tc.expected {
					t.Errorf("expected %v, got %v", tc.expected, is)
				}
			}
			if failed {
				t.Fail()
			}
		})
	}
}
