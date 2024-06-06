package seltabl

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SelectorInferface is an interface for running a goquery selector on a cellValue
type SelectorInferface interface {
	Run(cellValue *goquery.Selection) (*string, error)
}

// Selector is a struct for running a goquery selector on a cellValue
type selector struct {
	query string
}

// Run runs the selector on the cellValue and sets the cellText
// and returns the cellText
func (s *selector) Run(cellValue *goquery.Selection) (*string, error) {
	var cellText string
	var exists bool
	switch s.query {
	case innerTextSelector:
		cellText = cellValue.Text()
		cellText = strings.TrimSpace(cellText)
		if cellValue.Length() == 0 {
			return nil, fmt.Errorf("failed to find selector: %s", s.query)
		}
	case attrSelector:
		cellText, exists = cellValue.Attr(s.query)
		if !exists {
			return nil, fmt.Errorf("failed to find selector: %s", s.query)
		}
	default:
		print("default")
	}
	return &cellText, nil
}
