package parsers

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetAllSelectors retrieves all selectors from the given HTML document
func GetAllSelectors(doc *goquery.Document) ([]string, error) {
	strs := []string{}
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		str := getSelectorsFromSelection(s)
		if str != "" {
			if !contains(strs, str) {
				strs = append(strs, str)
			}
		}
	})
	if len(strs) == 0 {
		return nil, fmt.Errorf("no selectors found in document")
	}
	return strs, nil
}

// getSelectorsFromSelection returns the CSS selector for the given goquery selection
func getSelectorsFromSelection(s *goquery.Selection) string {
	if s.Length() == 0 {
		return ""
	}
	// Get the parent node
	parent := s.Parent()
	// Recursive call for the parent
	parentSelector := getSelectorsFromSelection(parent)
	// Get the selector for the current node
	currentSelector := singleSelector(s)
	// Combine the parent and current selectors
	if parentSelector != "" && currentSelector != "" {
		return parentSelector + " " + currentSelector
	} else if parentSelector != "" && currentSelector == "" {
		return parentSelector
	}
	return currentSelector
}

// singleSelector returns a single CSS selector for the given node
func singleSelector(selection *goquery.Selection) string {
	var selector string
	id, exists := selection.Attr("id")
	if exists {
		selector = fmt.Sprintf("%s#%s", goquery.NodeName(selection), id)
	}
	attr, exists := selection.Attr("class")
	if exists {
		selector = fmt.Sprintf(
			"%s[class=%s]",
			goquery.NodeName(selection),
			strings.Join(strings.Fields(attr), "."),
		)
	}
	attr, exists = selection.Attr("name")
	if exists {
		selector = fmt.Sprintf(
			"%s[name=%s]",
			goquery.NodeName(selection),
			attr,
		)
	}
	attr, exists = selection.Attr("type")
	if exists {
		selector = fmt.Sprintf(
			"%s[type=%s]",
			goquery.NodeName(selection),
			attr,
		)
	}
	attr, exists = selection.Attr("placeholder")
	if exists {
		selector = fmt.Sprintf(
			"%s[placeholder=%s]",
			goquery.NodeName(selection),
			attr,
		)
	}
	attr, exists = selection.Attr("value")
	if exists {
		selector = fmt.Sprintf(
			"%s[value=%s]",
			goquery.NodeName(selection),
			attr,
		)
	}
	attr, exists = selection.Attr("src")
	if exists {
		selector = fmt.Sprintf("%s[src=%s]", goquery.NodeName(selection), attr)
	}
	_, exists = selection.Attr("href")
	if exists {
		selector = fmt.Sprintf(
			"%s[href]",
			goquery.NodeName(selection),
		)
	}
	if selector == "" {
		selector = goquery.NodeName(selection)
	}
	return selector
}
