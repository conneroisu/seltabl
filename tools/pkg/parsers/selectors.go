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
			found := false
			for _, str2 := range strs {
				if str2 == str {
					found = true
					break
				}
			}
			if !found {
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
	if parentSelector != "" && currentSelector != "" {
		return parentSelector + " " + currentSelector
	} else if parentSelector != "" {
		return parentSelector
	}
	return currentSelector
}

// singleSelector returns a single CSS selector for the given node
func singleSelector(selection *goquery.Selection) string {
	var selector string

	if id, exists := selection.Attr("id"); exists {
		selector = fmt.Sprintf("#%s", id)
	} else if class, exists := selection.Attr("class"); exists {
		selector = fmt.Sprintf("%s.%s", goquery.NodeName(selection), strings.Join(strings.Fields(class), "."))
	} else if attr, exists := selection.Attr("name"); exists {
		selector = fmt.Sprintf("[name=%s]", attr)
	} else if attr, exists := selection.Attr("type"); exists {
		selector = fmt.Sprintf("[type=%s]", attr)
	} else if attr, exists := selection.Attr("placeholder"); exists {
		selector = fmt.Sprintf("[placeholder=%s]", attr)
	} else if attr, exists := selection.Attr("value"); exists {
		selector = fmt.Sprintf("[value=%s]", attr)
	} else if attr, exists := selection.Attr("src"); exists {
		selector = fmt.Sprintf("[src=%s]", attr)
	} else if attr, exists := selection.Attr("href"); exists {
		selector = fmt.Sprintf("[href=%s]", attr)
	} else {
		selector = goquery.NodeName(selection)
	}

	return selector
}
