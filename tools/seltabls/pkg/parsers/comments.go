// Package parsers provides functions for parsing comments in Go source code.
//
// It includes functions for extracting URLs and ignore-elements from comments.
package parsers

import (
	"fmt"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

var (
	errIgnoreOnly = fmt.Errorf("ignore elements found but no URLs")
	errNoURLs     = fmt.Errorf("no URLs found")

	// @url: <url>
	urlPattern = regexp.MustCompile(`@url:\s*(\S+)`)
	// @ignore-elements: <element1>, <element2>, ...
	ignorePattern = regexp.MustCompile(`@ignore:\s*(.*)`)
	// @occurrences: <number>
	occurrencesPattern = regexp.MustCompile(`@occurrences:\s*(\d+)`)
)

// StructCommentData holds the parsed URLs and ignore-elements.
type StructCommentData struct {
	URLs           []string
	IgnoreElements [][]string
	Occurrences    []int
}

// ParseStructComments parses the comments from struct type declarations in the
// provided Go source code and extracts @url and @ignore-elements into separate
// arrays.
func ParseStructComments(src string) (StructCommentData, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("parsers.ParseStructComments recovered: %v", r)
		}
	}()
	node, err := decorator.Parse(src)
	if err != nil {
		return StructCommentData{}, err
	}
	var data StructCommentData
	// Inspect the AST to find struct type declarations and their comments.
	dst.Inspect(node, func(n dst.Node) bool {
		switch t := n.(type) {
		case *dst.GenDecl:
			// Look for type declarations.
			if t.Tok == token.TYPE {
				for _, spec := range t.Specs {
					typeSpec, ok := spec.(*dst.TypeSpec)
					if !ok {
						continue
					}
					// Check if the type specification is a struct
					if _, ok := typeSpec.Type.(*dst.StructType); ok && t.Decs.NodeDecs.Start != nil {
						for _, comment := range t.Decs.NodeDecs.Start {
							text := strings.TrimSpace(comment)
							// Extract @url of type string
							if urlMatches := urlPattern.FindStringSubmatch(text); len(urlMatches) > 1 {
								data.URLs = append(data.URLs, urlMatches[1])
							}
							// Extract @ignore of type []string
							if ignoreMatches := ignorePattern.FindStringSubmatch(text); len(ignoreMatches) > 1 {
								elements := strings.Split(ignoreMatches[1], ",")
								for _, elem := range elements {
									data.IgnoreElements = append(data.IgnoreElements, []string{strings.TrimSpace(elem)})
								}
							}
							// Extract @occurrences of type int
							if occurrencesMatches := occurrencesPattern.FindStringSubmatch(text); len(occurrencesMatches) > 1 {
								occurrences, err := strconv.Atoi(occurrencesMatches[1])
								if err != nil {
									return false
								}
								data.Occurrences = append(data.Occurrences, occurrences)
							}
						}
					}
				}
			}
		}
		return true
	})
	hasNoURLs := len(data.URLs) == 0
	if hasNoURLs && len(data.IgnoreElements) != 0 {
		return StructCommentData{}, errIgnoreOnly
	}
	if hasNoURLs {
		return StructCommentData{}, errNoURLs
	}
	return data, nil
}
