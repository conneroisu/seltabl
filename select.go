package seltabl

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// SelectorInferface is an interface for running a goquery selector on a cellValue
//
// It is an interface that defines a Select method that takes a cellValue (goquery.Selection)
// and returns a string of the applied selection and an error.
type SelectorInferface interface {
	Select(cellValue *goquery.Selection) (string, error)
}

// selector is a struct for running a goquery selector on a cellValue
//
// It is a struct that satisfies the SelectorInferface interface.
//
// It contains a control tag and a query selector.
type selector struct {
	control string
	query   string
}

// Select runs the selector on the cellValue and sets the cellText
// and returns the cellText
//
// It returns the output of running a  selector and an error if the selector is not supported or fails.
func (s selector) Select(cellValue *goquery.Selection) (string, error) {
	var cellText string
	var exists bool
	switch s.control {
	case ctlInnerTextSelector:
		cellText = cellValue.Text()
		cellText = strings.TrimSpace(cellText)
		if cellValue.Length() == 0 {
			return "", fmt.Errorf("failed to find selector: %s", s.control)
		}
	case ctlAttrSelector:
		cellText, exists = cellValue.Attr(s.query)
		if !exists {
			return "", fmt.Errorf("failed to find selector: %s", s.control)
		}
	default:
		return "", fmt.Errorf(
			"unsupported identifer: %s (identifers are %s)",
			s.control,
			strings.Join(cSels, " "),
		)
	}
	return cellText, nil
}
