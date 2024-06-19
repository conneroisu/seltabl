package parsers

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
)

var errNoURLs = errors.New("no urls found in file")
var errIgnoreOnly = errors.New("ignores only found in file")

// StructCommentData holds the parsed URLs and ignore-elements.
type StructCommentData struct {
	URLs           []string
	IgnoreElements []string
}

// ParseStructComments parses the comments from struct type declarations in the provided Go source code
// and extracts @url and @ignore-elements into separate arrays.
func ParseStructComments(src string) (StructCommentData, error) {
	fset := token.NewFileSet()
	// Parse the source code into an ast.File.
	node, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return StructCommentData{}, err
	}
	var data StructCommentData
	urlPattern := regexp.MustCompile(`@url:\s*(\S+)`)
	ignorePattern := regexp.MustCompile(`@ignore-elements:\s*(.*)`)
	// Inspect the AST to find struct type declarations and their comments
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.GenDecl:
			// Look for type declarations
			if t.Tok == token.TYPE {
				for _, spec := range t.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					// Check if the type specification is a struct
					if _, ok := typeSpec.Type.(*ast.StructType); ok {
						// Add comments associated with the struct type declaration
						if t.Doc != nil {
							for _, comment := range t.Doc.List {
								text := strings.TrimSpace(comment.Text)
								// Extract @url
								if urlMatches := urlPattern.FindStringSubmatch(text); len(urlMatches) > 1 {
									data.URLs = append(data.URLs, urlMatches[1])
								}
								// Extract @ignore-elements
								if ignoreMatches := ignorePattern.FindStringSubmatch(text); len(ignoreMatches) > 1 {
									elements := strings.Split(ignoreMatches[1], ",")
									for _, elem := range elements {
										data.IgnoreElements = append(data.IgnoreElements, strings.TrimSpace(elem))
									}
								}
							}
							if len(data.URLs) == 0 && len(data.IgnoreElements) == 0 {
							}
						}
					}
				}
			}
		}
		return true
	})
	if len(data.URLs) == 0 {
		return StructCommentData{}, errNoURLs
	}
	return data, nil
}
