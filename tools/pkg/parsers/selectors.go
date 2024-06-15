package parsers

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// GetAllSelectors retrieves all selectors from the given HTML document
func GetAllSelectors(doc *goquery.Document) []string {
	strs := []string{}
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		str := getSelectorsFromSelection(s)
		if str != "" {
			if !contains(strs, str) {
				strs = append(strs, str)
			}
		}
	})
	return strs
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
	// should output html body div#ContentArea table tbody tr.heading td a[href=https://example.com]
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

	if id, exists := selection.Attr("id"); exists {
		selector = fmt.Sprintf("%s#%s", goquery.NodeName(selection), id)
	} else if class, exists := selection.Attr("class"); exists {
		selector = fmt.Sprintf("%s.%s", goquery.NodeName(selection), strings.Join(strings.Fields(class), "."))
	} else if attr, exists := selection.Attr("name"); exists {
		selector = fmt.Sprintf("%s[name=%s]", goquery.NodeName(selection), attr)
	} else if attr, exists := selection.Attr("type"); exists {
		selector = fmt.Sprintf("%s[type=%s]", goquery.NodeName(selection), attr)
	} else if attr, exists := selection.Attr("placeholder"); exists {
		selector = fmt.Sprintf("%s[placeholder=%s]", goquery.NodeName(selection), attr)
	} else if attr, exists := selection.Attr("value"); exists {
		selector = fmt.Sprintf("%s[value=%s]", goquery.NodeName(selection), attr)
	} else if attr, exists := selection.Attr("src"); exists {
		selector = fmt.Sprintf("%s[src=%s]", goquery.NodeName(selection), attr)
	} else if attr, exists := selection.Attr("href"); exists {
		selector = fmt.Sprintf("%s[href=%s]", goquery.NodeName(selection), attr)
	} else {
		selector = goquery.NodeName(selection)
	}

	return selector
}
