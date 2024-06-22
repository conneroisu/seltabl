package analysis

import (
	"fmt"
	"strings"

	"github.com/conneroisu/seltabl/tools/seltabls/pkg/lsp"
)

// Hover returns a hover response for the given uri and position
func (s *State) Hover(
	id int,
	uri string,
	position lsp.Position,
) (*lsp.HoverResponse, error) {
	text := s.Documents[uri]
	split := strings.Split(text, "\n")
	line := split[position.Line]
	preN, preStr := preLine(line, "\"", position.Character)
	postN, postStr := postLine(line, "\"", position.Character)
	return &lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf(
				"preline:\n %s\n num \"'s: %d, postline:\n %s\n num \"'s: %d",
				preStr, preN, postStr, postN,
			) + fmt.Sprintf(
				"[%s](%s)",
				"google",
				"https://www.google.com",
			),
		},
	}, nil
}

// preLine returns the number of times the specified character is repeated in the given line
func preLine(line string, char string, col int) (count int, str string) {
	str = line[:col]
	for i := range str {
		if string(line[i]) == char {
			count++
		}
	}
	return count, str
}

// postLine returns the number of times the specified character is repeated in the given line
func postLine(line string, char string, col int) (count int, str string) {
	for i := range line[col:] {
		if string(line[i]) == char {
			count++
		}
	}
	return count, line[col:]
}
