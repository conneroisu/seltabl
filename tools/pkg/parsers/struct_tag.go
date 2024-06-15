package parsers

import (
	"fmt"
	"strings"

	"github.com/conneroisu/seltabl/tools/pkg/lsp"
)

const (
	// tagMarker is the marker used to indicate the start of a struct tag
	tagMarker = ">>>"
	// tagMarkerEnd is the marker used to indicate the end of a struct tag
	tagMarkerEnd = "<<<"

	// StateOutsideTag is the state for when the tag is outside of a struct tag
	StateOutsideTag = 0
	// StateInsideTag is the state for inside a tag
	StateInsideTag = 1
	// StateInsideDoubleQuotes is the state for inside double quotes
	StateInsideDoubleQuotes = 2
)

// PositionStatusInStructTag checks if a given position is within a struct tag
// in a given source code string.
//
// Returns the following status codes:
// 0 - not in golang struct tag
// 1 - in golang struct tag
// 2 - in struct tag and inside double quotes
//
// An example of a struct tag where examples are given as the output status code for the position:
//
//	type TableStruct struct {
//		A string `json:"apple" seltabl:"apple" hSel:"tr:nth-child(1) td:nth-child(1)" dSel:"tr td:nth-child(1)" ctl:"text"`
//		B string `json:"bee" seltabl:"bee" hSel:"tr:nth-child(1) td:nth-child(2)" dSel:"tr td:nth-child(2)" ctl:"text"`
//	}
func PositionStatusInStructTag(src string, pos lsp.Position) (int, error) {
	lines := strings.Split(src, "\n")
	if pos.Line-1 < 0 || pos.Line-1 >= len(lines) {
		return 0, fmt.Errorf("line number out of range")
	}
	line := lines[pos.Line-1]
	if pos.Character-1 < 0 || pos.Character-1 >= len(line) {
		return 0, fmt.Errorf("character number out of range")
	}
	highlightedLine := line[:pos.Character-1] + tagMarker + string(
		line[pos.Character-1],
	) + tagMarkerEnd + line[pos.Character:]

	return determinePositionStatus(highlightedLine, pos.Character)
}

// determinePositionStatus determines the position status of a given line
func determinePositionStatus(line string, offset int) (int, error) {
	state := StateOutsideTag
	if !strings.Contains(line, "`") {
		fmt.Println("not found a ` in the line")
		return 0, nil
	}
	for i := 0; i < offset; i++ {
		char := line[i]

		switch state {
		case StateOutsideTag:
			if char == '`' {
				state = StateInsideTag
			}

		case StateInsideTag:
			if char == '`' {
				state = StateOutsideTag
			} else if char == '"' {
				state = StateInsideDoubleQuotes
			}

		case StateInsideDoubleQuotes:
			if char == '"' {
				state = StateInsideTag
			}
		}
	}
	return state, nil
}
