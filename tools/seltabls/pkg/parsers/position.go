package parsers

import (
	"go/ast"
	"go/token"
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
	// StateInvalid is the state for when the position is invalid or not within a struct.
	StateInvalid
)

// String returns the string representation of the State.
func (s State) String() string {
	return [...]string{
		"StateInTag",
		"StateInTagValue",
		"StateAfterColon",
		"StateInvalid",
	}[s]
}

// PositionInStructTagValue checks if the given position is within the value of a struct tag.
func PositionInStructTagValue(
	node *ast.StructType,
	pos protocol.Position,
	fset *token.FileSet,
) (
	result string,
	exists bool,
) {
	closestTagValue := ""
	closestDistance := int(
		^uint(0) >> 1,
	) // Initialize with the maximum possible integer value
	for _, field := range node.Fields.List {
		if field.Tag != nil {
			inNode := IsPositionInNode(field.Tag, pos, fset)
			if inNode {
				tagValue := field.Tag.Value
				tagContent := strings.Trim(tagValue, "`")
				start := fset.Position(field.Tag.Pos())
				for i := 0; i < len(tagContent); i++ {
					if tagContent[i] == '"' {
						startQuote := i + 1
						endQuote := strings.Index(
							tagContent[startQuote:],
							"\"",
						)
						if endQuote == -1 {
							continue
						}
						endQuote += startQuote
						tagRow := start.Line
						tagColumnStart := start.Column + startQuote
						tagColumnEnd := start.Column + endQuote
						if int(pos.Line) == tagRow &&
							int(pos.Character) >= tagColumnStart &&
							int(pos.Character) <= tagColumnEnd {
							return tagContent[startQuote:endQuote], true
						}
						i = endQuote
					}
				}
			}
			// Calculate the distance to the current tag for the closest tag logic
			start := fset.Position(field.Tag.Pos())
			distance := (start.Line-int(pos.Line))*(start.Line-int(pos.Line)) + (start.Column-int(pos.Character))*(start.Column-int(pos.Character))
			if distance < closestDistance {
				closestDistance = distance
				closestTagValue = strings.Trim(field.Tag.Value, "`")
			}
		}
	}
	return closestTagValue, false
}

// IsPositionInNode checks if a given position is within the range of an AST node.
func IsPositionInNode(
	node ast.Node,
	pos protocol.Position,
	fset *token.FileSet,
) bool {
	start := fset.Position(node.Pos())
	end := fset.Position(node.End())
	// Check if the position is within the node's range
	if (int(pos.Line) > start.Line || (pos.Line == uint32(start.Line) && pos.Character >= uint32(start.Column))) &&
		(pos.Line < uint32(end.Line) || (pos.Line == uint32(end.Line) && pos.Character <= uint32(end.Column))) {
		return true
	}
	return false
}

// FindStructNodes finds the struct nodes in the AST.
func FindStructNodes(node ast.Node) (structNodes []*ast.StructType) {
	ast.Inspect(node, func(n ast.Node) bool {
		if n != nil {
			if s, ok := n.(*ast.StructType); ok {
				structNodes = append(structNodes, s)
				return false
			}
		}
		return true
	})
	return structNodes
}

// IsPositionInTag checks if the given position is within a struct tag.
func IsPositionInTag(
	node *ast.StructType,
	pos protocol.Position,
	fset *token.FileSet,
) bool {
	for _, field := range node.Fields.List {
		if field.Tag != nil {
			if IsPositionInNode(field.Tag, pos, fset) {
				return true
			}
		}
	}
	return false
}

// PositionBeforeValue returns the value of the position in a file
func PositionBeforeValue(
	pos protocol.Position,
	text *string,
) byte {
	split := strings.Split(*text, "\n")
	if int(pos.Line) > len(split) {
		return '\n'
	}
	line := split[pos.Line]
	if int(pos.Character) > len(line) {
		return '\n'
	}
	return line[pos.Character-1]
}
