package parsers

import (
	"strings"
	"testing"

	"github.com/conneroisu/seltabl/tools/pkg/lsp"
	"github.com/conneroisu/seltabl/tools/testdata"
)

// TestPositionStatusInStructTag tests the PositionStatusInStructTag function
func TestPositionStatusInStructTag(t *testing.T) {
	src := `
package main

type Person struct {
	Name string ` + "`json:\"name\"`" + `
	Age  int    ` + "`json:\"age\"`" + `
	Address string ` + "`json:\"address,omitempty\"`" + `
}

func main() {
	p := Person{Name: "Alice", Age: 30, Address: "123 Street"}
}
`
	tests := []struct {
		name      string
		line      int
		character int
		expected  int
	}{
		// tests for 1
		{"Inside double quotes of json:\"address,omitempty\"", 7, 20, 1},
		{"Inside double quotes of json:\"name\"", 5, 17, 1},
		{"Inside double quotes of json:\"name\"", 5, 19, 1},
		{"Inside double quotes of json:\"age\"", 6, 18, 1},
		// tests for 0
		{"Outside struct tag (Name field)", 5, 5, 0},
		{"Inside main function, not inside struct tag", 10, 5, 0},
		{"Inside struct field but not a tag", 4, 10, 0},
		{"Outside struct field but not a tag", 4, 5, 0},
		// tests for 2
		{"Inside double quotes of json:\"name\"", 5, 21, 2},
		{"Inside double quotes of json:\"age\"", 6, 21, 2},
		{"Inside double quotes of json:\"address,omitempty\"", 7, 24, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := lsp.Position{
				Line:      tt.line,
				Character: tt.character,
			}
			status, err := PositionStatusInStructTag(src, pos)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if status != tt.expected {
				highlightedSrc := highlightPosition(src, pos)
				t.Errorf("expected status %d, got %d\n%s", tt.expected, status, highlightedSrc)
			}
		})
	}
}

// highlightPosition adds markers at the specified position in the source code
// This is used to highlight the position in the source code for easier debugging
func highlightPosition(src string, pos lsp.Position) string {
	lines := strings.Split(src, "\n")
	if pos.Line-1 < 0 || pos.Line-1 >= len(lines) {
		return src
	}
	line := lines[pos.Line-1]
	if pos.Character-1 < 0 || pos.Character-1 >= len(line) {
		return src
	}
	highlightedLine := line[:pos.Character-1] + ">>>" + string(line[pos.Character-1]) + "<<<" + line[pos.Character:]
	lines[pos.Line-1] = highlightedLine
	return strings.Join(lines, "\n")
}

func TestPositionStatusInStructTag2(t *testing.T) {
	t.Run("Test ExtractSelectors", func(t *testing.T) {
		t.Parallel()
		src := testdata.MainExGo
		pos := lsp.Position{
			Line:      40,
			Character: 17,
		}
		status, err := PositionStatusInStructTag(src, pos)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status != 2 {
			highlightedSrc := highlightPosition(src, pos)
			t.Errorf("expected status %d, got %d\n%s", 1, status, highlightedSrc)
		}
	})
}
