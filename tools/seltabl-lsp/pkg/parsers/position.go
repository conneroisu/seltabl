package parsers

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/conneroisu/seltabl/tools/seltabl-lsp/pkg/lsp"
)

// isPositionInStructTagValue checks if the given position is within the value of a struct tag.
func isPositionInStructTagValue(node *ast.StructType, pos lsp.Position, fset *token.FileSet) bool {
	for _, field := range node.Fields.List {
		if field.Tag != nil {
			if IsPositionInNode(field.Tag, pos, fset) {
				tagValue := field.Tag.Value
				tagContent := strings.Trim(tagValue, "`")
				start := fset.Position(field.Tag.Pos())

				for i := 0; i < len(tagContent); i++ {
					if tagContent[i] == '"' {
						startQuote := i + 1
						endQuote := strings.Index(tagContent[startQuote:], "\"")
						if endQuote == -1 {
							continue
						}
						endQuote += startQuote
						tagRow := start.Line
						tagColumnStart := start.Column + startQuote
						tagColumnEnd := start.Column + endQuote

						if pos.Line == tagRow && pos.Character >= tagColumnStart && pos.Character <= tagColumnEnd {
							return true
						}
						i = endQuote
					}
				}
			}
		}
	}
	return false
}

// IsPositionInNode checks if a given position is within the range of an AST node.
func IsPositionInNode(node ast.Node, pos lsp.Position, fset *token.FileSet) bool {
	start := fset.Position(node.Pos())
	end := fset.Position(node.End())

	// Check if the position is within the node's range
	if (pos.Line > start.Line || (pos.Line == start.Line && pos.Character >= start.Column)) &&
		(pos.Line < end.Line || (pos.Line == end.Line && pos.Character <= end.Column)) {
		return true
	}
	return false
}

// FindStructNode finds the struct node in the AST.
func FindStructNode(node ast.Node) (structNodes []*ast.StructType) {
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
func IsPositionInTag(node *ast.StructType, pos lsp.Position, fset *token.FileSet) bool {
	for _, field := range node.Fields.List {
		if field.Tag != nil {
			if IsPositionInNode(field.Tag, pos, fset) {
				return true
			}
		}
	}
	return false
}
