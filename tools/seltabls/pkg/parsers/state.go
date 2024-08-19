package parsers

import (
	"fmt"
	"strings"

	"go.lsp.dev/protocol"
)

// State represents the state of a position within a struct.
type State int

const (
	// StateInTag is the state for when the position is within a struct tag.
	StateInTag State = iota
	// StateInTagValue is the state for when the position is within a struct tag value.
	StateInTagValue
	// StateAfterColon is the state for when the position is after a colon.
	StateAfterColon
	// StateOnURL is the state for when the position is on a url.
	StateOnURL
	// StateInvalid is the state for when the position is invalid or not within a struct.
	StateInvalid
)

// String returns the string representation of the State.
func (s State) String() string {
	return [...]string{
		"StateInTag",
		"StateInTagValue",
		"StateAfterColon",
		"StateOnURL",
		"StateInvalid",
	}[s]
}

// ParsePosState parses the state of a position within a struct.
func ParsePosState(pos protocol.Position, text *string) (State, error) {
	split := strings.Split(*text, "\n")
	if int(pos.Line) > len(split) {
		return StateInvalid, fmt.Errorf("position line is greater than text length")
	}
	line := split[pos.Line-1]
	if int(pos.Character) > len(line) {
		return StateInvalid, fmt.Errorf("position character is greater than text length")
	}

	if strings.Contains(line, "@url:") {
		return StateOnURL, nil
	}

	isInStruct := false
	for i := int(pos.Line) - 1; i >= 0; i-- {
		if strings.Contains(split[i], "type") && strings.Contains(split[i], "struct") {
			isInStruct = true
			break
		}
		if strings.TrimSpace(split[i]) == "}" {
			break
		}
	}
	if !isInStruct {
		return StateInvalid, fmt.Errorf("position is not within a defined struct")
	}
	if strings.Count(line, "`") != 2 {
		return StateInvalid, nil
	}
	// if there is an odd number of double quotes to the left of the position, then it is in a tag
	left := strings.Count(line[:pos.Character], "\"")
	right := strings.Count(line[pos.Character:], "\"")
	if left%2 == 1 && right%2 == 1 {
		return StateInTagValue, nil
	}
	// if the character before the position is a colon, then it is after a colon
	before := line[pos.Character-1]
	if before == ':' {
		return StateAfterColon, nil
	}
	return StateInTag, nil
}
