package seltabl

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SelectorInferface is an interface for running a goquery selector on a cellValue
type SelectorInferface interface {
	Select(cellValue *goquery.Selection) (*string, error)
}

// Selector is a struct for running a goquery selector on a cellValue
type selector struct {
	identifer string
	query     string
}

// Select runs the selector on the cellValue and sets the cellText
// and returns the cellText
func (s selector) Select(cellValue *goquery.Selection) (*string, error) {
	var cellText string
	var exists bool
	switch s.identifer {
	case innerTextSelector:
		cellText = cellValue.Text()
		cellText = strings.TrimSpace(cellText)
		if cellValue.Length() == 0 {
			return nil, fmt.Errorf("failed to find selector: %s", s.identifer)
		}
	case attrSelector:
		cellText, exists = cellValue.Attr(s.query)
		if !exists {
			return nil, fmt.Errorf("failed to find selector: %s", s.identifer)
		}
	default:
		print("default")
	}
	return &cellText, nil
}
